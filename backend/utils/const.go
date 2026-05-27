package utils

import "regexp"

// Pre-compiling these globally means they only parse ONCE when your app boots up.
var (
	// Enforces a minimum length of 8 characters
	MinLengthRegex = regexp.MustCompile(`.{8,}`)

	// Enforces at least one lowercase letter
	LowerRegex = regexp.MustCompile(`[a-z]`)

	// Enforces at least one uppercase letter
	UpperRegex = regexp.MustCompile(`[A-Z]`)

	// Enforces at least one digit/number
	DigitRegex = regexp.MustCompile(`\d`)

	// Enforces at least one special character symbol
	SpecialRegex = regexp.MustCompile(`[!@#$%^&*()_+\-=\[\]{};':"\\|,.<>\/?]`)
)

// Default Paginate Number
const DefaultPaginateNumber = 10
