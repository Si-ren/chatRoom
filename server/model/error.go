package model

import "errors"

var (
	USERNOTEXITS error = errors.New("user not exits")
	USEREXITS          = errors.New("user already exits")
	USERPWDERROR error = errors.New("user PWD Error")
)
