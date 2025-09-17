package db

import "embed"

//go:embed "migrations/*.sql" "seed/*.sql"
var EmbedMigrations embed.FS
