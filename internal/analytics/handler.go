package analytics

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
)

// Handler provides HTTP endpoints for analytics
type Handler struct {
	service Service
}

// NewHandler constructs a new analytics Handler
func NewHandler(s Service) *Handler {
	return &Handler{service: s}
}

// GenerateExcelReport builds an Excel report with key metrics and a chart, then streams it to the client
func (h *Handler) GenerateExcelReport(c *gin.Context) {
	// 1) Get overall dashboard stats
	stats, err := h.service.GetDashboardStats()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get dashboard stats"})
		return
	}

	// 2) Prepare revenue stats for the last 7 days
	endDate := time.Now()
	startDate := endDate.AddDate(0, 0, -6)
	req := RevenueReportRequest{
		Period: "daily",
		DateRange: &DateRangeFilter{
			StartDate: startDate,
			EndDate:   endDate,
		},
	}

	revStats, err := h.service.GetRevenueStats(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get revenue stats"})
		return
	}

	// 3) Create a new Excel file
	f := excelize.NewFile()

	// 4) Create and activate Overview sheet
	sheet := "Overview"
	sheetIndex, err := f.NewSheet(sheet)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create overview sheet"})
		return
	}
	f.SetActiveSheet(sheetIndex)

	// Write metrics table
	f.SetCellValue(sheet, "A1", "Metric")
	f.SetCellValue(sheet, "B1", "Value")
	metrics := [][]interface{}{
		{"Total Users", stats.TotalUsers},
		{"Total Hotels", stats.TotalHotels},
		{"Total Bookings", stats.TotalBookings},
		{"Total Revenue", stats.TotalRevenue},
		{"Average Rating", stats.AverageRating},
		{"Monthly Growth %", stats.MonthlyGrowth},
	}
	for i, row := range metrics {
		r := i + 2
		f.SetCellValue(sheet, fmt.Sprintf("A%d", r), row[0])
		f.SetCellValue(sheet, fmt.Sprintf("B%d", r), row[1])
	}

	// 5) Create Revenue sheet
	revSheet := "Revenue"
	if _, err := f.NewSheet(revSheet); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create revenue sheet"})
		return
	}
	// Write revenue data
	f.SetCellValue(revSheet, "A1", "Date")
	f.SetCellValue(revSheet, "B1", "Revenue")
	for i, rs := range revStats {
		r := i + 2
		f.SetCellValue(revSheet, fmt.Sprintf("A%d", r), rs.Date.Format("2006-01-02"))
		f.SetCellValue(revSheet, fmt.Sprintf("B%d", r), rs.Revenue)
	}

	// Add a line chart for revenue
	if err := f.AddChart(revSheet, "D2", &excelize.Chart{
		Type: excelize.Line,
		Series: []excelize.ChartSeries{{
			Name:       revSheet + "!$B$1",
			Categories: revSheet + "!$A$2:$A$" + fmt.Sprint(len(revStats)+1),
			Values:     revSheet + "!$B$2:$B$" + fmt.Sprint(len(revStats)+1),
		}},
	}); err != nil {
		// Ignore chart creation errors
	}

	// 6) Remove default sheet
	f.DeleteSheet("Sheet1")

	// 7) Stream the Excel file to client
	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Header("Content-Disposition", "attachment; filename=analytics_report.xlsx")
	if err := f.Write(c.Writer); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to write excel file"})
	}
}
