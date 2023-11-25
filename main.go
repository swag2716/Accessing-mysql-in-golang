package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/go-sql-driver/mysql"
)

type Album struct {
	ID     int64
	Title  string
	Artist string
	Price  float32
}

var db *sql.DB

func main() {
	cfg := mysql.Config{
		User:   os.Getenv("DBUSER"),
		Passwd: os.Getenv("DBPASS"),
		Net:    "tcp",
		Addr:   "127.0.0.1:3306",
		DBName: "recordings",
	}
	// cfg2 := mysql.Config{
	// 	User:   os.Getenv("DBUSER"),
	// 	Passwd: os.Getenv("DBPASS"),
	// 	DBName: "recordings2",
	// }

	var err error

	db, err = sql.Open("mysql", cfg.FormatDSN())

	if err != nil {
		log.Fatal(err)
	}

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Sql connected!")

	albums, err := albumsbyArtist("Swapnika")

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Albums found are : ", albums)

	album, err := albumById(5)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Album found by id is : ", album)
}

func albumsbyArtist(name string) ([]Album, error) {
	var albums []Album

	rows, err := db.Query("SELECT * FROM album WHERE artist = ?", name)

	if err != nil {
		return nil, fmt.Errorf("%s, %v", name, err)
	}

	defer rows.Close()

	for rows.Next() {
		var alb Album
		err := rows.Scan(&alb.ID, &alb.Title, &alb.Artist, &alb.Price)

		if err != nil {
			return nil, fmt.Errorf("%s, %v", name, err)
		}

		albums = append(albums, alb)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("albumsByArtist %q: %v", name, err)
	}
	return albums, nil
}

func albumById(id int64) (Album, error) {
	row := db.QueryRow("SELECT * FROM album WHERE id = ?", id)

	var alb Album

	err := row.Scan(&alb.ID, &alb.Title, &alb.Artist, &alb.Price)
	if err != nil {
		if err == sql.ErrNoRows {
			return alb, fmt.Errorf("albumsById %d: no such album", id)
		}
		return alb, fmt.Errorf("albumsByArtist %q: %v", id, err)
	}

	return alb, nil
}
