package validator

import (
	"fmt"
	"jastip-app/internal/customerror"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func ValidateID(id string) (primitive.ObjectID, error) {
	val, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return val, &customerror.Err{
			Code:   customerror.CodeErrInvalidRequest,
			Errors: fmt.Sprintf("invalid id: %s", id),
		}
	}
	return val, nil
}
