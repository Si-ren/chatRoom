package message

type User struct {
	UserID     int    `json:"userID"`
	UserPWD    string `json:"userPWD"`
	UserStatus int    `json:"userStatus"`
}
