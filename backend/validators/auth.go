package validators

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

func (r *LoginRequest) Validate() map[string]string {
	if err := validate.Struct(r); err != nil {
		return TranslateError(err)
	}
	return nil
}
