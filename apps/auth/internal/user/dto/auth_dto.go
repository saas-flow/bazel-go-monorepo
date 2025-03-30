package dto

type ValidatePasswordRequest struct {
	Password string `form:"password" validate:"required"`
}

type LookupRequest struct {
	Username string `form:"username" validate:"required"`
}

type SignInWithPasswordRequest struct {
	Username string `form:"username" validate:"required"`
	Password string `form:"password" validate:"required"`
}

type SignUpRequest struct {
	Username  string `form:"username" validate:"required"`
	Password  string `form:"password" validate:"required"`
	FirstName string `form:"first_name" validate:"required"`
	LastName  string `form:"last_name" validate:"required"`
}
