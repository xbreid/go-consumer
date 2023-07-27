generate-schema:
	go run -mod=mod entgo.io/ent/cmd/ent generate --feature sql/upsert ./ent/schema