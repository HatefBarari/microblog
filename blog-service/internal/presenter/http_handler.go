package presenter

import (
	"net/http"

	"github.com/HatefBarari/microblog-blog/internal/domain"
	"github.com/HatefBarari/microblog-blog/internal/usecase"
	"github.com/HatefBarari/microblog-shared/pkg/httputil"
	"github.com/labstack/echo/v4"
)

type HTTPHandler struct {
	articleUC  *usecase.ArticleUseCase
	categoryUC *usecase.CategoryUseCase
	commentUC  *usecase.CommentUseCase
	ratingUC   *usecase.RatingUseCase
}

func NewHTTPHandler(articleUC *usecase.ArticleUseCase, categoryUC *usecase.CategoryUseCase, commentUC *usecase.CommentUseCase, ratingUC *usecase.RatingUseCase) *HTTPHandler {
	return &HTTPHandler{
		articleUC:  articleUC,
		categoryUC: categoryUC,
		commentUC:  commentUC,
		ratingUC:   ratingUC,
	}
}

// ---------- Article ----------
func (h *HTTPHandler) CreateArticle(c echo.Context) error {
	userID := c.Get("userID").(string)
	var req usecase.CreateArticleRequest
	if err := httputil.BindAndValidate(c, &req); err != nil {
		return c.JSON(http.StatusBadRequest, httputil.NewError(400, err.Error()))
	}
	resp, err := h.articleUC.Create(c.Request().Context(), userID, req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, httputil.NewError(400, err.Error()))
	}
	return c.JSON(http.StatusCreated, httputil.OK(resp))
}

func (h *HTTPHandler) GetArticleBySlug(c echo.Context) error {
	slug := c.Param("slug")
	resp, err := h.articleUC.GetBySlug(c.Request().Context(), slug)
	if err != nil {
		return c.JSON(http.StatusNotFound, httputil.NewError(404, err.Error()))
	}
	return c.JSON(http.StatusOK, httputil.OK(resp))
}

func (h *HTTPHandler) ListArticles(c echo.Context) error {
	// TODO: parse query filter, page, pageSize
	filter := domain.ListFilter{Page: 1, PageSize: 10}
	list, total, err := h.articleUC.List(c.Request().Context(), filter)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, httputil.NewError(500, err.Error()))
	}
	return c.JSON(http.StatusOK, httputil.OK(map[string]interface{}{
		"list":  list,
		"total": total,
	}))
}

func (h *HTTPHandler) UpdateArticle(c echo.Context) error {
	userID := c.Get("userID").(string)
	id := c.Param("id")
	var req usecase.CreateArticleRequest
	if err := httputil.BindAndValidate(c, &req); err != nil {
		return c.JSON(http.StatusBadRequest, httputil.NewError(400, err.Error()))
	}
	resp, err := h.articleUC.Update(c.Request().Context(), userID, id, req)
	if err != nil {
		if err.Error() == "article not found" {
			return c.JSON(http.StatusNotFound, httputil.NewError(404, err.Error()))
		}
		return c.JSON(http.StatusBadRequest, httputil.NewError(400, err.Error()))
	}
	return c.JSON(http.StatusOK, httputil.OK(resp))
}

func (h *HTTPHandler) DeleteArticle(c echo.Context) error {
	userID := c.Get("userID").(string)
	id := c.Param("id")
	if err := h.articleUC.Delete(c.Request().Context(), userID, id); err != nil {
		if err.Error() == "article not found" {
			return c.JSON(http.StatusNotFound, httputil.NewError(404, err.Error()))
		}
		return c.JSON(http.StatusBadRequest, httputil.NewError(400, err.Error()))
	}
	return c.NoContent(http.StatusNoContent)
}

// ---------- Category ----------
func (h *HTTPHandler) ListCategoryTree(c echo.Context) error {
	tree, err := h.categoryUC.ListTree(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, httputil.NewError(500, err.Error()))
	}
	return c.JSON(http.StatusOK, httputil.OK(tree))
}

func (h *HTTPHandler) CreateCategory(c echo.Context) error {
	role := c.Get("role").(string)
	if role != "manager" && role != "admin" {
		return c.JSON(http.StatusForbidden, httputil.NewError(403, "forbidden"))
	}
	var req struct {
		Name     string `json:"name" validate:"required,min=2"`
		ParentID string `json:"parent_id,omitempty"`
	}
	if err := httputil.BindAndValidate(c, &req); err != nil {
		return c.JSON(http.StatusBadRequest, httputil.NewError(400, err.Error()))
	}
	resp, err := h.categoryUC.Create(c.Request().Context(), req.Name, req.ParentID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, httputil.NewError(400, err.Error()))
	}
	return c.JSON(http.StatusCreated, httputil.OK(resp))
}

// ---------- Comment ----------
func (h *HTTPHandler) ListComments(c echo.Context) error {
	articleID := c.Param("id")
	status := domain.CommentApproved // فقط تاییدشده‌ها برای مهمان
	list, err := h.commentUC.ListByArticle(c.Request().Context(), articleID, status)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, httputil.NewError(500, err.Error()))
	}
	return c.JSON(http.StatusOK, httputil.OK(list))
}

func (h *HTTPHandler) CreateComment(c echo.Context) error {
	userID := c.Get("userID").(string)
	articleID := c.Param("id")
	var req usecase.CreateCommentRequest
	if err := httputil.BindAndValidate(c, &req); err != nil {
		return c.JSON(http.StatusBadRequest, httputil.NewError(400, err.Error()))
	}
	resp, err := h.commentUC.Create(c.Request().Context(), userID, articleID, req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, httputil.NewError(400, err.Error()))
	}
	return c.JSON(http.StatusCreated, httputil.OK(resp))
}

func (h *HTTPHandler) UpdateCommentStatus(c echo.Context) error {
	role := c.Get("role").(string)
	if role != "manager" && role != "admin" {
		return c.JSON(http.StatusForbidden, httputil.NewError(403, "forbidden"))
	}
	id := c.Param("id")
	var req struct {
		Status domain.CommentStatus `json:"status" validate:"required,oneof=approved rejected"`
	}
	if err := httputil.BindAndValidate(c, &req); err != nil {
		return c.JSON(http.StatusBadRequest, httputil.NewError(400, err.Error()))
	}
	if err := h.commentUC.UpdateStatus(c.Request().Context(), id, req.Status); err != nil {
		return c.JSON(http.StatusBadRequest, httputil.NewError(400, err.Error()))
	}
	return c.JSON(http.StatusOK, httputil.OK(map[string]string{"message": "updated"}))
}

// ---------- Rating ----------
func (h *HTTPHandler) RateArticle(c echo.Context) error {
	userID := c.Get("userID").(string)
	articleID := c.Param("id")
	var req usecase.RatingRequest
	if err := httputil.BindAndValidate(c, &req); err != nil {
		return c.JSON(http.StatusBadRequest, httputil.NewError(400, err.Error()))
	}
	if err := h.ratingUC.RateArticle(c.Request().Context(), userID, articleID, req.Stars); err != nil {
		return c.JSON(http.StatusBadRequest, httputil.NewError(400, err.Error()))
	}
	return c.JSON(http.StatusCreated, httputil.OK(map[string]string{"message": "rated"}))
}

func (h *HTTPHandler) DeleteRating(c echo.Context) error {
	userID := c.Get("userID").(string)
	articleID := c.Param("id")
	if err := h.ratingUC.Delete(c.Request().Context(), userID, articleID, "article"); err != nil {
		return c.JSON(http.StatusBadRequest, httputil.NewError(400, err.Error()))
	}
	return c.NoContent(http.StatusNoContent)
}