package presenter

import (
	"net/http"

	"github.com/HatefBarari/microblog-auth/internal/usecase"
	"github.com/HatefBarari/microblog-shared/pkg/httputil"
	"github.com/labstack/echo/v4"
)

type HTTPHandler struct {
	uc      *usecase.UserUseCase
	emailUC *usecase.EmailUseCase
}

func NewHTTPHandler(uc *usecase.UserUseCase, emailUC *usecase.EmailUseCase) *HTTPHandler {
	return &HTTPHandler{
		uc:      uc,
		emailUC: emailUC,
	}
}

func (h *HTTPHandler) Register(c echo.Context) error {
	var req usecase.RegisterRequest
	if err := httputil.BindAndValidate(c, &req); err != nil {
		return c.JSON(http.StatusBadRequest, httputil.NewError(400, err.Error()))
	}
	resp, err := h.uc.Register(c.Request().Context(), req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, httputil.NewError(400, err.Error()))
	}
	return c.JSON(http.StatusCreated, httputil.OK(resp))
}

func (h *HTTPHandler) Login(c echo.Context) error {
	var req usecase.LoginRequest
	if err := httputil.BindAndValidate(c, &req); err != nil {
		return c.JSON(http.StatusBadRequest, httputil.NewError(400, err.Error()))
	}
	resp, err := h.uc.Login(c.Request().Context(), req)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, httputil.NewError(401, err.Error()))
	}
	return c.JSON(http.StatusOK, httputil.OK(resp))
}

func (h *HTTPHandler) Refresh(c echo.Context) error {
	var req usecase.RefreshRequest
	if err := httputil.BindAndValidate(c, &req); err != nil {
		return c.JSON(http.StatusBadRequest, httputil.NewError(400, err.Error()))
	}
	resp, err := h.uc.Refresh(c.Request().Context(), req.RefreshToken)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, httputil.NewError(401, err.Error()))
	}
	return c.JSON(http.StatusOK, httputil.OK(resp))
}

// VerifyEmail verifies user email with token
func (h *HTTPHandler) VerifyEmail(c echo.Context) error {
	token := c.QueryParam("token")
	if token == "" {
		return c.JSON(http.StatusBadRequest, httputil.NewError(400, "verification token is required"))
	}

	err := h.emailUC.VerifyEmail(c.Request().Context(), token)
	if err != nil {
		return c.JSON(http.StatusBadRequest, httputil.NewError(400, err.Error()))
	}

	return c.JSON(http.StatusOK, httputil.OK(map[string]string{
		"message": "email verified successfully",
	}))
}

// ResendVerificationEmail resends verification email
func (h *HTTPHandler) ResendVerificationEmail(c echo.Context) error {
	userID := c.Get("userID").(string)
	email := c.Get("email").(string)

	err := h.emailUC.SendVerificationEmail(c.Request().Context(), userID, email)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, httputil.NewError(500, err.Error()))
	}

	return c.JSON(http.StatusOK, httputil.OK(map[string]string{
		"message": "verification email sent",
	}))
}

// SendPasswordResetEmail sends password reset email
func (h *HTTPHandler) SendPasswordResetEmail(c echo.Context) error {
	var req struct {
		Email string `json:"email" validate:"required,email"`
	}
	if err := httputil.BindAndValidate(c, &req); err != nil {
		return c.JSON(http.StatusBadRequest, httputil.NewError(400, err.Error()))
	}

	err := h.emailUC.SendPasswordResetEmail(c.Request().Context(), req.Email)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, httputil.NewError(500, err.Error()))
	}

	return c.JSON(http.StatusOK, httputil.OK(map[string]string{
		"message": "password reset email sent",
	}))
}

// ResetPassword resets user password with token
func (h *HTTPHandler) ResetPassword(c echo.Context) error {
	var req struct {
		Token       string `json:"token" validate:"required"`
		NewPassword string `json:"new_password" validate:"required,min=6"`
	}
	if err := httputil.BindAndValidate(c, &req); err != nil {
		return c.JSON(http.StatusBadRequest, httputil.NewError(400, err.Error()))
	}

	err := h.emailUC.ResetPassword(c.Request().Context(), req.Token, req.NewPassword)
	if err != nil {
		return c.JSON(http.StatusBadRequest, httputil.NewError(400, err.Error()))
	}

	return c.JSON(http.StatusOK, httputil.OK(map[string]string{
		"message": "password reset successfully",
	}))
}