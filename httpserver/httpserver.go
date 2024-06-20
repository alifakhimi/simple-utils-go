package httpserver

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	echoSwagger "github.com/swaggo/echo-swagger"

	"github.com/alifakhimi/simple-service-go/utils"
)

type (
	// HttpServers is a collection of Http Server
	HttpServers map[string]*HttpServer

	HttpServer struct {
		HttpServerConfig
		// echo is an instance of echo.labstack.com
		echo        *echo.Echo
		prefixGroup *echo.Group
	}

	HttpServerLogLevel uint8

	HttpServerConfig struct {
		// Address optionally specifies the TCP address for the server to listen on,
		// in the form "host:port". If empty, ":http" (port 80) is used.
		// The service names are defined in RFC 6335 and assigned by IANA.
		// See net.Dial for details of the address format.
		Address string `json:"address,omitempty"`
		// Prefix is an api path
		// like /api/v1
		Prefix string `json:"prefix,omitempty"`
		// Debug set develop log level
		Debug bool `json:"debug,omitempty"`
		// LogLevel
		LogLevel HttpServerLogLevel `json:"log_level,omitempty"`
	}
)

const (
	DEBUG HttpServerLogLevel = iota + 1
	INFO
	WARN
	ERROR
	OFF
)

func (h *HttpServer) Echo() *echo.Echo {
	return h.echo
}

func (h *HttpServer) PrefixGroup() *echo.Group {
	return h.prefixGroup
}

func (h *HttpServer) newEcho() (err error) {
	// Create echo instance
	h.echo = echo.New()
	h.echo.HideBanner = true
	h.echo.Logger.SetLevel(log.Lvl(h.LogLevel))

	// Create service middleware logger/recover related to Debug
	if h.Debug {
		h.echo.Use(middleware.Logger())
	} else {
		h.echo.Use(middleware.Recover())
	}

	h.prefixGroup = h.echo.Group(h.Prefix)

	// API Doc
	h.prefixGroup.GET("/swagger/*", echoSwagger.WrapHandler)

	// Add default routes
	h.prefixGroup.Any("/healthinfo", func(ctx echo.Context) error {
		return utils.Reply(
			ctx,
			http.StatusOK,
			nil,
			map[string]any{
				"server": map[string]any{
					"status": "running",
				},
			},
			nil,
		)
	})

	return nil
}

func (h *HttpServer) UnmarshalJSON(b []byte) error {
	var tmp = &HttpServer{}
	if err := json.Unmarshal(b, &tmp.HttpServerConfig); err != nil {
		return err
	}

	tmp.newEcho()

	*h = *tmp

	return nil
}

func (h *HttpServer) Run() error {
	if h.echo == nil {
		if err := h.newEcho(); err != nil {
			return err
		}
	}

	if err := h.echo.Start(h.Address); err != http.ErrServerClosed {
		return err
	}

	return nil
}

func (hs HttpServers) RunAll() error {
	var err error

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	for k, h := range hs {
		go func(k string, h *HttpServer) {
			fmt.Printf("%s service start \n", k)
			defer fmt.Printf("%s service stopped \n", k)

			if e := h.Run(); e != nil {
				err = errors.Join(e)
			}
		}(k, h)
	}

	// Wait for interrupt signal to gracefully shutdown the server with a timeout of 10 seconds.
	<-ctx.Done()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	for k, h := range hs {
		fmt.Printf("%s service is shutting down \n", k)
		if e := h.echo.Shutdown(ctx); e != nil {
			err = errors.Join(err, e)
		}
		fmt.Printf("%s service gracefully shutdown \n", k)
	}

	return err
}
