package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/mdeheij/kreditor/utils"
	"html/template"
)

var router *gin.Engine

func main() {
	router = gin.Default()

	//	router.NoRoute(homepage)

	router.LoadHTMLGlob("templates/*")
	router.Static("/assets", "assets")

	router.Use(func(context *gin.Context) {
		// add header Access-Control-Allow-Origin
		context.Writer.Header().Add("Access-Control-Allow-Origin", "*")
		context.Writer.Header().Add("Access-Control-Allow-Headers", "origin, content-type, accept")
		context.Writer.Header().Add("Access-Control-Allow-Methods", "POST, GET, PUT, DELETE, OPTIONS")
		context.Next()
	})
	loginInit(router)

	router.GET("/", splashPage)

	router.GET("/home", func(c *gin.Context) {
		c.Redirect(301, "/debts")
	})

	router.GET("/i/:link", routeViewInvoice)

	pagesGroup := router.Group("/", AuthRequired())
	{
		pagesGroup.GET("/debts", debtsPage)
		pagesGroup.GET("/invoices", invoices)
		pagesGroup.GET("/profile", profile)
		pagesGroup.GET("/contacts", contacts)
		pagesGroup.GET("/users", contacts)
	}
	apiGroup := router.Group("/api", AuthRequired())
	{

		apiGroup.GET("/debts", getAllDebts)
		apiGroup.GET("/debtlinks", routeGetLinkedDebts)

		apiGroup.OPTIONS("/debts", getAllDebts)
		apiGroup.POST("/debts", addDebt)
		apiGroup.POST("/debts/:id", updateDebt)
		apiGroup.DELETE("/debts/:id", deleteDebt)

		apiGroup.GET("/invoice", routeGetInvoices)
		apiGroup.POST("/invoice", routeCreateInvoice)
		apiGroup.POST("/invoice/:id", routeUpdateInvoice)
		apiGroup.DELETE("/invoice/:id", routeDeleteInvoice)

		apiGroup.GET("/profile", routeGetProfile)
		apiGroup.POST("/profile", routeUpdateProfile)
		apiGroup.POST("/profile/password", routeUpdateProfilePassword)

		apiGroup.GET("/contact", routeGetContacts)
		apiGroup.POST("/contact", routeCreateContact)
		apiGroup.POST("/contact/:id", routeUpdateContact)
		apiGroup.DELETE("/contact/:id", routeDeleteContact)
	}
	adminGroup := router.Group("/admin", AuthRequired())
	{
		adminGroup.GET("/", admin)
		adminGroup.GET("/user", routeGetUsers)
		adminGroup.POST("/user", routeCreateUser)
		adminGroup.POST("/user/id", routeUpdateUser)
		adminGroup.DELETE("/user/:id", routeDeleteUser)
	}

	router.Run()
	// router.Run(":3000") for a hard coded port
}

type navigationItem struct {
	Href string
	Name string
	Icon string
}

func getNavigation(CurrentUser User) []navigationItem {
	var items []navigationItem

	items = append(items, navigationItem{Href: "/debts", Name: "Debts", Icon: "fa-eur"})
	items = append(items, navigationItem{Href: "/invoices", Name: "Invoices", Icon: "fa-paperclip"})
	items = append(items, navigationItem{Href: "/contacts", Name: "Contacts", Icon: "fa-user-plus"})

	if CurrentUser.Administrator {
		items = append(items, navigationItem{Href: "/admin/", Name: "Admin", Icon: "fa-sliders"})
	}

	return items
}

func getAllDebts(c *gin.Context) {
	debts, err := GetDebts(getUserID(c))
	if err == nil {
		c.JSON(200, debts)
	} else {
		c.JSON(500, debts)
	}
}

/*
Router Invoice Code
*/

func routeGetInvoices(c *gin.Context) {
	invoices, err := GetInvoices(getUserID(c))
	if err == nil {
		c.JSON(200, invoices)
	} else {
		c.JSON(500, invoices)
	}
}

func routeDeleteInvoice(c *gin.Context) {
	var invoice Invoice
	id := c.Param("id")
	fmt.Println("[RouteDeleteInvoice] Id", id)
	err := Database.Where("id = ?", id).First(&invoice)

	fmt.Println("[RouteDeleteInvoice] [debug] Struct of &invoice", &invoice)

	if invoice.OwnerID == getUserID(c) {
		Database.Delete(&invoice)
	} else {
		//TODO: Error handling
		fmt.Println("[RouteDeleteInvoice] error occured", err)
		fmt.Println("[RouteDeleteInvoice] error occured; Id =", invoice.Id)
	}

	c.JSON(200, gin.H{
		"result": "done",
	})
}

func routeCreateInvoice(c *gin.Context) {
	postedInvoice := Invoice{}

	if c.BindJSON(&postedInvoice) == nil {
		invoice := NewInvoice(getUserID(c), postedInvoice.Debtor)
		invoice.IncludeLinks = postedInvoice.IncludeLinks
		invoice.Create()
		c.JSON(200, invoice)
	} else {
		c.String(500, "Something went wrong in binding (add)")
	}
}

func routeUpdateInvoice(c *gin.Context) {
	invoice := Invoice{}
	Database.Where("id = ? AND owner_id = ?", c.Param("id"), getUserID(c)).First(&invoice)

	if invoice.OwnerID == getUserID(c) {
		invoice.Update()
	}

	c.JSON(200, invoice)
}

func routeViewInvoice(c *gin.Context) {
	link := c.Param("link")
	invoice, _ := ViewInvoice(link)

	owner := GetUser(invoice.OwnerID)

	if len(invoice.Debts) == 0 {
		c.String(404, "Invoice could not be found! The link is probably expired, or all debts are already marked as paid.")
	} else {
		c.HTML(200, "publicinvoice.html", gin.H{
			"invoice":     invoice,
			"linkeddebts": GetLinkedDebts(invoice.OwnerID),
			//"printf":    fmt.Printf,
			"link":      link,
			"userCSS":   template.CSS(owner.CustomCss),
			"INGQRcode": utils.GetQRCode(owner.Name, owner.IBAN, invoice.TotalString),
			"contacts":  formatContactList(GetContacts(getUserID(c))),
		})
	}
}

/*
End of Invoice Routes
*/

/*
Router Contact Code
*/

func routeGetContacts(c *gin.Context) {
	contacts := GetContactObjects(getUserID(c))
	c.JSON(200, contacts)
}

func routeDeleteContact(c *gin.Context) {
	var contact Contact
	id := c.Param("id")
	fmt.Println("[RouteDeleteContact] Id", id)
	err := Database.Where("id = ?", id).First(&contact)

	fmt.Println("[RouteDeleteContact] [debug] Struct of &contact", &contact)

	if contact.OwnerID == getUserID(c) {
		Database.Delete(&contact)
	} else {
		fmt.Println("[RouteDeleteContact] error occured", err)
		fmt.Println("[RouteDeleteContact] error occured; Id =", contact.Id)
	}

	c.JSON(200, gin.H{
		"result": "done",
	})
}

func routeCreateContact(c *gin.Context) {
	postedContact := Contact{}

	if c.BindJSON(&postedContact) == nil {
		contact := NewContact(getUserID(c), postedContact.Name)
		contact.Create()
		c.JSON(200, contact)
	} else {
		c.String(500, "Something went wrong in binding (add)")
	}
}

func routeUpdateContact(c *gin.Context) {
	contact := Contact{}
	Database.Where("id = ? AND owner_id = ?", c.Param("id"), getUserID(c)).First(&contact)

	if contact.OwnerID == getUserID(c) {
		contact.Update()
	}

	c.JSON(200, contact)
}

/*
End of Contact Routes
*/

/*
Router Profile Code
*/

func routeGetProfile(c *gin.Context) {
	user := GetUser(getUserID(c))
	c.JSON(200, user)
}

func routeUpdateProfile(c *gin.Context) {
	postedUser := User{}

	if c.BindJSON(&postedUser) == nil {
		user := GetUser(getUserID(c))

		fmt.Println("User:", user)
		fmt.Println("PostedUser:", postedUser)

		user.Name = postedUser.Name
		user.EmailAddress = postedUser.EmailAddress
		user.IBAN = postedUser.IBAN
		user.CustomCss = postedUser.CustomCss
		user.SidebarColor = postedUser.SidebarColor
		user.InvoiceMessage = postedUser.InvoiceMessage

		user.Update()
		c.JSON(200, postedUser)
	} else {
		c.String(500, "Something went wrong in binding (add)")
	}

}
func routeUpdateProfilePassword(c *gin.Context) {
	type PasswordChange struct {
		Password string
		Result   string
	}
	var change PasswordChange

	if c.BindJSON(&change) == nil {

		user := GetUser(getUserID(c))
		user.Hash = HashPassword(change.Password)
		user.Update()

		change.Result = "Succesfull!"

		//TODO: built check for minimum length

		c.JSON(200, change)
	} else {
		c.String(500, "Something went wrong in binding (add)")
	}

}

/*
Router User Code
*/

func routeGetUsers(c *gin.Context) {
	users := GetUsers()
	c.JSON(200, users)
}

func routeDeleteUser(c *gin.Context) {
	var user User
	id := c.Param("id")

	Database.Where("id = ?", id).First(&user)

	fmt.Println("[RouteDeleteUser] [debug] Struct of &user", &user)

	if GetUser(getUserID(c)).Administrator == true && user.Administrator == false {
		Database.Delete(&user)
	} else {
		fmt.Println("Not deleting user! Current user does not have administrator privileges and/or selected user is an administrator!")
	}

	c.JSON(200, gin.H{
		"result": "done",
	})
}

func routeCreateUser(c *gin.Context) {
	postedUser := User{}

	if c.BindJSON(&postedUser) == nil {
		user := User{}

		user.Username = postedUser.Username
		user.Create()
		c.JSON(200, user)
	} else {
		c.String(500, "Something went wrong in binding (add)")
	}
}

func routeUpdateUser(c *gin.Context) {
	user := User{}
	Database.Where("id = ? AND owner_id = ?", c.Param("id"), getUserID(c)).First(&user)

	user.Update()

	c.JSON(200, user)
}

/*
End of User Routes
*/

/*
Start of Debt Routes
*/

func routeGetLinkedDebts(c *gin.Context) {
	c.JSON(200, GetLinkedDebts(getUserID(c)))
}

func addDebt(c *gin.Context) {

	id := c.Param("id")
	fmt.Println("[addDebt] Id", id)
	debt := Debt{}

	var postedDebt Debt

	if c.BindJSON(&postedDebt) == nil {
		fmt.Println("AmountChecker", postedDebt.Amount)

		//change 1 in current logged in user (owner_id)
		debt.Description = postedDebt.Description
		debt.Paid = postedDebt.Paid
		fmt.Println("[binding]", debt.Description, postedDebt.Description)
		debt.NewDebt(postedDebt.Amount, getUserID(c), postedDebt.Debtor)

	} else {
		c.String(500, "Something went wrong in binding (add)")
	}

	c.JSON(200, debt)
}
func updateDebt(c *gin.Context) {
	id := c.Param("id")
	fmt.Println("[updateDebt] Id", id)
	debt := Debt{}
	Database.Where("id = ? AND owner_id = ?", id, getUserID(c)).First(&debt)

	var postedDebt Debt

	if c.BindJSON(&postedDebt) == nil {
		debt.Debtor = postedDebt.Debtor
		debt.Description = postedDebt.Description
		debt.Amount = postedDebt.Amount
		debt.Paid = postedDebt.Paid
	} else {
		c.String(500, "Something went wrong in binding")
	}

	debt.UpdateDebt()

	// returnDebt := Debt{Id: "5fsdaap", Debtor: "Dion le Maître", Description: "Henkerman"}
	//hier moet dan een update komen...
	//returnDebt.NewDebt(666, 1, "HenkStiekema3")
	c.JSON(200, postedDebt)
}

func deleteDebt(c *gin.Context) {
	var debt Debt
	id := c.Param("id")
	fmt.Println("[deleteDebt] Id", id)
	err := Database.Where("id = ?", id).First(&debt)

	fmt.Println("[deleteDebt] [debug] Struct of &debt", &debt)

	if debt.OwnerID == getUserID(c) {
		Database.Delete(&debt)
	} else {
		fmt.Println("[deleteDebt] error occured", err)
		fmt.Println("[deleteDebt] error occured; Id =", debt.Id)
	}

	c.JSON(200, gin.H{
		"result": "done",
	})
}

func setLoginUser(c *gin.Context) {
	c.String(200, "<h1>Done!</h1>")
}

//PageVariables are given to a page with information such as current user and navigation
type PageVariables struct {
	UserID     int
	User       User
	Navigation []navigationItem
}

//GetPageVariables generates and returs PageVariables data
func GetPageVariables(c *gin.Context) PageVariables {
	fmt.Println("RequestVariables")
	userID := getUserID(c)
	user := GetUser(userID)
	fmt.Println("SidebarColor:", user.SidebarColor)
	vars := PageVariables{UserID: userID, User: user, Navigation: getNavigation(user)}
	return vars
}

func splashPage(c *gin.Context) {
	c.HTML(200, "splash.html", gin.H{
		"title":            "Kreditor",
		"controller":       "MainCtl",
		"RequestVariables": GetPageVariables(c),
	})
}
func debtsPage(c *gin.Context) {
	c.HTML(200, "index.html", gin.H{
		"title":            "Kreditor",
		"controller":       "MainCtl",
		"RequestVariables": GetPageVariables(c),
		"contacts":         GetContactObjects(getUserID(c)),
	})
}
func invoices(c *gin.Context) {
	c.HTML(200, "invoices.html", gin.H{
		"title":            "Kreditor",
		"controller":       "InvoiceController",
		"RequestVariables": GetPageVariables(c),
		"contacts":         GetContactObjects(getUserID(c)),
	})
}
func contacts(c *gin.Context) {
	c.HTML(200, "contacts.html", gin.H{
		"title":            "Kreditor",
		"controller":       "ContactController",
		"RequestVariables": GetPageVariables(c),
	})
}
func profile(c *gin.Context) {
	c.HTML(200, "profile.html", gin.H{
		"title":            "Kreditor",
		"controller":       "ProfileController",
		"RequestVariables": GetPageVariables(c),
	})
}
func admin(c *gin.Context) {
	c.HTML(200, "users.html", gin.H{
		"title":            "Kreditor",
		"controller":       "UserController",
		"RequestVariables": GetPageVariables(c),
	})
}
