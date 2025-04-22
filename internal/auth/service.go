package auth

import (
	"errors"
	"time"
	"travel-service/internal/db"
	"travel-service/models"

	"github.com/golang-jwt/jwt/v5"

	"crypto/sha256"
	"encoding/hex"
)

func HashPassword(password string) string {
	hash := sha256.Sum256([]byte(password))
	return hex.EncodeToString(hash[:])
}

type Service interface {
	Register(input RegisterInput) (string, error)
	Login(input LoginInput) (string, error)
	UpdateUser(userID uint, input UpdateUserInput) error
}

type service struct{}

func NewService() Service {
	return &service{}
}

func (s *service) Register(input RegisterInput) (string, error) {
	var exists models.User
	err := db.DB.Where("user_email = ?", input.Email).First(&exists).Error
	if err == nil {
		return "", errors.New("user already exists")
	}

	user := models.User{
		FirstName:    input.FirstName,
		LastName:     input.LastName,
		UserEmail:    input.Email,
		PasswordHash: HashPassword(input.Password),
		Role:         "client",
		PhotoURLs:    []string{},
	}

	if err := db.DB.Create(&user).Error; err != nil {
		return "", err
	}

	// Генерация токена
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"role":    user.Role,
		"exp":     time.Now().Add(time.Hour * 168).Unix(), // 7 дней
	})

	tokenString, err := token.SignedString([]byte("your_secret_key"))
	if err != nil {
		return "", errors.New("failed to generate token")
	}

	return tokenString, nil
}

func (s *service) Login(input LoginInput) (string, error) {
	var user models.User
	err := db.DB.Where("user_email = ?", input.Email).First(&user).Error
	if err != nil {
		return "", errors.New("user not found")
	}

	// Проверка пароля (если SHA256, проверь вручную)
	if HashPassword(input.Password) != user.PasswordHash {
		return "", errors.New("invalid credentials")
	}

	// Генерация токена
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"role":    user.Role,
		"exp":     time.Now().Add(time.Hour * 168).Unix(), // 168 часа
	})

	tokenString, err := token.SignedString([]byte("your_secret_key")) // 🔐 тут секрет
	if err != nil {
		return "", errors.New("failed to generate token")
	}

	return tokenString, nil
}

func (s *service) UpdateUser(userID uint, input UpdateUserInput) error {
	var user models.User
	if err := db.DB.First(&user, userID).Error; err != nil {
		return err
	}

	if input.FirstName != nil {
		user.FirstName = *input.FirstName
	}
	if input.LastName != nil {
		user.LastName = *input.LastName
	}
	if input.Patronymic != nil {
		user.Patronymic = *input.Patronymic
	}
	if input.UserPhone != nil {
		user.UserPhone = *input.UserPhone
	}
	if input.BirthDate != nil {
		t, err := time.Parse("2006-01-02", *input.BirthDate)
		if err != nil {
			return err
		}
		user.BirthDate = t
	}
	if input.Citizenship != nil {
		user.Citizenship = *input.Citizenship
	}
	if input.HasChildren != nil {
		user.HasChildren = *input.HasChildren
	}
	if input.ChildrenInfo != nil {
		user.ChildrenInfo = *input.ChildrenInfo
	}
	if input.PassportNumber != nil {
		user.PassportNumber = *input.PassportNumber
	}
	if input.PhotoURLs != nil {
		user.PhotoURLs = *input.PhotoURLs
	}
	if input.Role != nil {
		user.Role = *input.Role
	}

	return db.DB.Save(&user).Error
}
