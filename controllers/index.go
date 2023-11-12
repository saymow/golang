package controllers

import (
	"net/http"
	"strconv"

	"com.saymow/services"
	"com.saymow/structs/album"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	_ "github.com/joho/godotenv/autoload"
)

func HandleAlbumByArtist(c *gin.Context, db *pgx.Conn) {
	name := c.Param("artist")
	albums, _ := services.GetAlbumsByArtist(db, name)

	c.IndentedJSON(http.StatusOK, albums)
}

func HandleAddAlbum(c *gin.Context, db *pgx.Conn) {
	var createAlbumDTO album.CreateAlbumDTO

	if err := c.BindJSON(&createAlbumDTO); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Invalid params"})
	}
	if err := services.AddAlbum(db, createAlbumDTO); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Invalid params"})
		return
	}

	c.Status(http.StatusCreated)
}

func HandleAlbumById(c *gin.Context, db *pgx.Conn) {
	id, paramConversionError := strconv.Atoi(c.Param("id"))

	if paramConversionError != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Invalid params"})
		return
	}

	album, albumByIdError := services.GetAlbumById(db, id)

	if albumByIdError != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": albumByIdError.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, album)
}

func HandleAlbums(c *gin.Context, db *pgx.Conn) {
	albums, err := services.GetAlbums(db)

	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Unexpected error"})
		return
	}

	c.IndentedJSON(http.StatusOK, albums)
}
