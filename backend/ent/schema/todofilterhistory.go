package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

// TodoFilterHistory holds the schema definition for the TodoFilterHistory entity.
type TodoFilterHistory struct {
	ent.Schema
}

// Annotations of the TodoFilterHistory.
func (TodoFilterHistory) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "todo_filter_histories"},
	}
}

// Fields of the TodoFilterHistory.
func (TodoFilterHistory) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).Default(func() uuid.UUID { return uuid.Must(uuid.NewV7()) }),
		field.Int("user_id"),
		field.String("query").MaxLen(400),
		field.String("function_name").MaxLen(100).Optional(),
		field.JSON("args", map[string]interface{}{}).Optional(),
		field.JSON("result_todo_ids", []int{}).Optional(),
		field.Time("created_at").Default(time.Now).Immutable(),
	}
}

// Edges of the TodoFilterHistory.
func (TodoFilterHistory) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("user", User.Type).Ref("todo_filter_histories").Unique().Field("user_id").Required(),
	}
}
