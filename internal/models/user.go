package models

import (
	"errors"
)

type UserStatus string

const (
	UserStatusActive    UserStatus = "active"
	UserStatusSuspended UserStatus = "suspended"
	UserStatusDeleted   UserStatus = "deleted"
)

// validate the user status
func (s UserStatus) IsValid() bool {
	switch s {
	case UserStatusActive, UserStatusSuspended, UserStatusDeleted:
		return true
	}
	return false
}

type User struct {
	Base           `bson:",inline"`
	Email          string     `bson:"email" json:"email"`
	HashedPassword string     `bson:"hashed_password" json:"-"`
	Name           string     `bson:"name" json:"name"`
	Status         UserStatus `bson:"status" json:"status"` //active, suspended, deleted
}

// validation logic
func (u *User) Validate() error {
	if u.Email == "" {
		return errors.New("email is required")
	}
	if u.Name == "" {
		return errors.New("name is required")
	}

	if !u.Status.IsValid() {
		return errors.New("invalid user status")
	}
	return nil
}

// BeforeCreate hooks into the creation process
func (u *User) BeforeCreate() error {
	if u.Status == "" {
		u.Status = UserStatusActive //default
	}
	u.Base.BeforeCreate()
	return u.Validate()
}

// BeforeUpdate hooks into the update process
func (u *User) BeforeUpdate() error {
	u.Base.BeforeUpdate()
	return u.Validate()
}
