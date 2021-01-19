package dto

type SearchParams struct {
	Page      uint8   `json:"page"`
	PageSize  uint8   `json:"pageSize"`
	OrderBy   []Order `json:"orderBy"`
	Published bool    `json:"published"`
}

type Order struct {
	Field     string `json:"field"`
	Direction string `json:"direction"'`
}

type Authentication struct {
	Token string `json:"token"`
}

type RegistrationForm struct {
	Username        string `json:"username"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirmPassword"`
	DisplayName     string `json:"displayName"`
	Email           string `json:"email"`
}

type LoginForm struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
