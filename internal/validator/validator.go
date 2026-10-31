package validator

import (
	"net/mail"
	"net/url"
	"regexp"
	"strings"
	"unicode"
)

type Result struct {
	Valid  bool     `json:"valid"`
	Errors []string `json:"errors"`
}

var (
	phonePattern = regexp.MustCompile(`^\+?[0-9][0-9 .()-]{5,19}$`)
)

func NormalizeType(value string) string {
	return strings.ToLower(strings.TrimSpace(value))
}

// Validate validates a value using one of the supported validation types.
func Validate(kind, value string) (Result, bool) {
	kind = NormalizeType(kind)
	value = strings.TrimSpace(value)

	errors := []string{}
	switch kind {
	case "email":
		address, err := mail.ParseAddress(value)
		domain := ""
		if err == nil {
			at := strings.LastIndex(address.Address, "@")
			if at >= 0 {
				domain = address.Address[at+1:]
			}
		}
		if err != nil || !strings.Contains(domain, ".") || strings.HasPrefix(domain, ".") || strings.HasSuffix(domain, ".") {
			errors = append(errors, "Email address has an invalid format")
		}
	case "cuit":
		if !validCUIT(value) {
			errors = append(errors, "CUIT must contain 11 digits and have a valid checksum")
		}
	case "phone":
		if !phonePattern.MatchString(value) {
			errors = append(errors, "Phone must contain 6 to 20 characters and use a valid phone format")
		}
	case "url":
		parsed, err := url.ParseRequestURI(value)
		if err != nil || parsed.Host == "" || (parsed.Scheme != "http" && parsed.Scheme != "https") {
			errors = append(errors, "URL must be an absolute http or https URL")
		}
	case "password":
		if len([]rune(value)) < 8 {
			errors = append(errors, "Password must contain at least 8 characters")
		}
		var hasUpper, hasLower, hasDigit, hasSpecial bool
		for _, r := range value {
			switch {
			case unicode.IsUpper(r):
				hasUpper = true
			case unicode.IsLower(r):
				hasLower = true
			case unicode.IsDigit(r):
				hasDigit = true
			case !unicode.IsLetter(r) && !unicode.IsDigit(r):
				hasSpecial = true
			}
		}
		if !hasUpper {
			errors = append(errors, "Password must contain an uppercase letter")
		}
		if !hasLower {
			errors = append(errors, "Password must contain a lowercase letter")
		}
		if !hasDigit {
			errors = append(errors, "Password must contain a number")
		}
		if !hasSpecial {
			errors = append(errors, "Password must contain a special character")
		}
	default:
		return Result{}, false
	}

	return Result{Valid: len(errors) == 0, Errors: errors}, true
}

func validCUIT(value string) bool {
	digits := make([]int, 0, 11)
	for _, r := range value {
		if r < '0' || r > '9' {
			continue
		}
		digits = append(digits, int(r-'0'))
	}
	if len(digits) != 11 {
		return false
	}

	weights := []int{5, 4, 3, 2, 7, 6, 5, 4, 3, 2}
	sum := 0
	for i, weight := range weights {
		sum += digits[i] * weight
	}
	check := 11 - (sum % 11)
	if check == 11 {
		check = 0
	}
	if check == 10 {
		check = 9
	}
	return digits[10] == check
}
