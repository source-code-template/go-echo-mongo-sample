package user

import (
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/core-go/core"
	v "github.com/core-go/core/v10"
	"github.com/labstack/echo/v4"

	"go-service/internal/user/handler"
	"go-service/internal/user/repository/adapter"
	"go-service/internal/user/repository/query"
	"go-service/internal/user/service"
)

type UserTransport interface {
	All(echo.Context) error
	Load(echo.Context) error
	Create(echo.Context) error
	Update(echo.Context) error
	Patch(echo.Context) error
	Delete(echo.Context) error
	Search(echo.Context) error
}

func NewUserHandler(db *mongo.Database, logError core.Log, action *core.ActionConfig) (UserTransport, error) {
	validator, err := v.NewValidator()
	if err != nil {
		return nil, err
	}

	userRepository := adapter.NewUserAdapter(db, query.BuildQuery)
	userService := service.NewUserService(userRepository)
	userHandler := handler.NewUserHandler(userService, logError, validator.Validate, action)
	return userHandler, nil
}
