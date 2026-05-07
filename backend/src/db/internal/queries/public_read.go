package queries

import (
	"luna-backend/errors"
	"luna-backend/types"
	"net/http"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

// PublicSourceSummary is a minimal source row for anonymous read-only UI (no secrets).
type PublicSourceSummary struct {
	Id   types.ID
	Name string
	Type string
}

// ListPublicSourceSummaries returns every source in the system with id, name, type only.
func (q *Queries) ListPublicSourceSummaries() ([]PublicSourceSummary, *errors.ErrorTrace) {
	rows, err := q.Tx.Query(
		q.Context,
		`
		SELECT id, name, type::text
		FROM sources
		ORDER BY userid, display_order;
		`,
	)
	if err != nil {
		return nil, errors.New().Status(http.StatusInternalServerError).
			AddErr(errors.LvlDebug, err).
			Append(errors.LvlWordy, "Could not list sources for public read")
	}
	defer rows.Close()

	var out []PublicSourceSummary
	for rows.Next() {
		var idUUID uuid.UUID
		var row PublicSourceSummary
		if err := rows.Scan(&idUUID, &row.Name, &row.Type); err != nil {
			return nil, errors.New().Status(http.StatusInternalServerError).
				AddErr(errors.LvlDebug, err).
				Append(errors.LvlWordy, "Could not scan public source row")
		}
		row.Id = types.IdFromUuid(idUUID)
		out = append(out, row)
	}
	if err := rows.Err(); err != nil {
		return nil, errors.New().Status(http.StatusInternalServerError).
			AddErr(errors.LvlDebug, err).
			Append(errors.LvlWordy, "Could not iterate public sources")
	}
	return out, nil
}

// GetCalendarOwnerUserId returns the user id that owns the calendar (via its source).
func (q *Queries) GetCalendarOwnerUserId(calendarId types.ID) (types.ID, *errors.ErrorTrace) {
	var userUUID uuid.UUID
	err := q.Tx.QueryRow(
		q.Context,
		`
		SELECT s.userid
		FROM calendars c
		JOIN sources s ON c.source = s.id
		WHERE c.id = $1
		`,
		calendarId.UUID(),
	).Scan(&userUUID)
	switch err {
	case nil:
		return types.IdFromUuid(userUUID), nil
	case pgx.ErrNoRows:
		return types.EmptyId(), errors.New().Status(http.StatusNotFound).
			Append(errors.LvlDebug, "Calendar %v not found", calendarId).
			AltStr(errors.LvlPlain, "Calendar not found")
	default:
		return types.EmptyId(), errors.New().Status(http.StatusInternalServerError).
			AddErr(errors.LvlDebug, err).
			Append(errors.LvlWordy, "Could not resolve calendar owner")
	}
}

// CalendarBelongsToSource is true when the calendar row is under the given source.
func (q *Queries) CalendarBelongsToSource(calendarId, sourceId types.ID) (bool, *errors.ErrorTrace) {
	var exists bool
	err := q.Tx.QueryRow(
		q.Context,
		`SELECT EXISTS (SELECT 1 FROM calendars WHERE id = $1 AND source = $2)`,
		calendarId.UUID(),
		sourceId.UUID(),
	).Scan(&exists)
	if err != nil {
		return false, errors.New().Status(http.StatusInternalServerError).
			AddErr(errors.LvlDebug, err).
			Append(errors.LvlWordy, "Could not verify calendar source")
	}
	return exists, nil
}
