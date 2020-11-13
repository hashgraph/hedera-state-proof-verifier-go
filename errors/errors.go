package errors

import "errors"

var (
	ErrorInvalidSignaturesLength   = errors.New("INVALID_SIGNATURES_LENGTH")
	ErrorInvalidRecordFile         = errors.New("INVALID_RECORD_FILE")
	ErrorInvalidAddressBooksLength = errors.New("INVALID_ADDRESS_BOOKS_LENGTH")
)
