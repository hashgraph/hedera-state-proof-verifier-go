package errors

import "errors"

var (
	ErrorHashesNotMatch            = errors.New("HASHES_NOT_MATCH")
	ErrorInvalidSignaturesLength   = errors.New("INVALID_SIGNATURES_LENGTH")
	ErrorInvalidRecordFile         = errors.New("INVALID_RECORD_FILE")
	ErrorInvalidAddressBooksLength = errors.New("INVALID_ADDRESS_BOOKS_LENGTH")
	ErrorTransactionNotFound       = errors.New("TRANSACTION_NOT_FOUND")
	ErrorUnexpectedTypeDelimiter   = errors.New("UNEXPECTED_TYPE_DELIMITER")
	ErrorVerifySignature           = errors.New("VERIFY_SIGNATURE_FAIL")
)
