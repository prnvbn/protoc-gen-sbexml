package sbexml

import "errors"

var (
	errFileToGenerateNotFound = errors.New("file to generate not found")
	errUnsupportedEnumType    = errors.New("unsupported enum type")
	errUnsupportedMessageType = errors.New("unsupported message type")
	errUnsupportedFieldType   = errors.New("unsupported field type")
)
