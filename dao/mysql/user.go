package mysql

import (
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"web_app/models"
)

const secret = "wanglongzhi.com"

// CheckUserExist 在数据库中查询用户是否存在
func CheckUserExist(username string) (err error) {
	sqlStr := `select count(user_id) from users where username = ?`
	var count int
	if err := db.Get(&count, sqlStr, username); err != nil {
		return err
	}
	if count > 0 {
		return ErrorUserExist
	}
	return
}

// InsertUser 向数据库中插入一条新的用户记录
func InsertUser(user *models.User) (err error) {
	// 对密码进行加密
	user.Password = encryptPassword(user.Password)
	sqlStr := `insert into users(user_id, username, password) values(?,?,?)`
	_, err = db.Exec(sqlStr, user.UserID, user.Username, user.Password)
	return
}

func encryptPassword(oPassword string) string {
	h := md5.New()
	h.Write([]byte(secret))
	return hex.EncodeToString(h.Sum([]byte(oPassword)))
}

func Login(user *models.User) (err error) {
	// 1.判断用户是否存在
	oPassword := user.Password // 用户登录时自己传的密码
	sqlStr := `select user_id, username, password from users where username=?`
	err = db.Get(user, sqlStr, user.Username)
	if err == sql.ErrNoRows {
		return ErrorUserNotExist
	}
	if err != nil {
		// 查询数据库失败
		return err
	}
	// 判断加密后的密码是否正确
	password := encryptPassword(oPassword)
	if password != user.Password {
		return ErrorInvalidPassword
	}
	return
}

func GetUserById(authorID int64) (user *models.User, err error) {
	user = new(models.User)
	sqlStr := `select
			user_id, username
			from users where user_id = ?`
	err = db.Get(user, sqlStr, authorID)
	return
}
