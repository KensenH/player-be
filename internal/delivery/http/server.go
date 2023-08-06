package http

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"

	log "github.com/sirupsen/logrus"
)

type Option func(*Server)

type Server struct {
	e             echo.Echo
	PlayerHandler PlayerHandler
}

type PlayerHandler interface {
	SignUp(c echo.Context) error
	SignIn(c echo.Context) error
	SignOut(c echo.Context) error
	JwtMiddleware(next echo.HandlerFunc) echo.HandlerFunc
}

// create server
func New(playerHandler PlayerHandler, opts ...Option) *Server {
	server := Server{
		e:             *echo.New(),
		PlayerHandler: playerHandler,
	}

	// set route
	server.handler()

	return &server
}

// start server
func (s *Server) Serve() error {
	var err error
	//start
	go func() {
		if err := s.e.Start(":8080"); err != nil && err != http.ErrServerClosed {
			log.Errorf("error server: %s", err.Error())
		}
	}()

	//graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	log.Infoln("shutting down server")
	if err = s.e.Shutdown(ctx); err != nil {
		return errors.Wrap(err, "error shutting down server")
	}

	log.Infoln("server shutdown gracefully")

	return err
}
