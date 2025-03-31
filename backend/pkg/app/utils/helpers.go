package utils

import (
	"errors"
	"log"
	"os"
	"strings"
	"time"

	"github.com/gofrs/uuid"
)

// Generate a UUID
func GenerateToken() (string, error) {
	u2, err := uuid.NewV4()
	if err != nil {
		log.Fatalf("failed to generate UUID: %v", err)
		return "", err
	}
	return u2.String(), nil
}

// Allowed post image formats (JPEG, PNG & GIF)
func IsValidPostImageExtension(ext string) bool {
	allowed := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
		".gif":  true,
	}
	return allowed[strings.ToLower(ext)]
}

// GetOptionalString returns a pointer to the string if it's not empty, otherwise nil
func GetOptionalString(value string) *string {
	if value == "" {
		return nil
	}
	return &value
}

// validateAndFormatDate checks if the date is valid and formats it correctly
func ValidateAndFormatDate(inputDate string) (string, error) {
	parsedTime, err := time.Parse("2006-01-02T15:04", inputDate)
	if err != nil {
		return "", err
	}

	if parsedTime.Before(time.Now()) {
		return "", errors.New("the event date cannot be in the past")
	}

	return parsedTime.Format("2006-01-02 15:04:05"), nil
}

func CreateDirectories() error {
	err := os.MkdirAll("uploads/avatars", os.ModePerm)
	if err != nil {
		return err
	}

	err = os.MkdirAll("uploads/images", os.ModePerm)
	if err != nil {
		return err
	}
	return nil
}
