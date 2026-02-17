package handlers_test

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"
	"todo-app/ent"

	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v5"
	_ "github.com/mattn/go-sqlite3"
)

var testClient *ent.Client

func TestMain(m *testing.M) {
	// Load env
	envFile := fmt.Sprintf("../envs/%s.env", os.Getenv("APP_ENV"))
	if err := godotenv.Load(envFile); err != nil {
		fmt.Printf("Error loading .env file: %v\n", err)
		os.Exit(1)
	}

	// Open DB connection
	var err error
	driver := os.Getenv("TEST_DATABASE")
	dsn := os.Getenv("TEST_DATABASE_URL")
	testClient, err = ent.Open(driver, dsn)
	if err != nil {
		log.Fatalf("failed opening connection to sqlite: %v", err)
	}
	defer func() {
		if err := testClient.Close(); err != nil {
			log.Fatalf("failed closing connection to sqlite: %v", err)
		}
	}()

	// Run migration
	if err := testClient.Schema.Create(context.Background()); err != nil {
		log.Fatalf("failed creating schema resources: %v", err)
	}

	code := m.Run()
	os.Exit(code)
}

// withTx starts a transaction and rolls it back at the end.
// It provides the transactional client to the callback function.
func withTx(t *testing.T, fn func(*ent.Client)) {
	tx, err := testClient.Tx(context.Background())
	if err != nil {
		t.Fatalf("failed to start transaction: %v", err)
	}
	defer func() {
		if v := recover(); v != nil {
			_ = tx.Rollback()
			panic(v)
		}
		_ = tx.Rollback()
	}()
	fn(tx.Client())
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
