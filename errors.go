package binary

import "errors"

var (
	ErrBufTooSmall        = errors.New("buffer is too small")
	ErrNonStringTailZero  = errors.New("strings not found zero tail")
	ErrInvalidMapKey      = errors.New("invalid map key Type")
	ErrInvalidMapValue    = errors.New("invalid map value Type")
	ErrInvalidStructField = errors.New("invalid struct field Type")
	ErrInvalidStructValue = errors.New("invalid struct value Type")
	ErrMustScalarType     = errors.New("must scalar type")
	ErrInvalidElementType = errors.New("invalid element type")
)
