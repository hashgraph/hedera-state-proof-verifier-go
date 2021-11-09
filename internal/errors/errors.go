package errors

import "errors"

var (
	ErrorExtraDataInRecordFile                = errors.New("EXTRA_DATA_IN_RECORD_FILE")
	ErrorExtraDataInSignatureFile             = errors.New("EXTRA_DATA_IN_SIGNATURE_FILE")
	ErrorHashesNotMatch                       = errors.New("HASHES_NOT_MATCH")
	ErrorInvalidAddressBooksLength            = errors.New("INVALID_ADDRESS_BOOKS_LENGTH")
	ErrorInvalidChecksum                      = errors.New("INVALID_CHECKSUM")
	ErrorInvalidLength                        = errors.New("INVALID_LENGTH_FOUND")
	ErrorInvalidRecordFile                    = errors.New("INVALID_RECORD_FILE")
	ErrorInvalidSignatureFileType             = errors.New("INVALID_SIGNATURE_FILE_TYPE")
	ErrorInvalidSignaturesLength              = errors.New("INVALID_SIGNATURES_LENGTH")
	ErrorTransactionNotFound                  = errors.New("TRANSACTION_NOT_FOUND")
	ErrorUnexpectedSignatureFileTypeDelimiter = errors.New("UNEXPECTED_SIGNATURE_FILE_TYPE_DELIMITER")
	ErrorUnexpectedTypeDelimiter              = errors.New("UNEXPECTED_TYPE_DELIMITER")
	ErrorUnsupportedRecordFileMarker          = errors.New("UNSUPPORTED_RECORD_FILE_MARKER")
)
