package model

import (
	"go.uber.org/zap"
	"prompting/pkg/db/mysql"
	"time"
	"xorm.io/builder"
)

func init() {
	// 生成数据表
	err := mysql.Client.Sync2(new(User))
	if err != nil {
		zap.L().Error("User Init table err: ", zap.Error(err))
		return
	}
}

type User struct {
	Id        string    `xorm:"id pk" json:"id"`
	Username  string    `xorm:"username" json:"username"`
	Nickname  string    `xorm:"nickname" json:"nickname"`
	Password  string    `xorm:"password" json:"password"`
	Email     string    `xorm:"email" json:"email"`
	Phone     string    `xorm:"phone" json:"phone"`
	Avatar    string    `xorm:"avatar" json:"avatar"`
	Role      int64     `xorm:"role" json:"role"`
	CreatedAt time.Time `xorm:"createdAt" json:"createdAt"`
	UpdatedAt time.Time `xorm:"updatedAt" json:"updatedAt"`
}

// 插入一条数据
func (u *User) InsertOneUser() error {
	_, err := mysql.Client.Insert(&u)
	return err
}

// 删除一条数据
func (u *User) DeleteUserById() error {
	_, err := mysql.Client.ID(u.Id).Delete()
	return err
}

// 批量删除
func DeleteUsersById(ids []int64) (int64, error) {
	// 构建删除条件
	cond := builder.In("id", ids)

	// 执行删除
	affected, err := mysql.Client.Where(cond).Delete(&User{})
	return affected, err
}

// 查找
func (u *User) QueryUserById() error {
	err := mysql.Client.Where("id = ?", u.Id).Find(&u)
	return err
}

// 获得用户权限
func (u *User) QueryRole() (int64, error) {
	var roles []int64
	err := mysql.Client.Where("id = ?", u.Id).Cols("role").Find(&roles)
	return roles[0], err
}

// 修改用户信息
func (u *User) UpdateUserInfo() error {
	_, err := mysql.Client.ID(u.Id).Update(&u)
	return err
}
