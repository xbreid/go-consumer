// Code generated by ent, DO NOT EDIT.

package ent

import (
	"go-consumer/ent/accountgroup"
	"go-consumer/ent/schema"
	"time"

	"github.com/google/uuid"
)

// The init function reads all schema descriptors with runtime code
// (default values, validators, hooks and policies) and stitches it
// to their package variables.
func init() {
	accountgroupFields := schema.AccountGroup{}.Fields()
	_ = accountgroupFields
	// accountgroupDescCreatedAt is the schema descriptor for created_at field.
	accountgroupDescCreatedAt := accountgroupFields[3].Descriptor()
	// accountgroup.DefaultCreatedAt holds the default value on creation for the created_at field.
	accountgroup.DefaultCreatedAt = accountgroupDescCreatedAt.Default.(func() time.Time)
	// accountgroupDescID is the schema descriptor for id field.
	accountgroupDescID := accountgroupFields[0].Descriptor()
	// accountgroup.DefaultID holds the default value on creation for the id field.
	accountgroup.DefaultID = accountgroupDescID.Default.(func() uuid.UUID)
}
