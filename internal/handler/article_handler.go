package handler

import (
	"errors"
	"go_redbook/internal/pkg/req"
	"go_redbook/internal/pkg/response"
	"go_redbook/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

// ArticleHandler 只处理 HTTP 层逻辑：
// 参数绑定、参数校验、调用 service、组装响应。
type ArticleHandler struct {
	articleService service.ArticleService
}

func NewArticleHandler(articleService service.ArticleService) *ArticleHandler {
	return &ArticleHandler{
		articleService: articleService,
	}
}

func (h *ArticleHandler) RegisterCreateRoutes(group *gin.RouterGroup) {
	articles := group.Group("/articles")
	articles.POST("", h.Create)
	articles.POST("/:id/like", h.LikeYes)
	articles.DELETE("/:id/like", h.LikeNo)
	articles.POST("/:id/favorite", h.FavoriteYes)
	articles.DELETE("/:id/favorite", h.FavoriteNo)
}
func (h *ArticleHandler) RegisterPublicRoutes(group *gin.RouterGroup) {
	articles := group.Group("/articles")
	articles.GET("/list", h.List)
	articles.GET("/:id", h.GetByID)
}
func (h *ArticleHandler) RegisterPrivateRoutes(group *gin.RouterGroup) {
	articles := group.Group("/articles")
	articles.DELETE("/:id", h.DeleteByID)
	articles.PUT("/:id", h.Update)

}

// 创建文章
func (h *ArticleHandler) Create(c *gin.Context) {
	var r req.CreateArticleRequest
	if err := c.ShouldBindJSON(&r); err != nil {
		response.Error(c, http.StatusBadRequest, 40001, "参数错误: "+err.Error())
		return
	}

	currentUserInfo, ok := currentUser(c)
	if !ok {
		return
	}

	article, err := h.articleService.Create(c.Request.Context(), service.CreateArticleInput{
		UserID:    currentUserInfo.UserID,
		Title:     r.Title,
		Content:   r.Content,
		ImageURLs: r.Images,
		Topics:    r.Topics,
		Location:  r.Location,
		IsPublic:  r.IsPublic,
	})
	if err != nil {
		response.Error(c, http.StatusInternalServerError, 50001, "创建文章失败")
		return
	}
	response.Created(c, article)
}

// 文章列表
func (h *ArticleHandler) List(c *gin.Context) {
	page := queryInt(c, "page", 1)
	pageSize := queryInt(c, "page_size", 10)

	result, err := h.articleService.List(c.Request.Context(), page, pageSize)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, 50003, "查询文章列表失败")
		return
	}
	response.Success(c, gin.H{
		"items":     result.Items,
		"total":     result.Total,
		"page":      page,
		"page_size": pageSize,
	})
}

// 文章详情
func (h *ArticleHandler) GetByID(c *gin.Context) {
	id, err := parseUintParam(c, "id")
	if err != nil {
		response.Error(c, http.StatusBadRequest, 40003, "文章ID不正确")
		return
	}
	article, err := h.articleService.GetByID(c.Request.Context(), uint(id))
	if errors.Is(err, service.ErrNotFound) {
		response.Error(c, http.StatusNotFound, 40004, "文章不存在")
		return
	}
	if err != nil {
		response.Error(c, http.StatusInternalServerError, 50002, "查询文章失败")
		return
	}
	response.Success(c, article)
}

// 删除文章
func (h *ArticleHandler) DeleteByID(c *gin.Context) {
	id, err := parseUintParam(c, "id")
	if err != nil {
		response.Error(c, http.StatusBadRequest, 40003, "文章ID不正确")
		return
	}

	currentUserInfo, ok := currentUser(c)
	if !ok {
		return
	}

	err = h.articleService.DeleteById(c.Request.Context(), currentUserInfo.UserID, uint(id))
	if errors.Is(err, service.ErrNotFound) {
		response.Error(c, http.StatusNotFound, 40004, "文章不存在")
		return
	}
	if errors.Is(err, service.ErrForbidden) {
		response.Error(c, http.StatusForbidden, 40007, "无权删除他人文章")
		return
	}
	if err != nil {
		response.Error(c, http.StatusInternalServerError, 50005, "删除文章失败")
		return
	}
	response.Success(c, nil)
}

// 更新文章
func (h *ArticleHandler) Update(c *gin.Context) {
	id, err := parseUintParam(c, "id")
	if err != nil {
		response.Error(c, http.StatusBadRequest, 40003, "文章ID不正确")
		return
	}
	var r req.UpdateArticleRequest
	if err := c.ShouldBindJSON(&r); err != nil {
		response.Error(c, http.StatusBadRequest, 40003, "参数错误"+err.Error())
		return
	}
	currentUserInfo, ok := currentUser(c)
	if !ok {
		return
	}

	article, err := h.articleService.Update(c.Request.Context(), service.UpdateArticleInput{
		CurrentUserID: currentUserInfo.UserID,
		ArticleID:     uint(id),
		Title:         r.Title,
		Content:       r.Content,
		Location:      r.Location,
		IsPublic:      r.IsPublic,
		Images:        r.Images,
		Topics:        r.Topics,
	})

	if errors.Is(err, service.ErrNotFound) {
		response.Error(c, http.StatusNotFound, 40004, "文章不存在")
		return
	}
	if errors.Is(err, service.ErrForbidden) {
		response.Error(c, http.StatusNotFound, 40007, "无权修改他人文章")
		return
	}
	if err != nil {
		response.Error(c, http.StatusInternalServerError, 50010, "更新文章失败")
		return
	}
	response.Success(c, article)
}

// 文章点赞
func (h *ArticleHandler) LikeYes(c *gin.Context) {
	id, err := parseUintParam(c, "id")
	if err != nil {
		response.Error(c, http.StatusBadRequest, 40003, "文章ID不正确")
		return
	}

	currentUserInfo, ok := currentUser(c)
	if !ok {
		return
	}

	err = h.articleService.Like(c.Request.Context(), currentUserInfo.UserID, uint(id))
	if errors.Is(err, service.ErrNotFound) {
		response.Error(c, http.StatusNotFound, 40004, "文章不存在")
		return
	}
	if err != nil {
		response.Error(c, http.StatusInternalServerError, 50006, "点赞失败")
		return
	}
	response.Success(c, nil)
}

func (h *ArticleHandler) LikeNo(c *gin.Context) {
	id, err := parseUintParam(c, "id")
	if err != nil {
		response.Error(c, http.StatusBadRequest, 40003, "文章ID不正确")
		return
	}

	currentUserInfo, ok := currentUser(c)
	if !ok {
		return
	}

	err = h.articleService.Unlike(c.Request.Context(), currentUserInfo.UserID, uint(id))
	if errors.Is(err, service.ErrNotFound) {
		response.Error(c, http.StatusNotFound, 40004, "文章不存在")
		return
	}
	if err != nil {
		response.Error(c, http.StatusInternalServerError, 50007, "取消点赞失败")
		return
	}
	response.Success(c, nil)
}

// 文章收藏
func (h *ArticleHandler) FavoriteYes(c *gin.Context) {
	id, err := parseUintParam(c, "id")
	if err != nil {
		response.Error(c, http.StatusBadRequest, 40003, "文章ID不正确")
		return
	}

	currentUserInfo, ok := currentUser(c)
	if !ok {
		return
	}

	err = h.articleService.Favorite(c.Request.Context(), currentUserInfo.UserID, uint(id))
	if errors.Is(err, service.ErrNotFound) {
		response.Error(c, http.StatusNotFound, 40004, "文章不存在")
		return
	}
	if err != nil {
		response.Error(c, http.StatusInternalServerError, 50008, "收藏失败")
		return
	}
	response.Success(c, nil)
}

func (h *ArticleHandler) FavoriteNo(c *gin.Context) {
	id, err := parseUintParam(c, "id")
	if err != nil {
		response.Error(c, http.StatusBadRequest, 40003, "文章ID不正确")
		return
	}

	currentUserInfo, ok := currentUser(c)
	if !ok {
		return
	}

	err = h.articleService.Unfavorite(c.Request.Context(), currentUserInfo.UserID, uint(id))
	if errors.Is(err, service.ErrNotFound) {
		response.Error(c, http.StatusNotFound, 40004, "文章不存在")
		return
	}
	if err != nil {
		response.Error(c, http.StatusInternalServerError, 50009, "取消收藏失败")
		return
	}
	response.Success(c, nil)
}
