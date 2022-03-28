package model

import (
	"encoding/json"
	"fmt"
	"github.com/gomodule/redigo/redis"
)

var (
	UDAO *UserDAO
)

// UserDAO DAO mean that Data Access Object
type UserDAO struct {
	pool *redis.Pool
}

func NewUserDAO(pool *redis.Pool) *UserDAO {
	return &UserDAO{
		pool: pool,
	}
}

func GetUserDaoPool() {
	return
}

func (receiver *UserDAO) GetUserByID(conn redis.Conn, userID int) (user *User, err error) {
	//hset users 123 "{\"userID\":123,\"userPWD\":\"456\"}"
	//               "{\"userID\":123,\"userPWD\":\"456\"}"
	fmt.Println("userID is ", userID)
	res, err := redis.String(conn.Do("HGet", "users", userID))
	fmt.Println("res is", res, "err is ", err)
	if err != nil {
		if err == redis.ErrNil {
			err = USERNOTEXITS
		}
		return
	}
	user = &User{}
	err = json.Unmarshal([]byte(res), user)
	if err != nil {
		fmt.Println("json Unmarshal err: ", err)
		return
	}
	return
}

func (receiver *UserDAO) CheckPWD(userID int, userPWD string) (bool, error) {
	conn := receiver.pool.Get()
	defer conn.Close()
	fmt.Println("start check user pwd")
	user, err := receiver.GetUserByID(conn, userID)
	fmt.Println(user)
	if err != nil {
		return false, err
	}
	if userPWD != user.UserPWD {
		err = USERPWDERROR
		return false, nil
	}
	fmt.Println("check user successful")
	return true, nil

}
