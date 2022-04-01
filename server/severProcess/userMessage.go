package severProcess

import "fmt"

//因为UserMgr 实例在服务器端有且只有一个
//因为在很多的地方，都会使用到，因此，我们
//将其定义为全局变量
var (
	userMgr *UserMgr
)

type UserMgr struct {
	onlineUsers map[int]*UserProcess
}

//完成对userMgr初始化工作
func init() {
	userMgr = &UserMgr{
		onlineUsers: make(map[int]*UserProcess, 1024),
	}
}

// AddOnlineUser 完成对onlineUsers添加
func (this *UserMgr) AddOnlineUser(up *UserProcess) {
	this.onlineUsers[up.UserID] = up
}

// DelOnlineUser 删除
func (this *UserMgr) DelOnlineUser(userID int) {
	delete(this.onlineUsers, userID)
}

// GetAllOnlineUser 返回当前所有在线的用户
func (this *UserMgr) GetAllOnlineUser() map[int]*UserProcess {
	return this.onlineUsers
}

// GetOnlineUserByID 根据ID返回对应的值
func (this *UserMgr) GetOnlineUserByID(userID int) (up *UserProcess, err error) {

	//从map取出一个值
	up, ok := this.onlineUsers[userID]
	if !ok { //说明，你要查找的这个用户，当前不在线。
		err = fmt.Errorf("用户%d 不存在", userID)
		return
	}
	return
}
