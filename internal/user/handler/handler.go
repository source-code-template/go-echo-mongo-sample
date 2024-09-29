package handler

import (
	"fmt"
	"net/http"
	"reflect"

	"github.com/core-go/core"
	e "github.com/core-go/core/handler/echo"
	"github.com/core-go/search"
	"github.com/labstack/echo/v4"

	"go-service/internal/user/model"
	"go-service/internal/user/service"
)

type UserHandler struct {
	service  service.UserService
	Validate core.Validate
	*core.Attributes
	*search.Parameters
}

func NewUserHandler(service service.UserService, logError core.Log, validate core.Validate, action *core.ActionConfig) *UserHandler {
	userType := reflect.TypeOf(model.User{})
	parameters := search.CreateParameters(reflect.TypeOf(model.UserFilter{}), userType)
	attributes := core.CreateAttributes(userType, logError, action)
	return &UserHandler{service: service, Validate: validate, Attributes: attributes, Parameters: parameters}
}

func (h *UserHandler) All(c echo.Context) error {
	users, err := h.service.All(c.Request().Context())
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, users)
}

func (h *UserHandler) Load(c echo.Context) error {
	id, err := e.GetRequiredString(c)
	if err != nil {
		return err
	}
	user, err := h.service.Load(c.Request().Context(), id)
	if err != nil {
		h.Error(c.Request().Context(), fmt.Sprintf("Error to get user '%s': %s", id, err.Error()))
		return c.String(http.StatusInternalServerError, core.InternalServerError)
	}
	return c.JSON(core.IsFound(user), user)
}

func (h *UserHandler) Create(c echo.Context) error {
	var user model.User
	er1 := e.Decode(c, &user)
	if er1 != nil {
		return er1
	}

	errors, er2 := h.Validate(c.Request().Context(), &user)
	if e.HasError(c, errors, er2, h.Error, user, h.Log, h.Resource, h.Action.Create) {
		return er2
	}

	res, er3 := h.service.Create(c.Request().Context(), &user)
	return e.AfterCreated(c, &user, res, er3, h.Error)
}

func (h *UserHandler) Update(c echo.Context) error {
	var user model.User
	er1 := e.DecodeAndCheckId(c, &user, h.Keys, h.Indexes)
	if er1 != nil {
		return er1
	}

	errors, er2 := h.Validate(c.Request().Context(), &user)
	if e.HasError(c, errors, er2, h.Error, user, h.Log, h.Resource, h.Action.Update) {
		return er2
	}

	res, er3 := h.service.Update(c.Request().Context(), &user)
	if er3 != nil {
		h.Error(c.Request().Context(), er2.Error(), core.MakeMap(user))
		return c.String(http.StatusInternalServerError, core.InternalServerError)
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
	var user model.User
	jsonUser, er1 := e.BuildMapAndCheckId(c, &user, h.Keys, h.Indexes)
	if er1 != nil {
		return er1
	}

	errors, er2 := h.Validate(c.Request().Context(), &user)
	if e.HasError(c, errors, er2, h.Error, user, h.Log, h.Resource, h.Action.Patch) {
		return er2
	}

	res, er3 := h.service.Patch(c.Request().Context(), jsonUser)
	return e.AfterSaved(c, jsonUser, res, er3, h.Error)
}

func (h *UserHandler) Delete(c echo.Context) error {
	id, err := e.GetRequiredString(c)
	if err != nil {
		return err
	}

	res, err := h.service.Delete(c.Request().Context(), id)
	if err != nil {
		h.Error(c.Request().Context(), fmt.Sprintf("Error to delete user '%s': %s", id, err.Error()))
		return c.String(http.StatusInternalServerError, err.Error())
	}
	if res > 0 {
		return c.JSON(http.StatusOK, res)
	} else {
		return c.JSON(http.StatusNotFound, res)
	}
}

func (h *UserHandler) Search(c echo.Context) error {
	filter := model.UserFilter{Filter: &search.Filter{}}
	err := search.Decode(c.Request(), &filter, h.ParamIndex, h.FilterIndex)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	offset := search.GetOffset(filter.Limit, filter.Page)
	users, total, err := h.service.Search(c.Request().Context(), &filter, filter.Limit, offset)
	if err != nil {
		h.Error(c.Request().Context(), fmt.Sprintf("Error to search user %v: %s", filter, err.Error()))
		return c.String(http.StatusInternalServerError, core.InternalServerError)
	}
	return c.JSON(http.StatusOK, &search.Result{List: &users, Total: total})
}
