package policies

import "errors"

var (
	ErrUserAlreadyHasBroaderPolicy = errors.New("user already has broader policy; no need to add")
)
