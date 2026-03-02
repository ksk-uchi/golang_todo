package validators

type CreateTodoRequest struct {
	Title       string `json:"title" validate:"required,max=100"`
	Description string `json:"description" validate:"max=200"`
}

func (r *CreateTodoRequest) Validate() map[string]string {
	if err := validate.Struct(r); err != nil {
		return TranslateError(err)
	}
	return nil
}

type UpdateTodoRequest struct {
	Title       *string `json:"title" validate:"omitempty,max=100"`
	Description *string `json:"description" validate:"omitempty,max=200"`
}

func (r *UpdateTodoRequest) Validate() map[string]string {
	if err := validate.Struct(r); err != nil {
		return TranslateError(err)
	}
	return nil
}

type UpdateDoneStatusRequest struct {
	IsDone *bool `json:"is_done" validate:"required"`
}

func (r *UpdateDoneStatusRequest) Validate() map[string]string {
	if err := validate.Struct(r); err != nil {
		return TranslateError(err)
	}
	return nil
}
