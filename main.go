package main

import (
	"fmt"
	"net/http"

	"github.com/joho/godotenv"
	"github.com/gin-gonic/gin"
)

// album represents data about a record album.
type album struct {
	ID     string  `json:"id"`
	Title  string  `json:"title"`
	Artist string  `json:"artist"`
	Price  float64 `json:"price"`
}

// albums slice to seed record album data.
var albums = []album{
	{ID: "1", Title: "The Rise and Fall of a Midwest Princess", Artist: "Chappell Roan", Price: 32.49},
	{ID: "2", Title: "SOS", Artist: "SZA", Price: 34.98},
	{ID: "3", Title: "World Wide Whack", Artist: "Tierra Whack", Price: 29.98},
}

func main() {
	fmt.Printf("Hello world!")

	// Load environment variables from .env file
	err := godotenv.Load()
	handle(err)
	
	// Setup DB connection
	client := connectCosmosClient()
	h:= New(client)
	
	// Initialize gin router
	router := gin.Default()
	// Note: we're passing the *name* of the function, not the function itself which would be getAlbums()
	router.GET("/albums", h.getAlbums)
	router.GET("/albums/:id", h.getAlbumByID)
	router.POST("/albums", h.postAlbums)

	// Attach router to an HTTP server
	router.Run("localhost:8080")
}

// HANDLERS
// getAlbums responds with the list of all albums as JSON.
func (h handler) getAlbums(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, albums) // Could use JSON but this format doesn't take much extra space and is easier to read
}

// postAlbums adds an album from JSON received in the request body.
func (h handler) postAlbums(c *gin.Context) {
	var newAlbum album

	// Call BindJSON to bind the received JSON to newAlbum.
	if err := c.BindJSON(&newAlbum); err != nil {
		return
	}

	// Add the new album to the slice.
	albums = append(albums, newAlbum)
	c.IndentedJSON(http.StatusCreated, newAlbum)
}

// getAlbumByID locates the album whose ID value matches the id
// parameter sent by the client, then returns that album as a response.
func (h handler) getAlbumByID(c *gin.Context) {
	id := c.Param("id")

	container := getContainer(*h.client)

	album := getAlbumByIdFromCosmos(*container, id)

	if album != nil {
		c.IndentedJSON(http.StatusOK, album)
		return
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"Error": "Album not found."})
}
