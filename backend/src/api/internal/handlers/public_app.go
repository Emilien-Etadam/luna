package handlers

import (
	"luna-backend/api/internal/util"
	"luna-backend/cache"
	"luna-backend/errors"
	"luna-backend/types"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func requirePublicCalendarPublication(u *util.HandlerUtility) *errors.ErrorTrace {
	enabled, tr := u.Tx.Queries().GetPublicCalendarEnabled()
	if tr != nil {
		return tr
	}
	if !enabled {
		return errors.New().Status(http.StatusNotFound).
			Append(errors.LvlPlain, "Public calendar is disabled")
	}
	return nil
}

// GetPublicSources lists all sources with id, name, type only (no auth).
func GetPublicSources(c *gin.Context) {
	u := util.GetUtil(c)
	if tr := requirePublicCalendarPublication(u); tr != nil {
		u.Error(tr)
		return
	}

	rows, tr := u.Tx.Queries().ListPublicSourceSummaries()
	if tr != nil {
		u.Error(tr)
		return
	}

	out := make([]exposedSource, 0, len(rows))
	for _, r := range rows {
		out = append(out, exposedSource{
			Id:   r.Id,
			Name: r.Name,
			Type: r.Type,
		})
	}
	u.Success(&gin.H{"sources": out})
}

// GetPublicSourceCalendars lists calendars for a source (read-only, no secrets).
func GetPublicSourceCalendars(c *gin.Context) {
	u := util.GetUtil(c)
	if tr := requirePublicCalendarPublication(u); tr != nil {
		u.Error(tr)
		return
	}

	sourceId, tr := util.GetId(c, "source")
	if tr != nil {
		u.Error(tr)
		return
	}

	userId, tr := u.Tx.Queries().GetSourceOwner(sourceId)
	if tr != nil {
		u.Error(tr)
		return
	}

	source, tr := cache.GetCached(u.Config.Cache, userId, sourceId, u.Context, func() (types.Source, *errors.ErrorTrace) {
		return u.Tx.Queries().GetSource(userId, sourceId, u.Context, u.Config)
	})
	if tr != nil {
		u.Error(tr)
		return
	}

	calsFromSource, tr := source.GetCalendars(u.Tx.Queries())
	if tr != nil {
		u.Error(tr)
		return
	}

	cals, tr := u.Tx.Queries().MergeCalendarsOverridesReadOnly(calsFromSource)
	if tr != nil {
		u.Error(tr)
		return
	}

	converted := make([]exposedCalendar, len(cals))
	for i, cal := range cals {
		u.Config.Cache.Cache(userId, cal)
		converted[i] = exposedCalendar{
			Id:           cal.GetId(),
			Source:       cal.GetSource().GetId(),
			Name:         cal.GetName(),
			Desc:         "",
			Color:        cal.GetColor(),
			Overridden:   cal.GetOverridden(),
			CanEdit:      false,
			CanDelete:    false,
			CanAddEvents: false,
		}
	}
	u.Success(&gin.H{"calendars": converted})
}

// GetPublicCalendarEvents returns events for a calendar in a read-only transaction (no cache writes).
func GetPublicCalendarEvents(c *gin.Context) {
	u := util.GetUtil(c)
	if tr := requirePublicCalendarPublication(u); tr != nil {
		u.Error(tr)
		return
	}

	calendarId, tr := util.GetId(c, "calendar")
	if tr != nil {
		u.Error(tr)
		return
	}

	userId, tr := u.Tx.Queries().GetCalendarOwnerUserId(calendarId)
	if tr != nil {
		u.Error(tr)
		return
	}

	calendar, tr := cache.GetCached(u.Config.Cache, userId, calendarId, u.Context, func() (types.Calendar, *errors.ErrorTrace) {
		return u.Tx.Queries().GetCalendar(userId, calendarId, u.Context, u.Config)
	})
	if tr != nil {
		u.Error(tr)
		return
	}

	startStr := c.Query("start")
	startTime, err := time.Parse(time.RFC3339, startStr)
	if err != nil {
		u.Error(errors.New().
			AddErr(errors.LvlDebug, err).
			Append(errors.LvlPlain, "Missing or malformed start time"))
		return
	}
	endStr := c.Query("end")
	endTime, err := time.Parse(time.RFC3339, endStr)
	if err != nil {
		u.Error(errors.New().
			AddErr(errors.LvlDebug, err).
			Append(errors.LvlPlain, "Missing or malformed end time"))
		return
	}
	if startTime.After(endTime) {
		u.Error(errors.New().
			Append(errors.LvlPlain, "Start time must not be after end time"))
		return
	}
	if endTime.Sub(startTime) > time.Hour*24*365 {
		endTime = startTime.Add(time.Hour * 24 * 365)
	}

	eventsFromCal, tr := calendar.GetEvents(startTime, endTime, u.Tx.Queries())
	if tr != nil {
		u.Error(tr)
		return
	}

	expandedEvents := make([]types.Event, len(eventsFromCal))
	count := 0
	for _, event := range eventsFromCal {
		expanded, tr := types.ExpandRecurrence(event, &startTime, &endTime)
		if tr != nil {
			u.Error(tr)
			return
		}

		if len(expanded) > 1 {
			newRes := make([]types.Event, len(expandedEvents)-1+len(expanded))
			copy(newRes, expandedEvents[:count])
			expandedEvents = newRes
		}

		for _, e := range expanded {
			expandedEvents[count] = e
			count++
		}
	}

	events, tr := u.Tx.Queries().MergeEventOverridesReadOnly(expandedEvents[:count])
	if tr != nil {
		u.Error(tr)
		return
	}

	converted := make([]exposedEvent, 0, len(events))
	for _, event := range events {
		if event.GetName() == "" {
			continue
		}
		converted = append(converted, exposedEvent{
			Id:         event.GetId(),
			Calendar:   event.GetCalendar().GetId(),
			Name:       event.GetName(),
			Desc:       "",
			Color:      event.GetColor(),
			Date:       event.GetDate(),
			Overridden: event.GetOverridden(),
			CanEdit:    false,
			CanDelete:  false,
		})
	}

	u.Success(&gin.H{"events": converted})
}
