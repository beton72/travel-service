package analytics

import (
	"bytes"
	"fmt"
	"time"

	"github.com/drummonds/excelize/v2"
)

type ExcelService struct {
	service *Service
}

func NewExcelService(service *Service) *ExcelService {
	return &ExcelService{service: service}
}

func (s *ExcelService) GenerateExcelReport(hotelID uint, startDate, endDate time.Time) (*bytes.Buffer, error) {
	f := excelize.NewFile()
	defer f.Close()
	f.DeleteSheet("Sheet1")

	if err := s.createBookingsSheet(f, hotelID, startDate, endDate); err != nil {
		return nil, err
	}
	if err := s.createRevenueSheet(f, hotelID, startDate, endDate); err != nil {
		return nil, err
	}
	if err := s.createOccupancySheet(f, hotelID, startDate, endDate); err != nil {
		return nil, err
	}
	if err := s.createReviewsSheet(f, hotelID, startDate, endDate); err != nil {
		return nil, err
	}
	if err := s.createAverageCheckSheet(f, hotelID, startDate, endDate); err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	if err := f.Write(&buf); err != nil {
		return nil, fmt.Errorf("ошибка генерации Excel: %w", err)
	}
	return &buf, nil
}

func (s *ExcelService) createBookingsSheet(f *excelize.File, hotelID uint, startDate, endDate time.Time) error {
	sheet := "Bookings"
	f.NewSheet(sheet)
	f.SetCellValue(sheet, "A1", "Месяц")
	f.SetCellValue(sheet, "B1", "Бронирований")

	data, err := s.service.GetMonthlyBookings(hotelID, startDate, endDate)
	if err != nil {
		return err
	}

	for i, row := range data {
		f.SetCellValue(sheet, fmt.Sprintf("A%d", i+2), row.Month)
		f.SetCellValue(sheet, fmt.Sprintf("B%d", i+2), row.Count)
	}

	return s.addLineChart(f, sheet, "Количество бронирований", len(data), "D2")
}

func (s *ExcelService) createRevenueSheet(f *excelize.File, hotelID uint, startDate, endDate time.Time) error {
	sheet := "Revenue"
	f.NewSheet(sheet)
	f.SetCellValue(sheet, "A1", "Месяц")
	f.SetCellValue(sheet, "B1", "Доход ₽")

	data, err := s.service.GetMonthlyRevenue(hotelID, startDate, endDate)
	if err != nil {
		return err
	}

	for i, row := range data {
		f.SetCellValue(sheet, fmt.Sprintf("A%d", i+2), row.Month)
		f.SetCellValue(sheet, fmt.Sprintf("B%d", i+2), row.Revenue)
	}

	return s.addLineChart(f, sheet, "Доход", len(data), "D2")
}

func (s *ExcelService) createOccupancySheet(f *excelize.File, hotelID uint, startDate, endDate time.Time) error {
	sheet := "Occupancy"
	f.NewSheet(sheet)
	f.SetCellValue(sheet, "A1", "Месяц")
	f.SetCellValue(sheet, "B1", "Загруженность (%)")

	data, err := s.service.GetMonthlyOccupancy(hotelID, startDate, endDate)
	if err != nil {
		return err
	}

	for i, row := range data {
		f.SetCellValue(sheet, fmt.Sprintf("A%d", i+2), row.Month)
		f.SetCellValue(sheet, fmt.Sprintf("B%d", i+2), row.OccupancyRate)
	}

	return s.addLineChart(f, sheet, "Загруженность", len(data), "D2")
}

func (s *ExcelService) createReviewsSheet(f *excelize.File, hotelID uint, startDate, endDate time.Time) error {
	sheet := "Reviews"
	f.NewSheet(sheet)
	f.SetCellValue(sheet, "A1", "Месяц")
	f.SetCellValue(sheet, "B1", "Средняя оценка")

	data, err := s.service.GetMonthlyReviews(hotelID, startDate, endDate)
	if err != nil {
		return err
	}

	for i, row := range data {
		f.SetCellValue(sheet, fmt.Sprintf("A%d", i+2), row.Month)
		f.SetCellValue(sheet, fmt.Sprintf("B%d", i+2), row.AverageRating)
	}

	return s.addLineChart(f, sheet, "Оценки", len(data), "D2")
}

func (s *ExcelService) createAverageCheckSheet(f *excelize.File, hotelID uint, startDate, endDate time.Time) error {
	sheet := "AverageCheck"
	f.NewSheet(sheet)
	f.SetCellValue(sheet, "A1", "Месяц")
	f.SetCellValue(sheet, "B1", "Средний чек ₽")

	data, err := s.service.GetAverageCheck(hotelID, startDate, endDate)
	if err != nil {
		return err
	}

	for i, row := range data {
		f.SetCellValue(sheet, fmt.Sprintf("A%d", i+2), row.Month)
		f.SetCellValue(sheet, fmt.Sprintf("B%d", i+2), row.AvgCheck)
	}

	return s.addLineChart(f, sheet, "Средний чек", len(data), "D2")
}

func (s *ExcelService) addLineChart(f *excelize.File, sheet, title string, rows int, pos string) error {
	categories := fmt.Sprintf("%s!$A$2:$A$%d", sheet, rows+1)
	values := fmt.Sprintf("%s!$B$2:$B$%d", sheet, rows+1)

	return f.AddChart(sheet, pos, &excelize.Chart{
		Type: excelize.Line,
		Series: []excelize.ChartSeries{
			{
				Name:       fmt.Sprintf("%s!$B$1", sheet),
				Categories: categories,
				Values:     values,
			},
		},
		Title: excelize.ChartTitle{Name: title},
	})
}
