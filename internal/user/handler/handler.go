package handler

import (
	"context"
	"fmt"
	"net/http"
	"reflect"

	"github.com/core-go/core"
	s "github.com/core-go/search"
	"github.com/labstack/echo/v4"

	"go-service/internal/user/model"
	"go-service/internal/user/service"
)

type UserHandler struct {
	service     service.UserService
	Validate    func(context.Context, interface{}) ([]core.ErrorMessage, error)
	Error       func(context.Context, string, ...map[string]interface{})
	Map         map[string]int
	ParamIndex  map[string]int
	FilterIndex int
}

func NewUserHandler(service service.UserService, logError func(context.Context, string, ...map[string]interface{}), validate func(context.Context, interface{}) ([]core.ErrorMessage, error)) *UserHandler {
	userType := reflect.TypeOf(model.User{})
	_, jsonMap, _ := core.BuildMapField(userType)
	paramIndex, filterIndex := s.BuildAttributes(reflect.TypeOf(model.UserFilter{}))
	return &UserHandler{service: service, Validate: validate, Map: jsonMap, Error: logError, ParamIndex: paramIndex, FilterIndex: filterIndex}
}

func (h *UserHandler) All(c echo.Context) error {
	users, err := h.service.All(c.Request().Context())
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, users)
}

func (h *UserHandler) Load(c echo.Context) error {
	id := c.Param("id")
	if len(id) == 0 {
		return c.String(http.StatusBadRequest, "Id cannot be empty")
	}

	user, err := h.service.Load(c.Request().Context(), id)
	if err != nil {
		h.Error(c.Request().Context(), fmt.Sprintf("Error to get user %s: %s", id, err.Error()))
		return c.String(http.StatusInternalServerError, err.Error())
	}
	if user == nil {
		return c.JSON(http.StatusNotFound, user)
	}
	return c.JSON(http.StatusOK, user)
}

func (h *UserHandler) Create(c echo.Context) error {
	var user model.User
	er1 := c.Bind(&user)

	defer c.Request().Body.Close()
	if er1 != nil {
		return c.String(http.StatusInternalServerError, er1.Error())
	}

	errors, er2 := h.Validate(c.Request().Context(), &user)
	if er2 != nil {
		h.Error(c.Request().Context(), er2.Error(), core.MakeMap(user))
		return c.String(http.StatusInternalServerError, core.InternalServerError)
	}
	if len(errors) > 0 {
		return c.JSON(http.StatusUnprocessableEntity, errors)
	}

	res, er3 := h.service.Create(c.Request().Context(), &user)
	if er3 != nil {
		return c.String(http.StatusInternalServerError, er3.Error())
	}
	if res > 0 {
		return c.JSON(http.StatusCreated, user)
	} else {
		return c.JSON(http.StatusConflict, res)
	}
}

func (h *UserHandler) Update(c echo.Context) error {
	var user model.User
	er1 := c.Bind(&user)
	defer c.Request().Body.Close()

	if er1 != nil {
		return c.String(http.StatusInternalServerError, er1.Error())
	}

	id := c.Param("id")
	if len(id) == 0 {
		return c.String(http.StatusBadRequest, "Id cannot be empty")
	}

	if len(user.Id) == 0 {
		user.Id = id
	} else if id != user.Id {
		return c.String(http.StatusBadRequest, "Id not match")
	}

	errors, er2 := h.Validate(c.Request().Context(), &user)
	if er2 != nil {
		h.Error(c.Request().Context(), er2.Error(), core.MakeMap(user))
		return c.String(http.StatusInternalServerError, core.InternalServerError)
	}
	if len(errors) > 0 {
		return c.JSON(http.StatusUnprocessableEntity, errors)
	}

	res, er3 := h.service.Update(c.Request().Context(), &user)
	if er3 != nil {
		return c.String(http.StatusInternalServerError, er3.Error())
	}
	if res > 0 {
		return c.JSON(http.StatusOK, user)
	} else if res == 0 {
		return c.JSON(http.StatusNotFound, res)
	} else {
		return c.JSON(http.StatusConflict, res)
	}
}

func (h *UserHandler) Patch(c echo.Context) error {
	id := c.Param("id")
	if len(id) == 0 {
		return c.String(http.StatusBadRequest, "Id cannot be empty")
	}

	r := c.Request()
	var user model.User
	userType := reflect.TypeOf(user)
	_, jsonMap, _ := core.BuildMapField(userType)
	body, er0 := core.BuildMapAndStruct(r, &user)
	if er0 != nil {
		return c.String(http.StatusInternalServerError, er0.Error())
	}
	if len(user.Id) == 0 {
		user.Id = id
	} else if id != user.Id {
		return c.String(http.StatusBadRequest, "Id not match")
	}
	json, er1 := core.BodyToJsonMap(r, user, body, []string{"id"}, jsonMap)
	if er1 != nil {
		return c.String(http.StatusInternalServerError, er1.Error())
	}

	errors, er2 := h.Validate(c.Request().Context(), &user)
	if er2 != nil {
		h.Error(c.Request().Context(), er2.Error(), core.MakeMap(user))
		return c.String(http.StatusInternalServerError, core.InternalServerError)
	}
	if len(errors) > 0 {
		return c.JSON(http.StatusUnprocessableEntity, errors)
	}

	res, er3 := h.service.Patch(r.Context(), json)
	if er3 != nil {
		return c.String(http.StatusInternalServerError, er3.Error())
	}
	if res > 0 {
		return c.JSON(http.StatusOK, json)
	} else if res == 0 {
		return c.JSON(http.StatusNotFound, res)
	} else {
		return c.JSON(http.StatusConflict, res)
	}
}

func (h *UserHandler) Delete(c echo.Context) error {
	id := c.Param("id")
	if len(id) == 0 {
		return c.String(http.StatusBadRequest, "Id cannot be empty")
	}

	res, err := h.service.Delete(c.Request().Context(), id)
	if err != nil {
		h.Error(c.Request().Context(), fmt.Sprintf("Error to delete user %s: %s", id, err.Error()))
		return c.String(http.StatusInternalServerError, err.Error())
	}
	if res > 0 {
		return c.JSON(http.StatusOK, res)
	} else {
		return c.JSON(http.StatusNotFound, res)
	}
}

func (h *UserHandler) Search(c echo.Context) error {
	filter := model.UserFilter{Filter: &s.Filter{}}
	err := s.Decode(c.Request(), &filter, h.ParamIndex, h.FilterIndex)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	offset := s.GetOffset(filter.Limit, filter.Page)
	users, total, err := h.service.Search(c.Request().Context(), &filter, filter.Limit, offset)
	if err != nil {
		h.Error(c.Request().Context(), fmt.Sprintf("Error to search user %v: %s", filter, err.Error()))
		return c.String(http.StatusInternalServerError, core.InternalServerError)
	}
	return c.JSON(http.StatusOK, &s.Result{List: &users, Total: total})
}
