package migrations

import . "github.com/grafana/grafana/pkg/services/sqlstore/migrator"

func addPreferencesMigrations(mg *Migrator) {
	preferencesV1 := Table{
		Name: "preferences",
		Columns: []*Column{
			{Name: "id", Type: DB_BigInt, IsPrimaryKey: true, IsAutoIncrement: true},
			{Name: "org_id", Type: DB_BigInt, Nullable: true},
			{Name: "user_id", Type: DB_BigInt, Nullable: true},
			{Name: "home_dashboard_id", Type: DB_BigInt, Nullable: true},
			{Name: "theme", Type: DB_NVarchar, Length: 255, Nullable: true},
			{Name: "created", Type: DB_DateTime, Nullable: false},
			{Name: "updated", Type: DB_DateTime, Nullable: false},
		},
		Indices: []*Index{
			{Cols: []string{"org_id"}},
			{Cols: []string{"user_id"}},
		},
	}

	mg.AddMigration("create preferences table v1", NewAddTableMigration(preferencesV1))
	addTableIndicesMigrations(mg, "v1", preferencesV1)
}
