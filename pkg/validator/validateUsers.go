package validator

type UserValidatorParams struct {
	Email    string `json:"email"`
	Username string `json:"username,omitempty"`
	Password string `json:"password"`
}

func ValidateUserEmail(v *Validator, email string) {
	v.Check(email != "", "email", "email must not empty")
	v.Check(Matches(email, EmailRX), "email", "must be a valid email address")
}

func ValidateUserPassword(v *Validator, password string) {
	v.Check(password != "", "password", "password must not empty")
	v.Check(len(password) >= 8, "password", "must be at least 8 bytes long")
	v.Check(len(password) <= 16, "password", "must not be more than 16 bytes long")
}

func ValidateUserUsername(v *Validator, username string) {
	v.Check(username != "", "username", "username must not empty")
	v.Check(len(username) <= 50, "username", "must be more than 50 bytes long")
}

func ValidateUser(v *Validator, u *UserValidatorParams) {
	if u.Username == "" {
		ValidateUserEmail(v, u.Email)
		ValidateUserPassword(v, u.Password)
	} else {
		ValidateUserEmail(v, u.Email)
		ValidateUserPassword(v, u.Password)
		ValidateUserUsername(v, u.Username)
	}
}
