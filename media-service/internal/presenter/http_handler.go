package presenter

import (
	"net/http"
	"strconv"

	"github.com/HatefBarari/microblog-media/internal/usecase"
	"github.com/HatefBarari/microblog-shared/pkg/httputil"
	"github.com/labstack/echo/v4"
)

type HTTPHandler struct {
	mediaUC *usecase.MediaUseCase
}

func NewHTTPHandler(mediaUC *usecase.MediaUseCase) *HTTPHandler {
	return &HTTPHandler{
		mediaUC: mediaUC,
	}
}

// Upload media file
func (h *HTTPHandler) Upload(c echo.Context) error {
	// Check user role - only admin, manager, or author can upload
	role := c.Get("role").(string)
	if role != "admin" && role != "manager" && role != "author" {
		return c.JSON(http.StatusForbidden, httputil.NewError(403, "insufficient permissions"))
	}

	userID := c.Get("userID").(string)
	
	// Get uploaded file
	file, err := c.FormFile("file")
	if err != nil {
		return c.JSON(http.StatusBadRequest, httputil.NewError(400, "no file uploaded"))
	}

	// Upload file
	resp, err := h.mediaUC.Upload(c.Request().Context(), userID, file)
	if err != nil {
		return c.JSON(http.StatusBadRequest, httputil.NewError(400, err.Error()))
	}

	return c.JSON(http.StatusCreated, httputil.OK(resp))
}

// Get media by ID
func (h *HTTPHandler) GetByID(c echo.Context) error {
	id := c.Param("id")
	
	resp, err := h.mediaUC.GetByID(c.Request().Context(), id)
	if err != nil {
		if err.Error() == "media not found" {
			return c.JSON(http.StatusNotFound, httputil.NewError(404, err.Error()))
		}
		return c.JSON(http.StatusInternalServerError, httputil.NewError(500, err.Error()))
	}

	return c.JSON(http.StatusOK, httputil.OK(resp))
}

// List user's media
func (h *HTTPHandler) List(c echo.Context) error {
	userID := c.Get("userID").(string)
	
	// Parse query parameters
	page, _ := strconv.Atoi(c.QueryParam("page"))
	if page <= 0 {
		page = 1
	}
	
	pageSize, _ := strconv.Atoi(c.QueryParam("page_size"))
	if pageSize <= 0 {
		pageSize = 20
	}
	if pageSize > 100 {
		pageSize = 100
	}

	filter := usecase.ListFilter{
		UploaderID: userID,
		Page:       page,
		PageSize:   pageSize,
	}

	resp, err := h.mediaUC.List(c.Request().Context(), filter)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, httputil.NewError(500, err.Error()))
	}

	return c.JSON(http.StatusOK, httputil.OK(resp))
}

// Delete media
func (h *HTTPHandler) Delete(c echo.Context) error {
	userID := c.Get("userID").(string)
	id := c.Param("id")
	
	err := h.mediaUC.Delete(c.Request().Context(), userID, id)
	if err != nil {
		if err.Error() == "media not found" {
			return c.JSON(http.StatusNotFound, httputil.NewError(404, err.Error()))
		}
		if err.Error() == "forbidden" {
			return c.JSON(http.StatusForbidden, httputil.NewError(403, err.Error()))
		}
		return c.JSON(http.StatusBadRequest, httputil.NewError(400, err.Error()))
	}

	return c.NoContent(http.StatusNoContent)
}

// Serve media file
func (h *HTTPHandler) Serve(c echo.Context) error {
	filename := c.Param("filename")
	
	// Get media info from database first
	media, err := h.mediaUC.GetByID(c.Request().Context(), filename)
	if err != nil {
		return c.JSON(http.StatusNotFound, httputil.NewError(404, "media not found"))
	}

	// For now, redirect to the URL
	// In production, you might want to serve the file directly with proper headers
	return c.Redirect(http.StatusMovedPermanently, media.URL)
}
