package main

import (
// _ "github.com/go-sql-driver/mysql"
//	"gopkg.in/gorp.v1"
//	"github.com/satori/go.uuid"
// "os"
)

type Role struct {
	Id          int
	ManageUsers bool
}

type User struct {
	//Id      int
	Id             int `json:"id"`
	Name           string
	Username       string
	Hash           string `json:"-"`
	EmailAddress   string
	IBAN           string
	InvoiceMessage string
	AvatarUrl      string //TODO: support this in invoices and/or application || EMPTY == GRAVATAR
	CustomCssUrl   string //TODO: support this in invoices and/or application
	CustomCss      string `sql:"type:text"`
	SidebarColor   string
	Administrator  bool
}

func (u *User) Create() {
	Database.Create(&u)
}

func (u *User) Delete() error {
	Database.First(&u)
	Database.Delete(&u)
	return nil
}

func (u *User) Update() error {
	Database.Save(&u)
	return nil
}

func GetUserByUsername(username string) (user User) {
	Database.Where("username = ?", username).First(&user)
	return user
}

func GetUser(id int) (user User) {
	Database.First(&user, id)
	return user
}

func GetUsers() (users []User) {
	Database.Find(&users)
	return users
}
