package db

import (
	"fmt"

	"github.com/broothie/queuecumber/model"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type NotFound struct {
	EnvCollection string
	Lookup        string
}

func (e NotFound) Error() string {
	return fmt.Sprintf("'%s' not found in collection '%s'", e.Lookup, e.EnvCollection)
}

func IsNotFound(err error) bool {
	_, isNotFound := err.(NotFound)
	return isNotFound
}

func (db *DB) notFound(collection model.Collection, lookup string) NotFound {
	return NotFound{EnvCollection: db.fullCollectionName(collection), Lookup: lookup}
}

func isCodeNotFound(err error) bool {
	return status.Code(err) == codes.NotFound
}
