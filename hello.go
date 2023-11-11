package main

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	_ "github.com/joho/godotenv/autoload"
)

type CreateAlbumDTO struct {
	Title  string  `json:"title"`
	Artist string  `json:"artist"`
	Price  float64 `json:"price"`
}

type Album struct {
	ID     int     `json:"id"`
	Title  string  `json:"title"`
	Artist string  `json:"artist"`
	Price  float64 `json:"price"`
}

var db *pgx.Conn

func getAlbumsByArtist(name string) ([]Album, error) {
	var albums []Album

	rows, err := db.Query(context.TODO(), "SELECT * FROM album WHERE artist=$1", name)

	if err != nil {
		return nil, fmt.Errorf("getAlbumsByArtist %q: %v", name, err)
	}

	defer rows.Close()

	for rows.Next() {
		var alb Album

		if err := rows.Scan(&alb.ID, &alb.Title, &alb.Artist, &alb.Price); err != nil {
			return nil, fmt.Errorf("getAlbumsByArtist %q: %v", name, err)
		}

		albums = append(albums, alb)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("getAlbumsByArtist %q: %v", name, err)
	}

	return albums, nil
}

func handleAlbumByArtist(c *gin.Context) {
	name := c.Param("artist")
	albums, _ := getAlbumsByArtist(name)

	c.IndentedJSON(http.StatusOK, albums)
}

func addAlbum(createAlbumDTO CreateAlbumDTO) error {
	_, err := db.Exec(context.TODO(), "INSERT INTO album (title, artist, price) VALUES ($1, $2, $3)", createAlbumDTO.Title, createAlbumDTO.Artist, createAlbumDTO.Price)

	return err
}

func handleAddAlbum(c *gin.Context) {
	var createAlbumDTO CreateAlbumDTO

	if err := c.BindJSON(&createAlbumDTO); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Invalid params"})
	}
	if err := addAlbum(createAlbumDTO); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Invalid params"})
		return
	}

	c.Status(http.StatusCreated)
}

func getAlbumById(id int) (Album, error) {
	var album Album
	row := db.QueryRow(context.TODO(), "SELECT * FROM album WHERE id = $1", id)

	if scanErr := row.Scan(&album.ID, &album.Title, &album.Artist, &album.Price); scanErr != nil {
		if scanErr == sql.ErrNoRows {
			return album, fmt.Errorf("getAlbumById: Album not found for id %v", id)
		}

		return album, fmt.Errorf("getAlbumById: %v, %v", id, scanErr)
	}

	return album, nil
}

func handleAlbumById(c *gin.Context) {
	id, paramConversionError := strconv.Atoi(c.Param("id"))

	if paramConversionError != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Invalid params"})
		return
	}

	album, albumByIdError := getAlbumById(id)

	if albumByIdError != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": albumByIdError.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, album)
}

func getAlbums() ([]Album, error) {
	var albums = []Album{}
	rows, err := db.Query(context.Background(), "SELECT * FROM album")

	if err != nil {
		return nil, fmt.Errorf("getAlbums %v", err)
	}

	defer rows.Close()

	for rows.Next() {
		var alb Album

		if err := rows.Scan(&alb.ID, &alb.Title, &alb.Artist, &alb.Price); err != nil {
			return nil, fmt.Errorf("getAlbums %v", err)
		}

		albums = append(albums, alb)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("getAlbums %v", err)
	}

	return albums, nil
}

func handleAlbums(c *gin.Context) {
	albums, err := getAlbums()

	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Unexpected error"})
		return
	}

	c.IndentedJSON(http.StatusOK, albums)
}

func main() {
	fmt.Println(os.Getenv("DB_CONN"))

	var err error
	db, err = pgx.Connect(context.Background(), os.Getenv("DB_CONN"))

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to the database %v\n", err)
		os.Exit(1)
	}
	if pingErr := db.Ping(context.Background()); pingErr != nil {
		fmt.Fprintf(os.Stderr, "Unable to ping the database %v\n", pingErr)
		os.Exit(1)
	}
	defer db.Close(context.Background())

	router := gin.Default()

	router.GET("/albums", handleAlbums)
	router.POST("/albums", handleAddAlbum)
	router.GET("/albums/artist/:artist", handleAlbumByArtist)
	router.GET("/albums/:id", handleAlbumById)

	router.Run("localhost:8080")
}
