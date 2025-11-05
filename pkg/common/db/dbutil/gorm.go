package dbutil

import (
	"github.com/openimsdk/tools/errs"
	"go.mongodb.org/mongo-driver/mongo"
)

func IsDBNotFound(err error) bool {
	return errs.Unwrap(err) == mongo.ErrNoDocuments
}
