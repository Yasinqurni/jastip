package request

type RegisterRequest struct {
	Name        string `json:"name" binding:"required"`
	Address     string `json:"address" binding:"required"`
	Email       string `json:"email" binding:"required"`
	PhoneNumber string `json:"phone_number" binding:"required"`
	Password    string `json:"password" binding:"required"`
}

// GetJsonFieldName will return json tag name
func (req *RegisterRequest) GetJsonFieldName(field string) string {
	return map[string]string{
		"Name":        "name",
		"Address":     "address",
		"Email":       "email",
		"PhoneNumber": "phone_number",
		"Password":    "password",
	}[field]
}

// ErrMessages contains map of error messages
func (req *RegisterRequest) ErrMessages() map[string]map[string]string {
	return map[string]map[string]string{
		"name": {
			"required": "name is required",
		},
		"address": {
			"required": "address is required",
		},
		"email": {
			"required": "email is required",
		},
		"phone_number": {
			"required": "phone number is required",
		},
		"password": {
			"required": "password is required",
		},
	}
}
