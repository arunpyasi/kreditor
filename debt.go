package main

import (
	//	"database/sql"
	"fmt"
	// _ "github.com/go-sql-driver/mysql"
	//	"gopkg.in/gorp.v1"
	"github.com/satori/go.uuid"
	// "os"
	"time"
)

type Debt struct {
	Id          int `json:"id"`
	Created     time.Time
	Amount      float64
	Paid        bool
	OwnerID     int
	Debtor      string
	Description string
	Permalink   string
}

type linkedDebt struct {
	Id          int `json:"id"`
	Created     time.Time
	Amount      float64
	Owner       string
	Description string
}

//var dbmap *gorp.DbMap

func init() {

}

// func initDb() *gorp.DbMap {
// 	fmt.Println("initialize database")
// 	// connect to db using standard Go database/sql API
// 	// use whatever database/sql driver you wish
// 	connectionString := "root:" + os.Getenv("KREDITOR_MYSQL_PASS") + "@tcp(localhost:3306)/kreditor?charset=utf8&parseTime=True"
// 	db, err := sql.Open("mysql", connectionString)
// 	if err != nil {
// 		//fmt.Println(err, "sql.Open failed")
// 		panic(err)
// 	}
// 	// construct a gorp DbMap
// 	dbmap := &gorp.DbMap{Db: db, Dialect: gorp.MySQLDialect{"InnoDB", "UTF-8"}}
//
// 	dbmap.AddTableWithName(Debt{}, "debt").SetKeys(false, "id")
// 	return dbmap
// }

func (d Debt) Delete() error {
	Database.First(&d)
	Database.Delete(&d)
	return nil
}

func (d *Debt) NewDebt(amount float64, ownerid int, debtor string) error {
	//d := new(Debt)
	//d.Id = uuid.NewV4().String()
	d.Created = time.Now()
	d.Amount = amount
	d.OwnerID = ownerid
	d.Debtor = debtor
	d.Permalink = uuid.NewV4().String()

	fmt.Println("Dit is NewDebt:", d)
	Database.Create(&d)
	return nil
}

func (d *Debt) UpdateDebt() error {
	fmt.Println("Updating debt with", d.Id)
	fmt.Println("Updating debt Description", d.Description)

	Database.Save(&d)
	return nil
}

func GetDebts(ownerid int) (debts []Debt, err error) {
	Database.Where("owner_id = ?", ownerid).Find(&debts)
	return debts, nil
}

func GetDebtsByName(ownerid int, debtor string) (debts []Debt, err error) {
	Database.Where("owner_id = ? AND debtor = ? AND paid = 0", ownerid, debtor).Find(&debts)
	return debts, nil
}

func GetLinkedDebts(ownerid int) []linkedDebt {
	//myOwnName := "Mike de Heij"
	myOwnName := GetUser(ownerid).Name
	fmt.Println("My own name is", myOwnName)

	var contactsWithLink []Contact
	Database.Where("owner_id = ? AND user_id > 0", ownerid).Find(&contactsWithLink)

	var linkedDebts []linkedDebt

	for _, c := range contactsWithLink {
		var debts []Debt
		fmt.Println("Looping over", c.Name, "which is", c.UserID)
		Database.Where("owner_id = ? AND debtor = ? AND paid = 0", c.UserID, myOwnName).Find(&debts)
		for _, debt := range debts {
			var linkedDebt linkedDebt
			ownerUser := GetUser(debt.OwnerID)

			linkedDebt.Id = debt.Id
			linkedDebt.Amount = debt.Amount
			linkedDebt.Created = debt.Created
			linkedDebt.Description = debt.Description

			if len(ownerUser.Name) > 0 {
				linkedDebt.Owner = ownerUser.Name
			} else {
				linkedDebt.Owner = ownerUser.Username
			}

			linkedDebts = append(linkedDebts, linkedDebt)
		}
	}

	return linkedDebts
}
