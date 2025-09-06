package definitions

import "errors"

var (
	ErrInvalidToken       = errors.New("invalid token")
	ErrIPXEScriptNotFound = errors.New("iPXE script not found")
)
