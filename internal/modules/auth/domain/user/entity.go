package user

type User struct {
	ID           string
	Username     string
	PasswordHash string `json:"-"`
}
