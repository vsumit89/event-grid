package utils

import "regexp"

func ValidateEmail(email string) bool {
	// Define the regex pattern for email validation
	pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	// Compile the regex pattern
	re := regexp.MustCompile(pattern)
	// Match the email against the pattern
	return re.MatchString(email)
}
