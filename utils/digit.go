package utils

func IsHexDigit[T RuneOrByte](c T) bool {
	return (c >= '0' && c <= '9') || (c >= 'a' && c <= 'f') || (c >= 'A' && c <= 'F')
}

func IsUpperHexDigit[T RuneOrByte](c T) bool {
	return (c >= '0' && c <= '9') || (c >= 'A' && c <= 'F')
}

func IsLowerHexDigit[T RuneOrByte](c T) bool {
	return (c >= '0' && c <= '9') || (c >= 'a' && c <= 'f')
}
