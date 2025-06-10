package hotel

import (
	"log"
	"net/http"
	"travel-service/internal/db"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) CreateHotel(c *gin.Context) {
	userIDInterface, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	// Получаем роль пользователя
	roleInterface, roleExists := c.Get("userRole")
	role, ok := roleInterface.(string)

	if !roleExists || !ok || role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "доступ только для администраторов"})
		return
	}

	userID := userIDInterface.(uint)

	var input CreateHotelInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}

	if err := h.service.CreateHotel(input, userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create hotel", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Hotel created successfully"})
}

func (h *Handler) GetHotels(c *gin.Context) {
	hotels, err := h.service.GetHotels()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get hotels"})
		return
	}
	c.JSON(http.StatusOK, hotels)
}

func (h *Handler) AddAdminToHotel(c *gin.Context) {
	var input AddAdminInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Некорректный ввод"})
		return
	}

	if err := h.service.AddAdminToHotel(input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Администратор добавлен к отелю"})
}

func (h *Handler) GetHotelByID(c *gin.Context) {
	id := c.Param("id")

	// Инкремент просмотров
	if err := h.service.IncrementViews(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update views"})
		return
	}

	// Получение отеля
	hotel, err := h.service.GetHotelByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "hotel not found"})
		return
	}

	c.JSON(http.StatusOK, hotel)
}

func (h *Handler) FilterHotelsByPrice(c *gin.Context) {
	var input FilterHotelsInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hotels, err := h.service.FilterHotelsByPriceRange(input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, hotels)
}

func (h *Handler) SearchAvailableRooms(c *gin.Context) {
	var input SearchRoomsInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	rooms, err := h.service.SearchAvailableRooms(input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, rooms)
}

func (h *Handler) GetMyHotels(c *gin.Context) {
	userID := c.GetUint("userID")

	hotels, err := h.service.GetHotelsByAdmin(userID)
	if err != nil {
		log.Println("ошибка получения отелей:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "не удалось получить отели"})
		return
	}

	c.JSON(http.StatusOK, hotels)
}

func IsRoomBookedToday(roomID uint) (bool, error) {
	var count int64
	err := db.DB.
		Table("bookings").
		Where("room_id = ? AND start_date <= CURRENT_DATE AND end_date >= CURRENT_DATE", roomID).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}


