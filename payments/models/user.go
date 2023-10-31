package models

import (
	"errors"
	"regexp"
)

type User struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	Email        string `json:"email"`
	PasswordHash string `json:"password_hash"`
	IsActive     bool   `json:"is_active"`
}

var (
	errNameIsEmpty     = errors.New("name of the user is empty")
	errEmailIsEmpty    = errors.New("email of the user is empty")
	errEmailIsNotValid = errors.New("email of the user is not valid")
	errPasswordIsEmpty = errors.New("password of the user is empty")
)

func (u *User) SignUpValidation() (err error) {
	if isEmpty(u.Name) {
		return errNameIsEmpty
	}
	if isEmpty(u.Email) {
		return errEmailIsEmpty
	}
	if !isEmailValid(u.Email) {
		return errEmailIsNotValid
	}
	if isEmpty(u.PasswordHash) {
		return errPasswordIsEmpty
	}
	return nil
}

func (u *User) LoginValidation() (err error) {
	if isEmpty(u.Email) {
		return errEmailIsEmpty
	}
	if !isEmailValid(u.Email) {
		return errEmailIsNotValid
	}
	if isEmpty(u.PasswordHash) {
		return errPasswordIsEmpty
	}
	return nil
}

func isEmpty(s string) bool {
	return s == ""
}

func isEmailValid(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	return emailRegex.MatchString(email)
}
