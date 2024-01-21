package request

type LoginRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// GetJsonFieldName will return json tag name
func (req *LoginRequest) GetJsonFieldName(field string) string {
	return map[string]string{
		"Email":    "email",
		"Password": "password",
	}[field]
}

// ErrMessages contains map of error messages
func (req *LoginRequest) ErrMessages() map[string]map[string]string {
	return map[string]map[string]string{
		"email": {
			"required": "email is required",
		},
		"password": {
			"required": "password is required",
		},
	}
}
