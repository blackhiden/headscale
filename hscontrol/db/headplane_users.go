package db

import (
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

const (
	// BcryptCost is the cost factor for bcrypt password hashing.
	BcryptCost = 12
)

var (
	ErrHeadplaneUserNotFound         = errors.New("headplane user not found")
	ErrHeadplaneUserExists           = errors.New("headplane user already exists")
	ErrInvalidUsername               = errors.New("invalid username")
	ErrInvalidPassword               = errors.New("invalid password")
	ErrHeadplaneUserNotAuthenticated = errors.New("user not authenticated")
	ErrHeadplaneUserForbidden        = errors.New("user forbidden")
	ErrHeadplanePasswordAuthRequired = errors.New("password authentication required")
)

// HeadplaneUser represents a Headplane UI user with credentials and role.
type HeadplaneUser struct {
	ID           uint   `gorm:"primaryKey"`
	Username     string `gorm:"uniqueIndex;not null"`
	PasswordHash string `gorm:"not null"`
	Role         string `gorm:"default:'user'"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// TableName specifies the table name for GORM.
func (HeadplaneUser) TableName() string {
	return "headplane_users"
}

// SetPassword hashes and sets the password for the user.
func (u *HeadplaneUser) SetPassword(password string) error {
	if password == "" {
		return ErrInvalidPassword
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), BcryptCost)
	if err != nil {
		return err
	}

	u.PasswordHash = string(hash)
	return nil
}

// CheckPassword verifies a password against the stored hash.
func (u *HeadplaneUser) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password))
	return err == nil
}

// IsAdmin returns true if the user has admin role.
func (u *HeadplaneUser) IsAdmin() bool {
	return u.Role == "admin"
}

// CreateHeadplaneUser creates a new Headplane user with the given username, password, and role.
func CreateHeadplaneUser(tx *gorm.DB, username, password, role string) (*HeadplaneUser, error) {
	if username == "" {
		return nil, ErrInvalidUsername
	}

	if password == "" {
		return nil, ErrInvalidPassword
	}

	// Check if user already exists
	var existing HeadplaneUser
	err := tx.Where("username = ?", username).First(&existing).Error
	if err == nil {
		return nil, ErrHeadplaneUserExists
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	user := &HeadplaneUser{
		Username: username,
		Role:     role,
	}

	if err := user.SetPassword(password); err != nil {
		return nil, err
	}

	if err := tx.Create(user).Error; err != nil {
		return nil, err
	}

	return user, nil
}

// GetHeadplaneUserByUsername retrieves a Headplane user by username.
func GetHeadplaneUserByUsername(tx *gorm.DB, username string) (*HeadplaneUser, error) {
	var user HeadplaneUser
	err := tx.Where("username = ?", username).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrHeadplaneUserNotFound
		}
		return nil, err
	}
	return &user, nil
}

// GetHeadplaneUserByID retrieves a Headplane user by ID.
func GetHeadplaneUserByID(tx *gorm.DB, id uint) (*HeadplaneUser, error) {
	var user HeadplaneUser
	err := tx.First(&user, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrHeadplaneUserNotFound
		}
		return nil, err
	}
	return &user, nil
}

// ListHeadplaneUsers lists all Headplane users.
func ListHeadplaneUsers(tx *gorm.DB) ([]HeadplaneUser, error) {
	var users []HeadplaneUser
	err := tx.Order("created_at ASC").Find(&users).Error
	if err != nil {
		return nil, err
	}
	return users, nil
}

// DeleteHeadplaneUser deletes a Headplane user by ID.
func DeleteHeadplaneUser(tx *gorm.DB, id uint) error {
	result := tx.Delete(&HeadplaneUser{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrHeadplaneUserNotFound
	}
	return nil
}

// UpdateHeadplaneUserPassword updates a user's password.
func UpdateHeadplaneUserPassword(tx *gorm.DB, id uint, newPassword string) error {
	user, err := GetHeadplaneUserByID(tx, id)
	if err != nil {
		return err
	}

	if err := user.SetPassword(newPassword); err != nil {
		return err
	}

	return tx.Model(user).Update("password_hash", user.PasswordHash).Error
}

// UpdateHeadplaneUser updates a user's username and/or role.
func UpdateHeadplaneUser(tx *gorm.DB, id uint, username, role string) (*HeadplaneUser, error) {
	user, err := GetHeadplaneUserByID(tx, id)
	if err != nil {
		return nil, err
	}

	updates := make(map[string]interface{})

	// Only update username if provided and different
	if username != "" && username != user.Username {
		// Check if new username is already taken
		var existing HeadplaneUser
		err := tx.Where("username = ? AND id != ?", username, id).First(&existing).Error
		if err == nil {
			return nil, ErrHeadplaneUserExists
		}
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		updates["username"] = username
	}

	// Only update role if provided and different
	if role != "" && role != user.Role {
		if role != "user" && role != "admin" {
			return nil, errors.New("invalid role: must be 'user' or 'admin'")
		}
		updates["role"] = role
	}

	// If no updates, return current user
	if len(updates) == 0 {
		return user, nil
	}

	if err := tx.Model(user).Updates(updates).Error; err != nil {
		return nil, err
	}

	// Reload user to get updated values
	return GetHeadplaneUserByID(tx, id)
}

// CountAdminUsers returns the number of users with admin role.
func CountAdminUsers(tx *gorm.DB) (int64, error) {
	var count int64
	err := tx.Model(&HeadplaneUser{}).Where("role = ?", "admin").Count(&count).Error
	return count, err
}
