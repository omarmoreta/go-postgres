package main

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5"
)

type Album struct {
    ID     int64
    Title  string
    Artist string
    Price  float32
}

func main() {
	conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(context.Background())

	albums, err := albumsByArtists("Gerry Mulligan", conn)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to get albums: %v\n", err)
		return
	}

	for _, album := range albums {
		fmt.Printf("ID: %d, Title: %s, Artist: %s, Price: %f\n", album.ID, album.Title, album.Artist, album.Price)
	}
}

func albumsByArtists(name string, conn *pgx.Conn) ([]Album, error) {
	rows, err := conn.Query(context.Background(), "select * from album where artist = $1", name)
	if err != nil {
		return nil, fmt.Errorf("albumByArtist %q: %v", name, err)
	}
	defer rows.Close()

	var albums []Album
	for rows.Next() {
		var alb Album
		if err := rows.Scan(&alb.ID, &alb.Title, &alb.Artist, &alb.Price); err != nil {
			return nil, fmt.Errorf("albumByArtist %q: %v", name, err)
		}
		albums = append(albums, alb)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("albumByArtist %q: %v", name, err)
	}

	return albums, nil
}