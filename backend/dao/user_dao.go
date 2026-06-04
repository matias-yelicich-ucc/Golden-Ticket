package dao

import (
	"golden-ticket/backend/domain"
)

// UserDAO defines the data access operations for users
type UserDAO interface {
	Create(user *domain.User) error
	GetByEmail(email string) (*domain.User, error)
}

type userDAOImpl struct{}

// NewUserDAO creates a new instance of UserDAO
func NewUserDAO() UserDAO {
	return &userDAOImpl{}
}

// Create inserts a new user in the database
func (d *userDAOImpl) Create(user *domain.User) error {
	return DB.Create(user).Error
}

// GetByEmail retrieves a user by their email
func (d *userDAOImpl) GetByEmail(email string) (*domain.User, error) {
	var user domain.User
	err := DB.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}
