package db

import (
	"database/sql"
	"log"
	_ "github.com/go-sql-driver/mysql"
	"fmt"
	"strconv"
	"strings"
	"os"
)

var (
    Dbhost = os.Getenv("DBHOST")
    Dbuser = os.Getenv("DBUSER")
    Dbpassword = os.Getenv("DBPASSWORD")
    Dbname = os.Getenv("DBNAME")
)

func InsertDb(linkid string, links []struct{Date string "json:\"date\""; Clicks int "json:\"clicks\""}) {
    log.Printf("Going to insert link %s data into Mysql", linkid)
    log.Printf("Data to insert %+v", links)
	db, err := sql.Open("mysql", Dbuser + ":" + Dbpassword + "@tcp(" + Dbhost + ":3306)/" + Dbname)
    if err != nil {
		log.Panic("Impossible to create the connection to Mysql: %s", err)
    }
    defer db.Close()
	date := strings.Split(links[0].Date, "T")
	clicks := strconv.Itoa(links[0].Clicks)

	query := "INSERT INTO `links` (`linkid`, `date`, `clicks`) VALUES ('" + linkid + "', '" + date[0] + "', " + clicks + ")"
	fmt.Println(query)
	insertResult, err := db.Exec(query)
	if err !=nil {
		panic(err.Error())
	}
	id, err := insertResult.LastInsertId()
	if err != nil {
		log.Fatalf("impossible to retrieve last inserted id: %s", err)
	}
	log.Printf("inserted id: %d", id)
}
