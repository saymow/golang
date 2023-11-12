package main

import (
	"context"
	"fmt"
	"os"

	"com.saymow/controllers"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	_ "github.com/joho/godotenv/autoload"
)

var db *pgx.Conn

func main() {
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

	router.GET("/albums", func(ctx *gin.Context) { controllers.HandleAlbums(ctx, db) })
	router.POST("/albums", func(ctx *gin.Context) { controllers.HandleAddAlbum(ctx, db) })
	router.GET("/albums/artist/:artist", func(ctx *gin.Context) { controllers.HandleAlbumByArtist(ctx, db) })
	router.GET("/albums/:id", func(ctx *gin.Context) { controllers.HandleAlbumById(ctx, db) })

	router.Run("localhost:8080")
}
