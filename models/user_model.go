package models

import (
	"database/sql"
	"fmt"

	"github.com/lethekhoi/twofactor/entities"
)

type UserModel struct {
	Db *sql.DB
}

//Create create new User with UserName and Pasword
func (userModel UserModel) Create(user *entities.User) (err error) {
	_, err = userModel.Db.Exec("insert into USERS(USER_NAME, PASSWORD) values(:1,:2)", user.Username, user.Password)

	if err != nil {
		fmt.Println("d")
		return err
	}
	return nil
}

//Create create new User with UserName and Pasword
func (userModel UserModel) Login(username, password string) (user entities.User, err error) {
	err = userModel.Db.QueryRow("select * from USERS where USER_NAME = :1 and PASSWORD = :2", username, password).Scan(&user.Username, &user.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			return user, sql.ErrNoRows
		}
		{
			return user, err
		}
	}
	return user, nil
}
