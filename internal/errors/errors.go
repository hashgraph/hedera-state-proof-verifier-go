package errors

import "errors"

var (
	ErrorExtraDataInRecordFile       = errors.New("EXTRA_DATA_IN_RECORD_FILE")
	ErrorHashesNotMatch              = errors.New("HASHES_NOT_MATCH")
	ErrorInvalidSignaturesLength     = errors.New("INVALID_SIGNATURES_LENGTH")
	ErrorInvalidRecordFile           = errors.New("INVALID_RECORD_FILE")
	ErrorInvalidAddressBooksLength   = errors.New("INVALID_ADDRESS_BOOKS_LENGTH")
	ErrorTransactionNotFound         = errors.New("TRANSACTION_NOT_FOUND")
	ErrorUnexpectedTypeDelimiter     = errors.New("UNEXPECTED_TYPE_DELIMITER")
	ErrorUnsupportedRecordFileMarker = errors.New("UNSUPPORTED_RECORD_FILE_MARKER")
	ErrorVerifyMetadataSignature     = errors.New("VERIFY_METADATA_SIGNATURE")
	ErrorVerifySignature             = errors.New("VERIFY_SIGNATURE_FAIL")
)
