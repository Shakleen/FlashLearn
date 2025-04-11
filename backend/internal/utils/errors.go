package utils

import "errors"

var (
	ErrDatabaseNotExist      = errors.New("database doesn't exist")
	ErrTableNotExist         = errors.New("table doesn't exist")
	ErrRecordNotExist        = errors.New("record doesn't exist")
	ErrMaxLengthExceeded     = errors.New("max length exceeded")
	ErrDuplicateKeyViolation = errors.New("duplicate key violation")
	ErrDeckNotExist          = errors.New("deck doesn't exist")
)
