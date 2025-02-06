package sql

import (
	"context"
	"fmt"
	"log/slog"
	"slices"
	"strings"

	"github.com/jackc/pgx/v5"
)

type SchemaKey struct{}

func ChangeSchema(ctx context.Context, conn *pgx.Conn) bool {
	var schemas []string
	switch v := ctx.Value(SchemaKey{}).(type) {
	case string:
		schemas = append(schemas, v)
		if v != "public" {
			schemas = append(schemas, "public")
		}
	case []string:
		schemas = append(schemas, v...)
		if !slices.Contains(schemas, "public") {
			schemas = append(schemas, "public")
		}
	default:
		return true
	}

	switchQuery := fmt.Sprintf("SET search_path = %s", strings.Join(schemas, ", "))
	if _, err := conn.Exec(ctx, switchQuery); err != nil {
		slog.Error("Failed to change schema on database connection", "schemas", schemas, "error", err)
		return false
	}
	return true
}

func RestoreSchema(conn *pgx.Conn) bool {
	if _, err := conn.Exec(context.Background(), "SET search_path TO DEFAULT"); err != nil {
		slog.Error("Failed to restore schema to database connection")
		return false
	}
	return true
}
