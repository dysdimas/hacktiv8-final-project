package controllers

type UserDto struct {
	Email    string
	Password string
	Name     string
}

var Users = []UserDto{}
