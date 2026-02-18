package handlers_test

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"
	"todo-app/ent"

	entsql "entgo.io/ent/dialect/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v5"
	_ "github.com/mattn/go-sqlite3"
)

var (
	testClient *ent.Client
	testDB     *sql.DB
)

func TestMain(m *testing.M) {
	// Load env
	envFile := fmt.Sprintf("../envs/%s.env", os.Getenv("APP_ENV"))
	if err := godotenv.Load(envFile); err != nil {
		fmt.Printf("Error loading .env file: %v\n", err)
		os.Exit(1)
	}

	// Open DB connection with sql.Open to keep reference to *sql.DB
	var err error
	driverName := os.Getenv("TEST_DATABASE")
	dsn := os.Getenv("TEST_DATABASE_URL")

	testDB, err = sql.Open(driverName, dsn)
	if err != nil {
		log.Fatalf("failed opening connection to %s: %v", driverName, err)
	}

	defer func() {
		if err := testDB.Close(); err != nil {
			log.Fatalf("failed closing connection to %s: %v", driverName, err)
		}
	}()

	// Create ent.Client from sql.DB
	drv := entsql.OpenDB(driverName, testDB)
	testClient = ent.NewClient(ent.Driver(drv))

	// Run migration
	if err := testClient.Schema.Create(context.Background()); err != nil {
		log.Fatalf("failed creating schema resources: %v", err)
	}

	code := m.Run()
	os.Exit(code)
}

// cleanupDatabase cleans up the database by truncating (or deleting) all tables.
// This is necessary because we cannot use transaction rollback for isolation
// when the application code itself manages transactions.
func cleanupDatabase(t *testing.T) {
	ctx := context.Background()
	driver := os.Getenv("TEST_DATABASE")

	switch driver {
	case "mysql":
		// Disable FK checks
		if _, err := testDB.ExecContext(ctx, "SET FOREIGN_KEY_CHECKS = 0"); err != nil {
			t.Fatalf("failed to disable foreign keys: %v", err)
		}

		// Get all table names
		rows, err := testDB.QueryContext(ctx, "SHOW TABLES")
		if err != nil {
			t.Fatalf("failed to show tables: %v", err)
		}
		defer func() {
			if err := rows.Close(); err != nil {
				t.Errorf("failed closing rows: %v", err)
			}
		}()

		var tables []string
		for rows.Next() {
			var table string
			if err := rows.Scan(&table); err != nil {
				t.Fatalf("failed to scan table name: %v", err)
			}
			tables = append(tables, table)
		}

		// Truncate all tables
		for _, table := range tables {
			if _, err := testDB.ExecContext(ctx, fmt.Sprintf("TRUNCATE TABLE `%s`", table)); err != nil {
				t.Fatalf("failed to truncate table %s: %v", table, err)
			}
		}

		// Enable FK checks
		if _, err := testDB.ExecContext(ctx, "SET FOREIGN_KEY_CHECKS = 1"); err != nil {
			t.Fatalf("failed to enable foreign keys: %v", err)
		}

	case "sqlite3":
		// Disable FK checks
		if _, err := testDB.ExecContext(ctx, "PRAGMA foreign_keys = OFF"); err != nil {
			t.Fatalf("failed to disable foreign keys: %v", err)
		}

		// Get all table names
		rows, err := testDB.QueryContext(ctx, "SELECT name FROM sqlite_master WHERE type='table' AND name NOT LIKE 'sqlite_%'")
		if err != nil {
			t.Fatalf("failed to list tables: %v", err)
		}
		defer func() {
			if err := rows.Close(); err != nil {
				t.Errorf("failed closing rows: %v", err)
			}
		}()

		var tables []string
		for rows.Next() {
			var table string
			if err := rows.Scan(&table); err != nil {
				t.Fatalf("failed to scan table name: %v", err)
			}
			tables = append(tables, table)
		}

		// Delete from all tables
		for _, table := range tables {
			if _, err := testDB.ExecContext(ctx, fmt.Sprintf("DELETE FROM `%s`", table)); err != nil {
				t.Fatalf("failed to delete table %s: %v", table, err)
			}
		}
		// Reset autoincrement logic: ignore error as persistent table might not exist
		_, _ = testDB.ExecContext(ctx, "DELETE FROM sqlite_sequence")

		// Enable FK checks
		if _, err := testDB.ExecContext(ctx, "PRAGMA foreign_keys = ON"); err != nil {
			t.Fatalf("failed to enable foreign keys: %v", err)
		}

	default:
		t.Fatalf("unsupported driver: %s", driver)
	}
}

type avoidCSRF struct{}

func (a avoidCSRF) SetCSRFCookie(req *http.Request) {
	cookie := &http.Cookie{Name: "csrf_token", Value: "test"}
	req.AddCookie(cookie)
}

func (a avoidCSRF) SetCSRFHeader(req *http.Request) {
	req.Header.Set("X-CSRF-Token", "test")
}

func createToken(t *testing.T, userID int) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": float64(userID), // jwt deserializes numbers as float64
		"exp":     time.Now().Add(time.Hour).Unix(),
	})
	tokenString, err := token.SignedString([]byte("secret"))
	if err != nil {
		t.Fatalf("failed to sign token: %v", err)
	}
	return tokenString
}

func createAuthenticatedRequest(t *testing.T, method, target, body string, userID int) (*http.Request, *httptest.ResponseRecorder) {
	var req *http.Request
	if body != "" {
		req = httptest.NewRequest(method, target, strings.NewReader(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	} else {
		req = httptest.NewRequest(method, target, nil)
	}
	if userID != 0 {
		cookie := &http.Cookie{Name: "token", Value: createToken(t, userID)}
		req.AddCookie(cookie)
	}

	csrf := avoidCSRF{}
	csrf.SetCSRFCookie(req)
	csrf.SetCSRFHeader(req)

	rec := httptest.NewRecorder()
	return req, rec
}
