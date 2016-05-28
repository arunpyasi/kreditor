package main

import (
	"fmt"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/mdeheij/gin-csrf"
	"github.com/mdeheij/kreditor/config"
	"golang.org/x/crypto/bcrypt"
	"strings"
)

// CsrfOptions stores the options to use for CSRF protection.
var CsrfOptions csrf.Options

//AuthRequired is authentication middleware for user authenticaton.
func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		fmt.Println("Aahh AUTHREQUIRED")
		username := getLoginUsername(c)

		//_, err := getUserByUsername(username)
		userID := getUserID(c)

		fmt.Println("userID is nu ", userID)
		fmt.Println("username is nu ", username)

		if userID > 0 {
			c.Next()
		} else {
			c.Redirect(302, "/auth/?pleaseloginfirst")
			c.Abort()
		}
	}
}

// HashPassword hashes the provided password.
// The hashing uses bcrypt with a default cost of 10.
func HashPassword(password string) string {
	// Hashing the password with the default cost of 10
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}

	return string(hashedPassword)
}

// CheckPassword compares whether or not the provided password matches the saved password hash for the provided username.
func CheckPassword(username string, password string) (bool, User) {
	user := GetUserByUsername(username)

	//if err == nil {
	compareErr := bcrypt.CompareHashAndPassword([]byte(user.Hash), []byte(password))

	if compareErr != nil {
		fmt.Println("compareErr")
		fmt.Println(compareErr)
		return false, user
	} else {
		fmt.Println("NO compareErr")
		return true, user
	}

	//}

	//	fmt.Println(err)
	return false, User{}
}

func getUserID(c *gin.Context) int {
	fmt.Println("getUserID before session")
	session := sessions.Default(c)
	fmt.Println("getUserID after session")

	if session.Get("userID") != nil {
		return session.Get("userID").(int)
	} else {
		fmt.Println("Session is wel nil, dus 0 returnen")
		return 0
	}
}

func getLoginUsername(c *gin.Context) string {
	session := sessions.Default(c)
	username := session.Get("username")

	if username != nil {
		return strings.ToLower(username.(string))
	}

	return ""
}

func loginInit(r *gin.Engine) {
	//TODO: move this component to configuration file
	var secret string

	secret = config.C.Secret

	var CookieConfig sessions.Options
	CookieConfig.Path = "/"
	CookieConfig.MaxAge = 86400 * 30

	CsrfOptions = csrf.Options{
		Secret: secret,
		ErrorFunc: func(c *gin.Context) {
			c.String(400, "Please try again later")
			c.Abort()
		},
	}

	store := sessions.NewCookieStore([]byte(secret))
	store.Options(CookieConfig)

	r.Use(sessions.Sessions("kreditor", store))
	r.Use(csrf.Middleware(CsrfOptions))

	group := r.Group("/auth")
	{
		group.GET("/", loginPage)
		group.GET("/logout", logout)
		group.GET("/session", getSession)
		group.POST("/", func(c *gin.Context) {
			username := strings.ToLower(c.PostForm("username"))
			password := c.PostForm("password")

			check, user := CheckPassword(username, password)

			if check {
				session := sessions.Default(c)
				fmt.Println("Setting username " + username)
				session.Set("userID", user.Id)
				session.Save()
				fmt.Println("Creating session for " + getLoginUsername(c))
				c.Redirect(302, "/debts")
			} else {
				fmt.Println("Invalid credentials!")
				fmt.Println("If this is a new password, please set the following hash:", HashPassword(password))
				c.Redirect(302, "/auth/?invalidpassword")
			}
		})
	}
}

// getUserByUsername returns the user details if the username matches with the known usernames.
// An empty result and an error will be returned if the user is not found.
// func getUserByUsername(username string) (User, error) {
//
// 	var users []User
//
// 	users = append(users, CurrentUser)
//
// 	for _, u := range users {
// 		if u.Username == username {
// 			fmt.Println("User obj", u, "for", username)
// 			return u, nil
// 		}
// 	}
//
// 	return User{}, fmt.Errorf("Cannot find user %s", username)
// }

func loginPage(c *gin.Context) {
	// fmt.Println(getLoginUsername(c))

	//Clear session anyway on login page
	destroySession(c)
	c.HTML(200, "login.html", gin.H{
		"user": getLoginUsername(c),
		"csrf": csrf.GetToken(c),
	})
}

func destroySession(c *gin.Context) {
	session := sessions.Default(c)
	session.Set("userID", nil)
	session.Save()
	session.Clear()
}

func logout(c *gin.Context) {
	destroySession(c)
	c.Redirect(302, "/")
}

func getSession(c *gin.Context) {
	c.JSON(200, gin.H{
		"username": getLoginUsername(c),
		"userID":   getUserID(c),
	})
}
