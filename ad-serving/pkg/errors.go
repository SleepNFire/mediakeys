package pkg

import "errors"

var (
	ErrAdvertAlreadyExist    = errors.New("advertisement with the same ID already exists")
	ErrNotFound              = errors.New("errors the id does not exist")
	ErrEnvironnementVariable = errors.New("failed to read environnement variable")
	ErrInternalError         = errors.New("internal error")
	ErrRedisUnaccessible     = errors.New("error during connecting to Redis")
	ErrObjectUnknown         = errors.New("error: object unknown")
)

type RestMessage struct {
	Message string
}
