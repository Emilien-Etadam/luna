package versions

import (
	"luna-backend/db/internal/migrations/internal/registry"
	migrationTypes "luna-backend/db/internal/migrations/types"
	"luna-backend/errors"
	"luna-backend/types"
)

func init() {
	registry.RegisterMigration(types.Ver(0, 1, 1), func(q *migrationTypes.MigrationQueries) *errors.ErrorTrace {
		_, err := q.Tx.Exec(
			q.Context,
			`
			INSERT INTO global_settings (key, value)
			VALUES ('public_calendar_enabled', 'false'::jsonb)
			ON CONFLICT (key) DO NOTHING;
			`,
		)
		if err != nil {
			return errors.New().
				AddErr(errors.LvlDebug, err).
				Append(errors.LvlDebug, "Could not insert public_calendar_enabled global setting")
		}
		return nil
	})
}
