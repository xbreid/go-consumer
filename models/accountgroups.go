package models

import (
	"context"
	"fmt"
	"go-consumer/ent"
	"go-consumer/ent/accountgroup"
	"log"
)

type Models struct {
	AccountGroup AccountGroup
}

type AccountGroup struct {
	DisplayName string `json:"display_name"`
	ExternalId  string `json:"external_id"`
}

func UpsertAccountGroup(ctx context.Context, client *ent.Client, input AccountGroup) (*ent.AccountGroup, error) {
	id, err := client.AccountGroup.
		Create().
		SetDisplayName(input.DisplayName).
		SetExternalID(input.ExternalId).
		OnConflictColumns(accountgroup.FieldExternalID).
		UpdateNewValues().
		ID(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed creating user: %w", err)
	}
	log.Println("account group was upsert: ", id)

	current := client.AccountGroup.GetX(ctx, id)
	return current, nil
}
