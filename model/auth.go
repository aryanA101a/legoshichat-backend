package model



type AuthResponse struct {
	User  User   `json:"user"`
	Token string `json:"token"`
}
type CreateAccountRequest struct {
	Name        string `json:"name" validate:"required,min=1"`
	PhoneNumber uint64 `json:"phoneNumber" validate:"required,min=1000000000,max=9999999999"`
	Password    string `json:"password" validate:"required,min=8"`
}
type LoginRequest struct {
	PhoneNumber uint64 `json:"phoneNumber" validate:"required,min=1000000000,max=9999999999"`
	Password    string `json:"password" validate:"required,min=8"`
}
type Account struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	PhoneNumber uint64 `json:"phoneNumber"`
	Password    string `json:"password"`
}
type User struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	PhoneNumber uint64 `json:"phoneNumber"`
}
