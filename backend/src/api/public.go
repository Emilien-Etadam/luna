package api

import (
	middleware "luna-backend/api/internal"
	"luna-backend/api/internal/handlers"

	"github.com/gin-gonic/gin"
)

func attachPublicRoutes(endpoints *gin.RouterGroup) {
	public := endpoints.Group("/public", middleware.PublicMinuteRateLimit())
	public.GET("/all.ics", handlers.GetPublicCalendarICS)
	public.HEAD("/all.ics", handlers.HeadPublicCalendarICS)
	public.GET("/sources", handlers.GetPublicSources)
	public.GET("/sources/:sourceId/calendars", handlers.GetPublicSourceCalendars)
	public.GET("/calendars/:calendarId/events", handlers.GetPublicCalendarEvents)
}
