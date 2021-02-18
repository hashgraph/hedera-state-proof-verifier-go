package constants

const (
	Sha384Length           = 48
	Sha384WithRsaType      = 1
	Sha384WithRsaMaxLength = 384
)

const (
	ByteSize = 1
	IntSize  = 4
	LongSize = 8
)

const (
	// the sum of the length field and the checksum field
	SimpleSum = 101
)

const (
	RecordFileMarker      = 2
	SignatureFileV2Marker = 3
)

const (
	RecordFileFormatV1 = 1
	RecordFileFormatV2 = 2
	RecordFileFormatV5 = 5
)

const (
	SignatureFileFormatV4 = 4
	SignatureFileFormatV5 = 5
)

const (
	MaxRecordLength      = 64 * 1024
	MaxTransactionLength = 64 * 1024
)

const (
	// version, hapiVersion, previous hash marker, SHA-384 hash length
	PreV5HeaderLength = IntSize + IntSize + ByteSize + Sha384Length
	// version, hapi version major/minor/patch, object stream version
	V5StartHashOffset = IntSize + IntSize + IntSize + IntSize + IntSize
)
