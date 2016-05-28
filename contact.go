package main

import (
	"fmt"
)

type Contact struct {
	Id      int `json:"id"`
	OwnerID int
	Name    string
	UserID  int
	IBAN    string
}

type ContactList struct {
	Contacts []Contact
}

func init() {

	//	fmt.Println(getContacts())
}

func NewContact(ownerID int, name string) Contact {
	c := Contact{}

	c.OwnerID = ownerID
	c.Name = name
	c.UserID = -1

	return c
}

func (c *Contact) Create() {
	Database.Create(&c)
}

func (c *Contact) Update() error {
	Database.Save(&c)
	return nil
}

//TODO: Create []string and create []Contact methods
func GetContacts(ownerID int) []string {

	var names []string
	Database.Model(&Contact{}).Where("owner_id = ?", ownerID).Pluck("name", &names)

	return names
}

func GetContactObjects(ownerid int) (contacts []Contact) {
	Database.Where("owner_id = ?", ownerid).Find(&contacts)
	return contacts
}

func appendIfMissing(slice []string, i string) []string {
	for _, ele := range slice {
		if ele == i {
			return slice
		}
	}
	return append(slice, i)
}

func formatContactList(source []string) string {

	fmt.Println("Start contact test")

	var dump string
	for i, obj := range source {
		addition := "{value:'" + obj + "'}"
		if i != 0 {
			dump = dump + "," + addition
		} else {
			dump = addition
		}
	}
	return "{source:[" + dump + "]}"

}
