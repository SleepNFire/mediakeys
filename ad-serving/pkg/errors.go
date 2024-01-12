package pkg

import "errors"

var (
	ErrAdvertAlreadyExist    = errors.New("advertisement with the same ID already exists")
	ErrEnvironnementVariable = errors.New("failed to read environnement variable")
	ErrInternalError         = errors.New("internal error")
	ErrRedisUnaccessible     = errors.New("error during connecting to Redis")
)
