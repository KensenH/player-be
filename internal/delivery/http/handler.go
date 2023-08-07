package http

import (
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	log "github.com/sirupsen/logrus"

	resp "player-be/internal/entity/response"

	"player-be/internal/pkg/valmsg"
)

// set routes
func (s *Server) handler() {
	s.e.HTTPErrorHandler = errorHandler

	s.e.Use(logger())

	s.e.Pre(s.PlayerHandler.JwtMiddleware)

	s.e.GET("/", defaultRoute)

	playerBe := s.e.Group("/player-be")
	apiV1 := playerBe.Group("/api/v1")

	//player
	playerV1 := apiV1.Group("/player")
	playerV1.POST("/signup", s.PlayerHandler.SignUp)
	playerV1.POST("/signin", s.PlayerHandler.SignIn)
	playerV1.GET("/signout", s.PlayerHandler.SignOut)
	playerV1.GET("/test", defaultRoute)

}

// default route
func defaultRoute(c echo.Context) error {
	return c.JSON(200, resp.Response{
		Data:   "PLAYER_BE OK",
		Errors: []resp.Error{},
	})
}

// error handler
func errorHandler(err error, c echo.Context) {
	var ve validator.ValidationErrors

	report, ok := err.(*echo.HTTPError)
	if !ok {
		report = echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	if errors.As(err, &ve) {
		for _, fe := range ve {
			report.Message = valmsg.MsgForTag(fe)
			break
		}
	} else {
		report.Message = fmt.Sprintf("%s", err.Error())
	}

	c.Logger().Error(report)
	c.JSON(report.Code, report)
}

// logger
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
