package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

// AccountGroup holds the schema definition for the AccountGroup entity.
type AccountGroup struct {
	ent.Schema
}

// Fields of the AccountGroup.
func (AccountGroup) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Default(uuid.New).
			StorageKey("oid"),
		field.String("display_name").
			Optional(),
		field.String("external_id").Unique(),
	}
}

// Edges of the AccountGroup.
func (AccountGroup) Edges() []ent.Edge {
	return nil
}
