package main

import (
	"context"

	"github.com/core-go/config"
	"github.com/core-go/core/header/echo"
	"github.com/core-go/core/random"
	svr "github.com/core-go/core/server"
	"github.com/core-go/core/strings"
	mid "github.com/core-go/log/echo"
	"github.com/core-go/log/zap"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"go-service/internal/app"
)

func main() {
	var cfg app.Config
	err := config.Load(&cfg, "configs/config")
	if err != nil {
		panic(err)
	}

	e := echo.New()
	log.Initialize(cfg.Log)
	logger := mid.NewMaskLogger(cfg.MiddleWare.Request, Mask, Mask)
	echoLogger := mid.NewEchoLogger(cfg.MiddleWare, log.InfoFields, logger, MaskLog)

	e.Use(echoLogger.BuildContextWithMask)
	e.Use(echoLogger.Logger)
	headerHandler := header.NewHeaderHandler(cfg.Response, GenerateId)
	e.Use(headerHandler.HandleHeader)
	e.Use(middleware.Recover())

	err = app.Route(context.Background(), e, cfg)
	if err != nil {
		panic(err)
	}
	e.Logger.Fatal(e.Start(svr.Addr(cfg.Server.Port)))
}

func GenerateId() string {
	return random.Random(16)
}
func MaskLog(name, s string) string {
	if name == "mobileNo" {
		return strings.Mask(s, 2, 2, "x")
	} else {
		return strings.Mask(s, 0, 5, "x")
	}
}
func Mask(obj map[string]interface{}) {
	v, ok := obj["phone"]
	if ok {
		s, ok2 := v.(string)
		if ok2 && len(s) > 3 {
			obj["phone"] = strings.Mask(s, 0, 3, "*")
		}
	}
}
