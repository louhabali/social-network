package utils

import (
	"errors"
	"regexp"
	"strings"
	"time"
)

// Validate password (at least 8 characters, 1 uppercase, 1 lowercase, 1 number)
func IsValidPassword(password string) bool {
	if len(password) < 8 {
		return false
	}
	hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(password)
	hasLower := regexp.MustCompile(`[a-z]`).MatchString(password)
	hasDigit := regexp.MustCompile(`\d`).MatchString(password)
	return hasUpper && hasLower && hasDigit
}

// Validate email format
func IsValidEmail(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}

// Validate first and last name format
func IsValidName(name string) bool {
	nameRegex := regexp.MustCompile(`^[a-zA-Z]{1,10}$`)
	return nameRegex.MatchString(name)
}

// Validate date of birth (minimum age 5)
func IsValidDateOfBirth(dateStr string) bool {
	dob, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return false
	}

	age := time.Now().Year() - dob.Year()
	if time.Now().Month() < dob.Month() || (time.Now().Month() == dob.Month() && time.Now().Day() < dob.Day()) {
		age--
	}

	return age >= 5
}

// Allowed profile avatar formats (JPEG & PNG)
func IsValidAvatarExtension(ext string) bool {
	allowed := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
	}
	return allowed[strings.ToLower(ext)]
}

// ValidatePostImage checks if an uploaded image meets the criteria
func ValidatePostImage(fileExt string, fileSize int64) error {
	if fileSize > 4*1024*1024 {
		return errors.New("file too large, max size is 4MB")
	}
	if !IsValidPostImageExtension(fileExt) {
		return errors.New("invalid image type, supported formats: JPEG, PNG, GIF")
	}
	return nil
}
