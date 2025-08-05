package utils

import (
	"fmt"
	"regexp"
)

func IsValidPassword(password string) bool {
	fmt.Printf("Validating password : [%s]\n", password) // for debug

	// Check minimum length
	if len(password) < 8 {
		fmt.Printf("Password too short: %d characters\n", len(password))
		return false
	}

	// Check for at least one lowercase letter
	lowercaseRegex := regexp.MustCompile(`[a-z]`)
	if !lowercaseRegex.MatchString(password) {
		fmt.Println("Missing lowercase letter")
		return false
	}

	// Check for at least one uppercase letter
	uppercaseRegex := regexp.MustCompile(`[A-Z]`)
	if !uppercaseRegex.MatchString(password) {
		fmt.Println("Missing uppercase letter")
		return false
	}

	// Check for at least one digit
	digitRegex := regexp.MustCompile(`\d`)
	if !digitRegex.MatchString(password) {
		fmt.Println("Missing digit")
		return false
	}

	// Check for at least one special character
	specialRegex := regexp.MustCompile(`[@$!%*?&#]`)
	if !specialRegex.MatchString(password) {
		fmt.Println("Missing special character")
		return false
	}

	// Check that password only contains allowed characters
	allowedRegex := regexp.MustCompile(`^[A-Za-z\d@$!%*?&#]+$`)
	if !allowedRegex.MatchString(password) {
		fmt.Println("Contains invalid characters")
		return false
	}

	fmt.Println("Password validation passed!")
	return true
}
