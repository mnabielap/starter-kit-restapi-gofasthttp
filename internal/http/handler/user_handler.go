package handler

import (
	"encoding/json"
	"strconv"

	"github.com/mnabielap/starter-kit-restapi-gofasthttp/internal/model"
	"github.com/mnabielap/starter-kit-restapi-gofasthttp/internal/service"
	"github.com/mnabielap/starter-kit-restapi-gofasthttp/pkg/utils"
	routing "github.com/qiangxue/fasthttp-routing"
	"github.com/valyala/fasthttp"
)

type UserHandler struct {
	userService *service.UserService
}

func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

// CreateUser godoc
// @Summary      Create a new user (Admin)
// @Description  Create a new user manually. Requires 'admin' role.
// @Tags         Users
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        user body model.User true "User Data"
// @Success      201  {object}  model.UserResponse
// @Failure      400  {object}  utils.Response
// @Failure      403  {object}  utils.Response
// @Router       /users [post]
func (h *UserHandler) CreateUser(ctx *fasthttp.RequestCtx) {
	var user model.User
	if err := json.Unmarshal(ctx.PostBody(), &user); err != nil {
		utils.WriteError(ctx, fasthttp.StatusBadRequest, "Invalid request body")
		return
	}

	if validationErrors := utils.ValidateStruct(&user); validationErrors != nil {
		utils.WriteJSON(ctx, fasthttp.StatusBadRequest, map[string]interface{}{
			"status": "error",
			"errors": validationErrors,
		})
		return
	}

	createdUser, err := h.userService.CreateUser(&user)
	if err != nil {
		utils.WriteError(ctx, fasthttp.StatusBadRequest, err.Error())
		return
	}

	utils.WriteSuccess(ctx, fasthttp.StatusCreated, createdUser.ToResponse())
}

// GetUsers godoc
// @Summary      Get all users
// @Description  Get a paginated list of users. Requires 'admin' role.
// @Tags         Users
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        page  query    int  false  "Page number" default(1)
// @Param        limit query    int  false  "Page limit"  default(10)
// @Success      200   {object} map[string]interface{}
// @Failure      403   {object} utils.Response
// @Router       /users [get]
func (h *UserHandler) GetUsers(ctx *fasthttp.RequestCtx) {
	page, _ := strconv.Atoi(string(ctx.QueryArgs().Peek("page")))
	limit, _ := strconv.Atoi(string(ctx.QueryArgs().Peek("limit")))
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	users, total, err := h.userService.GetUsers(page, limit)
	if err != nil {
		utils.WriteError(ctx, fasthttp.StatusInternalServerError, err.Error())
		return
	}

	utils.WriteSuccess(ctx, fasthttp.StatusOK, map[string]interface{}{
		"results": users,
		"page":    page,
		"limit":   limit,
		"total":   total,
	})
}

// GetUser godoc
// @Summary      Get a user by ID
// @Description  Get details of a specific user. Requires 'admin' role.
// @Tags         Users
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        userId path      int  true  "User ID"
// @Success      200    {object}  model.UserResponse
// @Failure      404    {object}  utils.Response
// @Failure      403    {object}  utils.Response
// @Router       /users/{userId} [get]
func (h *UserHandler) GetUser(c *routing.Context) error {
	id, err := strconv.Atoi(c.Param("userId"))
	if err != nil {
		utils.WriteError(c.RequestCtx, fasthttp.StatusBadRequest, "Invalid user ID")
		return nil
	}

	user, err := h.userService.GetUserByID(uint(id))
	if err != nil {
		utils.WriteError(c.RequestCtx, fasthttp.StatusNotFound, err.Error())
		return nil
	}

	utils.WriteSuccess(c.RequestCtx, fasthttp.StatusOK, user.ToResponse())
	return nil
}

// UpdateUser godoc
// @Summary      Update a user
// @Description  Update user details. Requires 'admin' role.
// @Tags         Users
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        userId path      int  true  "User ID"
// @Param        user   body      map[string]interface{} true "Update Data"
// @Success      200    {object}  model.UserResponse
// @Failure      403    {object}  utils.Response
// @Router       /users/{userId} [patch]
func (h *UserHandler) UpdateUser(c *routing.Context) error {
	id, err := strconv.Atoi(c.Param("userId"))
	if err != nil {
		utils.WriteError(c.RequestCtx, fasthttp.StatusBadRequest, "Invalid user ID")
		return nil
	}

	var updateBody map[string]interface{}
	if err := json.Unmarshal(c.PostBody(), &updateBody); err != nil {
		utils.WriteError(c.RequestCtx, fasthttp.StatusBadRequest, "Invalid request body")
		return nil
	}

	user, err := h.userService.UpdateUser(uint(id), updateBody)
	if err != nil {
		utils.WriteError(c.RequestCtx, fasthttp.StatusBadRequest, err.Error())
		return nil
	}

	utils.WriteSuccess(c.RequestCtx, fasthttp.StatusOK, user.ToResponse())
	return nil
}

// DeleteUser godoc
// @Summary      Delete a user
// @Description  Soft delete a user by ID. Requires 'admin' role.
// @Tags         Users
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        userId path      int  true  "User ID"
// @Success      204
// @Failure      403    {object}  utils.Response
// @Router       /users/{userId} [delete]
func (h *UserHandler) DeleteUser(c *routing.Context) error {
	id, err := strconv.Atoi(c.Param("userId"))
	if err != nil {
		utils.WriteError(c.RequestCtx, fasthttp.StatusBadRequest, "Invalid user ID")
		return nil
	}

	if err := h.userService.DeleteUser(uint(id)); err != nil {
		utils.WriteError(c.RequestCtx, fasthttp.StatusNotFound, err.Error())
		return nil
	}

	c.SetStatusCode(fasthttp.StatusNoContent)
	return nil
}