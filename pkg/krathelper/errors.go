package krathelper

import (
	"github.com/go-kratos/kratos/v2/errors"
)

// Common errors for use across services
var (
	ErrForbidden      = errors.Forbidden("FORBIDDEN", "access denied")
	ErrNotFound       = errors.NotFound("NOT_FOUND", "resource not found")
	ErrUnauthorized   = errors.Unauthorized("UNAUTHORIZED", "unauthorized access")
	ErrBadRequest     = errors.BadRequest("BAD_REQUEST", "invalid request")
	ErrInternalServer = errors.InternalServer("INTERNAL_SERVER", "internal server error")
)
