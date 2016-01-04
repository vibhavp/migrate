package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"

	_ "github.com/lib/pq"
)

var (
	dbHost     = flag.String("host", "localhost", "DB Host")
	dbPort     = flag.Int("port", 0, "port")
	dbName     = flag.String("database", "", "Database name")
	dbUser     = flag.String("user", "", "DB Username")
	dbPassword = flag.String("password", "", "DB Password")

	connect string
)

func init() {
	flag.Parse()
	if *dbPort == 0 || *dbName == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	connect = fmt.Sprintf("host=%s port=%d dbname=%s user=%s password=%s sslmode=disable",
		*dbHost, *dbPort, *dbName, *dbUser, *dbPassword)
}

func main() {
	db, err := sql.Open("postgres", connect)
	if err != nil {
		log.Fatal(err)
	}

	rows, err := db.Query("SELECT whitelist, id FROM lobbies")
	if err != nil {
		log.Fatal(err)
	}

	var oldWhitelists []int
	var lobbyIDs []int

	for rows.Next() {
		var whitelistID uint
		var lobbyID uint

		rows.Scan(&whitelistID, &lobbyID)

		lobbyIDs = append(lobbyIDs, int(lobbyID))
		oldWhitelists = append(oldWhitelists, int(whitelistID))
	}

	db.Exec("ALTER TABLE lobbies DROP COLUMN whitelist")
	db.Exec("ALTER TABLE lobbies ADD whitelist varchar(255)")

	for i, lobbyID := range lobbyIDs {
		log.Printf("UPDATE lobbies SET whitelist = %s WHERE id = %d\n",
			strconv.Itoa(oldWhitelists[i]),
			lobbyID)

		_, err := db.Exec("UPDATE lobbies SET whitelist = $1 WHERE id = $2", strconv.Itoa(oldWhitelists[i]), lobbyID)
		if err != nil {
			log.Fatal(err)
		}
	}

	db.Close()
	log.Println("done!")
}
