package handlers

import (
	"fmt"
	"luna-backend/api/internal/util"
	"luna-backend/errors"
	"luna-backend/types"
	"net/http"
	"strings"
	"time"

	ics "github.com/arran4/golang-ical"
	"github.com/gin-gonic/gin"
)

func HeadPublicCalendarICS(c *gin.Context) {
	u := util.GetUtil(c)
	ip := util.DetermineClientAddress(c).String()
	u.Logger.Infof("public calendar HEAD ip=%s ua=%s", ip, sanitizeUserAgent(c.Request.User-Agent()))

	enabled, tr := u.Tx.Queries().GetPublicCalendarEnabled()
	if tr != nil {
		u.Error(tr)
		return
	}
	if !enabled {
		u.SuccessRawBytes(http.StatusNotFound, []byte{}, "text/plain; charset=utf-8")
		return
	}

	u.GinContext.Header("Cache-Control", "public, max-age=300")
	u.SuccessRawBytes(http.StatusOK, []byte{}, "text/calendar; charset=utf-8")
}

func GetPublicCalendarICS(c *gin.Context) {
	u := util.GetUtil(c)
	ip := util.DetermineClientAddress(c).String()
	u.Logger.Infof("public calendar ICS ip=%s ua=%s", ip, sanitizeUserAgent(c.Request.User-Agent()))

	enabled, tr := u.Tx.Queries().GetPublicCalendarEnabled()
	if tr != nil {
		u.Error(tr)
		return
	}
	if !enabled {
		u.SuccessRawBytes(http.StatusNotFound, []byte{}, "text/plain; charset=utf-8")
		return
	}

	now := time.Now()
	from := now.AddDate(0, -6, 0)
	to := now.AddDate(1, 0, 0)

	events, tr := u.Tx.Queries().GetAllEventsAllCalendars(from, to, u.Context, u.Config)
	if tr != nil {
		u.Error(tr)
		return
	}

	body, err := buildMinimalPublicICS(events, now)
	if err != nil {
		u.Error(errors.New().Status(http.StatusInternalServerError).
			AddErr(errors.LvlDebug, err).
			Append(errors.LvlPlain, "Could not build calendar"))
		return
	}

	u.GinContext.Header("Cache-Control", "public, max-age=300")
	u.SuccessRawBytes(http.StatusOK, body, "text/calendar; charset=utf-8")
}

func sanitizeUserAgent(s string) string {
	s = strings.TrimSpace(s)
	if len(s) > 512 {
		return s[:512]
	}
	return s
}

func buildMinimalPublicICS(events []types.Event, stamp time.Time) ([]byte, error) {
	cal := ics.NewCalendarFor("Luna")
	cal.SetMethod(ics.MethodPublish)

	for _, ev := range events {
		if ev.GetName() == "" {
			continue
		}
		uid := fmt.Sprintf("%s@luna-public", ev.GetId().String())
		ve := cal.AddEvent(uid)
		ve.SetSummary(ev.GetName())
		ve.SetDtStampTime(stamp)
		ve.SetCreatedTime(stamp)

		date := ev.GetDate()
		if date.AllDay() {
			ve.SetAllDayStartAt(*date.Start())
			ve.SetAllDayEndAt(*date.End())
		} else {
			ve.SetStartAt(*date.Start())
			ve.SetEndAt(*date.End())
		}
	}

	return []byte(cal.Serialize()), nil
}
