package model

import (
	"chatRoom/common/message"
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

func (receiver *UserDAO) Register(user *message.User) (err error) {

	//先从UserDao 的连接池中取出一根连接
	conn := receiver.pool.Get()
	defer conn.Close()
	_, err = receiver.GetUserByID(conn, user.ID)
	if err == nil {
		err = USEREXITS
		return
	}
	//这时，说明id在redis还没有，则可以完成注册
	data, err := json.Marshal(user) //序列化
	if err != nil {
		return
	}
	//入库
	_, err = conn.Do("HSet", "users", user.ID, string(data))
	if err != nil {
		fmt.Println("保存注册用户错误 err=", err)
		return
	}
	return
}
