package services

import (
	"context"
	"database/sql"
	"fmt"

	"com.saymow/structs/album"
	"github.com/jackc/pgx/v5"
	_ "github.com/joho/godotenv/autoload"
)

func GetAlbumsByArtist(db *pgx.Conn, name string) ([]album.Album, error) {
	var albums []album.Album

	rows, err := db.Query(context.TODO(), "SELECT * FROM album WHERE artist=$1", name)

	if err != nil {
		return nil, fmt.Errorf("getAlbumsByArtist %q: %v", name, err)
	}

	defer rows.Close()

	for rows.Next() {
		var alb album.Album

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

func AddAlbum(db *pgx.Conn, createAlbumDTO album.CreateAlbumDTO) error {
	_, err := db.Exec(context.TODO(), "INSERT INTO album (title, artist, price) VALUES ($1, $2, $3)", createAlbumDTO.Title, createAlbumDTO.Artist, createAlbumDTO.Price)

	return err
}

func GetAlbumById(db *pgx.Conn, id int) (album.Album, error) {
	var album album.Album
	row := db.QueryRow(context.TODO(), "SELECT * FROM album WHERE id = $1", id)

	if scanErr := row.Scan(&album.ID, &album.Title, &album.Artist, &album.Price); scanErr != nil {
		if scanErr == sql.ErrNoRows {
			return album, fmt.Errorf("getAlbumById: Album not found for id %v", id)
		}

		return album, fmt.Errorf("getAlbumById: %v, %v", id, scanErr)
	}

	return album, nil
}

func GetAlbums(db *pgx.Conn) ([]album.Album, error) {
	var albums = []album.Album{}
	rows, err := db.Query(context.Background(), "SELECT * FROM album")

	if err != nil {
		return nil, fmt.Errorf("getAlbums %v", err)
	}

	defer rows.Close()

	for rows.Next() {
		var alb album.Album

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
