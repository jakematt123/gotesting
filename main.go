package main

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type book struct {
	ID       string "json:id"
	Title    string "json:title"
	Author   string "json:author"
	Quantity int    "json:quantity"
}

var books = []book{
	{ID: "1", Title: "The Hitchhiker's Guide to the Galaxy", Author: "Douglas Adams", Quantity: 5},
	{ID: "2", Title: "The Hobbit", Author: "J.R.R. Tolkien", Quantity: 3},
	{ID: "3", Title: "The Lord of the Rings", Author: "J.R.R. Tolkien", Quantity: 2},
}

func getBooks(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, books)
}

func homepage(c *gin.Context) {
	http.ServeFile(c.Writer, c.Request, "index.html")
	c.Status(http.StatusOK)
}

func newBook(c *gin.Context) {
	var newbook book

	if err := c.BindJSON(&newbook); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Invalid request"})
		return
	}

	books = append(books, newbook)
	c.IndentedJSON(http.StatusCreated, newbook)

}

func bookById(c *gin.Context) {
	id := c.Param("id")
	book, err := getBookByID(id)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Book not found"})
		return
	}

	c.IndentedJSON(http.StatusOK, book)
}

func getBookByID(id string) (*book, error) {
	for i, book := range books {
		if book.ID == id {
			return &books[i], nil
		}
	}

	return nil, errors.New("Book not found")
}

func checkoutBook(c *gin.Context) {
	id, ok := c.GetQuery("id")

	if ok == false {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Invalid request"})
		return
	}

	book, err := getBookByID(id)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Book not found"})
		return
	}

	if book.Quantity <= 0 {
		c.IndentedJSON(http.StatusConflict, gin.H{"message": "Book not available"})
		return
	}

	book.Quantity--
	c.IndentedJSON(http.StatusOK, book)

}

func main() {
	router := gin.Default()
	router.GET("/books", getBooks)
	router.GET("/", homepage)
	router.POST("/newbook", newBook)
	router.GET("/books/:id", bookById)
	router.PATCH("/checkout", checkoutBook)
	router.Run("localhost:8080")

}
