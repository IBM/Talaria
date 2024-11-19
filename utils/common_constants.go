package utils

import "strings"

type SecurityProtocol int

const (
	PLAINTEXT SecurityProtocol = iota
	SSL
	SASL_PLAINTEXT
	SASL_SSL
	UNDEFINED
)

// ParseSecurityProtocol parses the string p and returns the corresponding SecurityProtocol enum value and true
// or UNDEFINED and false if the value of string p is not recognized.
func ParseSecurityProtocol(p string) (SecurityProtocol, bool) {
	switch strings.ToUpper(p) {
	case "PLAINTEXT":
		return PLAINTEXT, true
	case "SSL":
		return SSL, true
	case "SASL_PLAINTEXT":
		return SASL_PLAINTEXT, true
	case "SASL_SSL":
		return SASL_SSL, true
	default:
		return UNDEFINED, false
	}
}
