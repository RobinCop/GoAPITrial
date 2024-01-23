package main

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type book struct {
	ID       string `json:"id"`
	Title    string `json:"title"`
	Author   string `json:"author"`
	Quantity int    `json:"quantity"`
}

var books = []book{
	{ID: "1", Title: "In Search of Lost Time", Author: "Marcel Proust", Quantity: 2},
	{ID: "2", Title: "The Great Gatsby", Author: "F. Scott Fitzgerald", Quantity: 5},
	{ID: "3", Title: "War and Peace", Author: "Leo Tolstoy", Quantity: 6},
}

func getBooks(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, books)
}

func bookById(c *gin.Context) {
	id := c.Param("id")
	book, err := getBookById(id)

	if err != nil {
		c.IndentedJSON(http.StatusNotFound, fmt.Sprintf("book with id: %v not found", id))
		return
	}
	c.IndentedJSON(http.StatusOK, book)
}

func checkoutOrReturnBook(c *gin.Context, checkout bool) {
	id, ok := c.GetQuery("id")
	if !ok {
		c.IndentedJSON(http.StatusBadRequest, fmt.Sprint("Missing id"))
		return
	}

	book, err := getBookById(id)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, fmt.Sprint("Book not found"))
		return
	}

	if checkout {
		if book.Quantity > 0 {
			book.Quantity -= 1
			c.IndentedJSON(http.StatusOK, book)
		} else {
			c.IndentedJSON(http.StatusInsufficientStorage, fmt.Sprint("Not enough books remaining"))
		}
	} else {
		book.Quantity += 1
		c.IndentedJSON(http.StatusOK, book)
	}
}

func checkoutBook(c *gin.Context) {
	checkoutOrReturnBook(c, true)
}

func returnBook(c *gin.Context) {
	checkoutOrReturnBook(c, false)
}

func getBookById(id string) (*book, error) {
	for i, b := range books {
		if b.ID == id {
			return &books[i], nil
		}
	}

	return nil, errors.New("book not found")
}

func createBook(c *gin.Context) {
	var newBook book

	if err := c.BindJSON(&newBook); err != nil {
		return
	}

	books = append(books, newBook)
	c.IndentedJSON(http.StatusCreated, newBook)
}

func main() {
	router := gin.Default()
	router.GET("/books", getBooks)          //curl localhost:8080/books
	router.GET("/books/:id", bookById)      //curl localhost:8080/books/2
	router.POST("/books", createBook)       //curl localhost:8080/books --include --header "Content-Type: application/json" -d @body.json --request "POST"
	router.PATCH("/checkout", checkoutBook) //curl localhost:8080/checkout?id=2 --request "PATCH"
	router.PATCH("/return", returnBook)     //curl localhost:8080/return?id=2 --request "PATCH"
	router.Run("localhost:8080")            //Random git test commit

}
