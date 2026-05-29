package repository

// UserRepository handles all user database queries.
type UserRepository struct{}

func NewUserRepository() *UserRepository {
	return &UserRepository{}
}

// TODO: FindByUsername, FindByID, Create
