package app

import (
	"context"

	"github.com/labstack/echo/v4"
)

func Route(ctx context.Context, e *echo.Echo, cfg Config) error {
	app, err := NewApp(ctx, cfg)
	if err != nil {
		return err
	}

	e.GET("/health", app.Health.Check)

	userPath := "/users"
	e.GET(userPath, app.User.All)
	e.GET(userPath+"/search", app.User.Search)
	e.POST(userPath+"/search", app.User.Search)
	e.GET(userPath+"/:id", app.User.Load)
	e.POST(userPath, app.User.Create)
	e.PUT(userPath+"/:id", app.User.Update)
	e.PATCH(userPath+"/:id", app.User.Patch)
	e.DELETE(userPath+"/:id", app.User.Delete)

	return nil
}
