package handler

import (
	"net/http"
	"strconv"

	"github.com/CarlosAlbertoFurtado/notifyhub/internal/application"
	"github.com/CarlosAlbertoFurtado/notifyhub/internal/domain"
	"github.com/gin-gonic/gin"
)

type NotificationHandler struct {
	sendUC *application.SendNotificationUseCase
	repo   domain.NotificationRepository
}

func NewNotificationHandler(sendUC *application.SendNotificationUseCase, repo domain.NotificationRepository) *NotificationHandler {
	return &NotificationHandler{sendUC: sendUC, repo: repo}
}

func (h *NotificationHandler) Send(c *gin.Context) {
	var input application.SendNotificationInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	n, err := h.sendUC.Execute(c.Request.Context(), input)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{"status": "queued", "data": n})
}

func (h *NotificationHandler) GetByID(c *gin.Context) {
	n, err := h.repo.FindByID(c.Request.Context(), c.Param("id"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "notification not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": n})
}

func (h *NotificationHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	result, err := h.repo.FindAll(c.Request.Context(), domain.ListParams{
		Page:    page,
		Limit:   limit,
		Channel: c.Query("channel"),
		Status:  c.Query("status"),
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list notifications"})
		return
	}

	c.JSON(http.StatusOK, result)
}

func (h *NotificationHandler) Stats(c *gin.Context) {
	stats, err := h.repo.Stats(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get stats"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": stats})
}
