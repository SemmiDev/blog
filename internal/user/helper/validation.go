package helper

import (
	"regexp"
)

// ValidateEmail validates an email address using regexp
var (
	MailRegex   = regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	NumberRegex = regexp.MustCompile(`^[0-9]+$`)
)
