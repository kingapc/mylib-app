package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	entities "github.com/rpinedafocus/mylib-entities"
	manage "github.com/rpinedafocus/mylib-entities/manage"
	security "github.com/rpinedafocus/mylib-entities/security"
)

type Credentials struct {
	USER     string `json:"user"`
	PASSWORD string `json:"password"`
}

type InfoLogin struct {
	FULL_NAME string
	USER_NAME string
	ROLE_ID   int
}

var isLogin1 InfoLogin

/**********	MAIN	************/
func main() {
	router := gin.Default()
	router.GET("/authors", getAuthors)
	router.GET("/books", getBooks)
	router.GET("/books/:id", getBook)
	router.GET("/login", login)
	router.POST("/manage/rent", setRent)
	router.POST("/manage/reserve", setReserve)

	router.Run("localhost:8080")
}

/*********** END POINTS	***********/
func getAuthors(c *gin.Context) {

	authors := entities.GetRows()
	fmt.Print(authors)
	c.IndentedJSON(http.StatusOK, authors)
}

//Get all books
func getBooks(c *gin.Context) {

	if isLogin1.ROLE_ID != 0 {
		isAccess := security.GetAccess(isLogin1.ROLE_ID, c.Request.URL.Path)

		if isAccess {

			books := entities.GetBooks()
			c.IndentedJSON(http.StatusOK, books)
		} else {
			c.IndentedJSON(http.StatusForbidden, gin.H{"message": "You don't have access to this option"})
		}
	} else {
		c.IndentedJSON(http.StatusForbidden, gin.H{"message": "You must log in before to access."})
	}
}

//Rent the book
func setRent(c *gin.Context) {

	if isLogin1.ROLE_ID != 0 {
		isAccess := security.GetAccess(isLogin1.ROLE_ID, c.Request.URL.Path)

		if isAccess {
			var nr manage.Status

			if err := c.BindJSON(&nr); err != nil {
				return
			}

			result := manage.GetRent(nr)

			if result {
				c.IndentedJSON(http.StatusCreated, gin.H{"message": "Rent successfully created"})
			} else {
				c.IndentedJSON(http.StatusCreated, gin.H{"message": "Unable to rent the book"})
			}
		} else {
			c.IndentedJSON(http.StatusForbidden, gin.H{"message": "You don't have access to this option"})
		}
	} else {
		c.IndentedJSON(http.StatusForbidden, gin.H{"message": "You must log in before to access."})
	}
}

//Reserve the book
func setReserve(c *gin.Context) {

	if isLogin1.ROLE_ID != 0 {
		isAccess := security.GetAccess(isLogin1.ROLE_ID, c.Request.URL.Path)

		if isAccess {

			books := entities.GetBooks()
			c.IndentedJSON(http.StatusOK, books)
		} else {
			c.IndentedJSON(http.StatusForbidden, gin.H{"message": "You don't have access to this option"})
		}
	} else {
		c.IndentedJSON(http.StatusForbidden, gin.H{"message": "You must log in before to access."})
	}
}

func getBook(c *gin.Context) {

	if isLogin1.ROLE_ID != 0 {
		isAccess := security.GetAccess(isLogin1.ROLE_ID, "/books/:id")

		if isAccess {
			var b entities.Books

			id, fail := strconv.ParseUint(c.Param("id"), 10, 32)

			if fail != nil {
				c.IndentedJSON(http.StatusForbidden, gin.H{"message": "Unable to parse the id"})
			}

			b, err := entities.GetBook(int(id))

			if !err {
				c.IndentedJSON(http.StatusForbidden, gin.H{"message": "Unable to get the detail"})
			} else {
				c.IndentedJSON(http.StatusOK, b)
			}
		} else {
			c.IndentedJSON(http.StatusForbidden, gin.H{"message": "You don't have access to this option"})
		}
	} else {
		c.IndentedJSON(http.StatusForbidden, gin.H{"message": "You must log in before to access."})
	}
}

//Login
func login(c *gin.Context) {
	var credentials Credentials

	if err := c.BindJSON(&credentials); err != nil {
		return
	}

	isLogin, err := security.Login(credentials.USER, credentials.PASSWORD)
	isLogin1.FULL_NAME = isLogin.FULL_NAME
	isLogin1.USER_NAME = isLogin.FULL_NAME
	isLogin1.ROLE_ID = isLogin.ROLE_ID

	if err {
		c.IndentedJSON(http.StatusForbidden, gin.H{"message": "Unable to login"})
	} else {
		c.IndentedJSON(http.StatusOK, isLogin)
	}
}
