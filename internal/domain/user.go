package domain

type UserCreate struct {
	Email      string `json:"email" bind:"required"`
	Phone      string `json:"phone" bind:"required"`
	Password   string `json:"password" bind:"required"`
	PassHash   []byte `json:"-"`
	FirstName  string `json:"first_name"`
	SecondName string `json:"second_name"`
	CreateUserAdress
}

type UserLogin struct {
	Email    string `json:"email" bind:"required"`
	Password string `json:"password" bind:"required"`
}

type UserLoginDAO struct {
	ID       int64  `db:"id"`
	Email    string `db:"email"`
	PassHash []byte `db:"password"`
}
