package validation

import (
	v1 "github.com/dinhcanh303/go-microservices/api/auth/v1"
	"github.com/dinhcanh303/go-microservices/pkg/error"
	"github.com/dinhcanh303/go-microservices/pkg/val"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
)

func ValidateSignIn(req *v1.SignInRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := val.ValidateEmail(req.GetEmail()); err != nil {
		violations = append(violations, error.FieldViolation("email", err))
	}

	if err := val.ValidatePassword(req.GetPassword()); err != nil {
		violations = append(violations, error.FieldViolation("password", err))
	}
	return violations
}

func ValidateSignUp(req *v1.SignUpRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := val.ValidateEmail(req.GetEmail()); err != nil {
		violations = append(violations, error.FieldViolation("email", err))
	}
	if err := val.ValidatePassword(req.GetPassword()); err != nil {
		violations = append(violations, error.FieldViolation("password", err))
	}
	if err := val.ValidateName(req.GetFirstName()); err != nil {
		violations = append(violations, error.FieldViolation("first_name", err))
	}
	if err := val.ValidateName(req.GetLastName()); err != nil {
		violations = append(violations, error.FieldViolation("last_name", err))
	}
	return violations
}
