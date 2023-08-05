package http

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"

	resp "player-be/internal/entity/response"
)

// set routes
func (s *Server) handler() {
	s.e.Use(logger())

	s.e.GET("/", defaultRoute)

}

// default route
func defaultRoute(c echo.Context) error {
	return c.JSON(200, resp.Response{
		Data:   "PLAYER_BE OK",
		Errors: []resp.Error{},
	})
}

// loggger
func logger() echo.MiddlewareFunc {
	return middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogURI:    true,
		LogStatus: true,
		LogValuesFunc: func(c echo.Context, values middleware.RequestLoggerValues) error {
			log.WithFields(logrus.Fields{
				"URI":    values.URI,
				"status": values.Status,
			}).Info("request")

			return nil
		},
	})
}
