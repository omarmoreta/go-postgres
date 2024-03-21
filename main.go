package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

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

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter artist's name to search for albums: ")
	artistName, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error reading input:", err)
		return
	}
	artistName = strings.TrimSpace(artistName)

	albums, err := albumsByArtists(artistName, conn)
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

	fmt.Println(albums)
	return albums, nil
}
