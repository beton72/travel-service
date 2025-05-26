package analytics

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service      *Service
	excelService *ExcelService
}

func NewHandler(service *Service, excelService *ExcelService) *Handler {
	return &Handler{
		service:      service,
		excelService: excelService,
	}
}

// GET /analytics/hotels/:id/export?start_date=YYYY-MM-DD&end_date=YYYY-MM-DD
func (h *Handler) ExportExcel(c *gin.Context) {
	userRole := c.GetString("userRole")
	if userRole != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "доступ только для администраторов"})
		return
	}

	hotelIDStr := c.Param("id")
	hotelID, err := strconv.ParseUint(hotelIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "неверный ID отеля"})
		return
	}

	startStr := c.Query("start_date")
	endStr := c.Query("end_date")

	if startStr == "" || endStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "start_date и end_date обязательны"})
		return
	}

	startDate, err := time.Parse("2006-01-02", startStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "неверный формат start_date (нужен YYYY-MM-DD)"})
		return
	}

	endDate, err := time.Parse("2006-01-02", endStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "неверный формат end_date (нужен YYYY-MM-DD)"})
		return
	}

	// Генерация Excel
	report, err := h.excelService.GenerateExcelReport(uint(hotelID), startDate, endDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("ошибка генерации отчёта: %v", err)})
		return
	}

	filename := fmt.Sprintf("hotel_%d_analytics_%s_%s.xlsx", hotelID, startStr, endStr)

	c.Header("Content-Disposition", "attachment; filename="+filename)
	c.Data(http.StatusOK, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", report.Bytes())
}
