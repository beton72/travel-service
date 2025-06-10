package review

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service Service
}

func NewHandler(s Service) *Handler {
	return &Handler{service: s}
}

func (h *Handler) CreateReview(c *gin.Context) {
	userID, _ := c.Get("userID")
	var input CreateReviewInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "неверные данные"})
		return
	}

	if err := h.service.CreateReview(userID.(uint), input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "отзыв успешно создан"})
}

func (h *Handler) GetHotelReviews(c *gin.Context) {
	hotelID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "неверный ID отеля"})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	reviews, err := h.service.GetHotelReviews(uint(hotelID), page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ошибка при получении отзывов"})
		return
	}

	c.JSON(http.StatusOK, reviews)
}

func (h *Handler) GetReviewStats(c *gin.Context) {
	hotelID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "неверный ID отеля"})
		return
	}

	stats, err := h.service.GetReviewStats(uint(hotelID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ошибка при получении статистики"})
		return
	}

	c.JSON(http.StatusOK, stats)
}

func (h *Handler) GetRandomReview(c *gin.Context) {
	hotelID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "неверный ID отеля"})
		return
	}

	review, err := h.service.GetRandomReview(uint(hotelID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ошибка при получении отзыва"})
		return
	}
	if review.ID == 0 {
		c.JSON(http.StatusOK, gin.H{})
		return
	}

	resp := ReviewResponse{
		ID:        review.ID,
		UserID:    review.UserID,
		Rating:    review.Rating,
		Text:      review.Text,
		CreatedAt: review.CreatedAt,
	}
	c.JSON(http.StatusOK, resp)
}
