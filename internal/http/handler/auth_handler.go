package handler

import (
	"encoding/json"

	"github.com/mnabielap/starter-kit-restapi-gofasthttp/internal/model"
	"github.com/mnabielap/starter-kit-restapi-gofasthttp/internal/service"
	"github.com/mnabielap/starter-kit-restapi-gofasthttp/pkg/utils"
	"github.com/valyala/fasthttp"
)

type AuthHandler struct {
	authService *service.AuthService
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

// Register godoc
// @Summary      Register a new user
// @Description  Register a new user and return tokens
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        user body model.User true "User Data"
// @Success      201  {object}  map[string]interface{}
// @Failure      400  {object}  utils.Response
// @Router       /auth/register [post]
func (h *AuthHandler) Register(ctx *fasthttp.RequestCtx) {
	var user model.User
	if err := json.Unmarshal(ctx.PostBody(), &user); err != nil {
		utils.WriteError(ctx, fasthttp.StatusBadRequest, "Invalid request body")
		return
	}

	// Validate Struct
	if validationErrors := utils.ValidateStruct(&user); validationErrors != nil {
		utils.WriteJSON(ctx, fasthttp.StatusBadRequest, map[string]interface{}{
			"status": "error",
			"errors": validationErrors,
		})
		return
	}

	user.Role = "user"

	createdUser, tokens, err := h.authService.Register(&user)
	if err != nil {
		utils.WriteError(ctx, fasthttp.StatusBadRequest, err.Error())
		return
	}

	utils.WriteSuccess(ctx, fasthttp.StatusCreated, map[string]interface{}{
		"user":   createdUser.ToResponse(),
		"tokens": tokens,
	})
}

// Login godoc
// @Summary      Login user
// @Description  Login with email and password to get tokens
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        request body LoginRequest true "Login Credentials"
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  utils.Response
// @Failure      401  {object}  utils.Response
// @Router       /auth/login [post]
func (h *AuthHandler) Login(ctx *fasthttp.RequestCtx) {
	var loginData LoginRequest

	if err := json.Unmarshal(ctx.PostBody(), &loginData); err != nil {
		utils.WriteError(ctx, fasthttp.StatusBadRequest, "Invalid request body")
		return
	}

	if validationErrors := utils.ValidateStruct(&loginData); validationErrors != nil {
		utils.WriteJSON(ctx, fasthttp.StatusBadRequest, map[string]interface{}{
			"status": "error",
			"errors": validationErrors,
		})
		return
	}

	user, tokens, err := h.authService.Login(loginData.Email, loginData.Password)
	if err != nil {
		utils.WriteError(ctx, fasthttp.StatusUnauthorized, err.Error())
		return
	}

	utils.WriteSuccess(ctx, fasthttp.StatusOK, map[string]interface{}{
		"user":   user.ToResponse(),
		"tokens": tokens,
	})
}

// Logout godoc
// @Summary      Logout user
// @Description  Invalidate refresh token
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        request body RefreshTokenRequest true "Refresh Token"
// @Success      204
// @Failure      404  {object}  utils.Response
// @Router       /auth/logout [post]
func (h *AuthHandler) Logout(ctx *fasthttp.RequestCtx) {
	var body RefreshTokenRequest
	json.Unmarshal(ctx.PostBody(), &body)

	if err := h.authService.Logout(body.RefreshToken); err != nil {
		utils.WriteError(ctx, fasthttp.StatusNotFound, "Token not found")
		return
	}
	ctx.SetStatusCode(fasthttp.StatusNoContent)
}

// RefreshTokens godoc
// @Summary      Refresh auth tokens
// @Description  Get new access and refresh tokens using a valid refresh token
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        request body RefreshTokenRequest true "Refresh Token"
// @Success      200  {object}  service.AuthTokens
// @Failure      401  {object}  utils.Response
// @Router       /auth/refresh-tokens [post]
func (h *AuthHandler) RefreshTokens(ctx *fasthttp.RequestCtx) {
	var body RefreshTokenRequest
	json.Unmarshal(ctx.PostBody(), &body)

	tokens, err := h.authService.RefreshAuth(body.RefreshToken)
	if err != nil {
		utils.WriteError(ctx, fasthttp.StatusUnauthorized, err.Error())
		return
	}

	utils.WriteSuccess(ctx, fasthttp.StatusOK, tokens)
}

// --- Request Structs for Swagger & Validation ---

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refreshToken"`
}