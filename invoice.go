package main

import (
	//	"database/sql"
	// _ "github.com/go-sql-driver/mysql"
	//	"gopkg.in/gorp.v1"
	"github.com/satori/go.uuid"
	// "os"
	"fmt"
	"time"
)

type Invoice struct {
	Id           int `json:"id"`
	Link         string
	Created      time.Time
	OwnerID      int
	Debtor       string
	IncludeLinks bool

	Owner       User         `gorm:"-"`
	Subtotal    float64      `gorm:"-"`
	LinkedTotal float64      `gorm:"-"`
	Total       float64      `gorm:"-"`
	TotalString string       `gorm:"-"`
	Debts       []Debt       `gorm:"-"`
	LinkedDebts []linkedDebt `gorm:"-"`
}

func NewInvoice(ownerID int, debtor string) (i Invoice) {
	i = Invoice{}
	i.Created = time.Now()
	i.Link = uuid.NewV4().String()
	i.OwnerID = ownerID
	i.Debtor = debtor

	return i
}

func (i *Invoice) Create() {
	Database.Create(&i)
}

func (i *Invoice) Delete() error {
	Database.First(&i)
	Database.Delete(&i)
	return nil
}

func (i *Invoice) Update() error {
	i.Link = uuid.NewV4().String()

	Database.Save(&i)
	return nil
}

func GetInvoices(ownerid int) (invoices []Invoice, err error) {
	Database.Where("owner_id = ?", ownerid).Find(&invoices)
	return invoices, nil
}

func getLinkedDebtsForInvoice(invoice *Invoice) {
	//	invoice.LinkedDebts = GetLinkedDebts(invoice.OwnerID)

	AllDebtsByOwner := GetLinkedDebts(invoice.OwnerID)
	var LinkedDebtsWithCurrentDebtor []linkedDebt

	for _, l := range AllDebtsByOwner {
		if l.Owner == invoice.Debtor {
			fmt.Println("l.Owner == invoice.Debtor:", l.Owner, invoice.Debtor)
			//Inverse linked debts for logic and simplicity for user on invoice
			l.Amount = l.Amount * -1
			LinkedDebtsWithCurrentDebtor = append(LinkedDebtsWithCurrentDebtor, l)
		} else {
			fmt.Println("l.Owner != invoice.Debtor:", l.Owner, invoice.Debtor)
		}
	}

	invoice.LinkedDebts = LinkedDebtsWithCurrentDebtor

	for _, linkeddebt := range invoice.LinkedDebts {
		invoice.LinkedTotal = invoice.LinkedTotal + linkeddebt.Amount
	}
}

func ViewInvoice(link string) (invoice Invoice, err error) {
	//link := "c2be4f75-84d1-4064-a2a9-6675031b6a9c"

	fmt.Println("[viewInvoice]", "voor db")
	Database.Where("link = ?", link).First(&invoice)
	fmt.Println("[viewInvoice]", "na db", invoice)

	invoice.Owner = GetUser(invoice.OwnerID)
	fmt.Println("Debtor", invoice.Debtor)

	fmt.Println("[---- GetDebtsByName]")
	invoice.Debts, err = GetDebtsByName(invoice.OwnerID, invoice.Debtor)

	for _, debt := range invoice.Debts {
		invoice.Subtotal = invoice.Subtotal + debt.Amount
	}

	if invoice.IncludeLinks {
		getLinkedDebtsForInvoice(&invoice)
	}

	//Add debts and linkeddebts total instead of substract because LinkedTotal is already negative
	invoice.Total = invoice.Subtotal + invoice.LinkedTotal
	invoice.TotalString = fmt.Sprintf("%.2f", invoice.Total)

	if len(invoice.LinkedDebts) == 0 {
		fmt.Println("Setting false for IncludeLinks because of", len(invoice.LinkedDebts), "based on", invoice.LinkedDebts)
		invoice.IncludeLinks = false
	}

	//GetLinkedDebts(ownerID)

	fmt.Println(invoice.Debts)
	fmt.Println("Sum", invoice.Total)

	return invoice, err
}
