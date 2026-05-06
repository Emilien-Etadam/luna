package handlers

import (
	"luna-backend/api/internal/util"
	"luna-backend/errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetAdminPublicCalendar(c *gin.Context) {
	u := util.GetUtil(c)

	enabled, tr := u.Tx.Queries().GetPublicCalendarEnabled()
	if tr != nil {
		u.Error(tr)
		return
	}

	u.Success(&gin.H{"enabled": enabled})
}

func PatchAdminPublicCalendar(c *gin.Context) {
	u := util.GetUtil(c)

	var body struct {
		Enabled bool `json:"enabled"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		u.Error(errors.New().Status(http.StatusBadRequest).
			AddErr(errors.LvlDebug, err).
			Append(errors.LvlPlain, "Invalid JSON body"))
		return
	}

	tr := u.Tx.Queries().SetPublicCalendarEnabled(body.Enabled)
	if tr != nil {
		u.Error(tr)
		return
	}

	InvalidatePublicICSCache()

	u.Success(&gin.H{"enabled": body.Enabled})
}
