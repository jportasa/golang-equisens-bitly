package main

// Host: db4free.net 
// DB: bitlyequisens


import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	//"github.com/aws/aws-lambda-go/lambda"
	"database/sql"
	"fmt"
	"strconv"
	"os"
	_ "github.com/go-sql-driver/mysql"
    "time"
    "net/url"
)

type link struct {
	UnitReference string `json:"unit_reference"`
	LinkClicks []struct {
		Date string `json:"date"`
		Clicks int `json:"clicks"`
	} `json:"link_clicks"`
	Units int `json:"units"`
	Unit string `json:"unit"`
}

type linkids struct {
	Links []struct {
		Id string `json:"id"`
	} `json:"links"`
}

var (
    token = os.Getenv("BITLYTOKEN")
    password = os.Getenv("DBPASSWORD")
)

func HandleRequest() {
	for _, linkid := range GetLinkIds(){
		client := &http.Client{}
		req, err := http.NewRequest("GET", "https://api-ssl.bitly.com/v4/bitlinks/bit.ly/" + linkid + "/clicks?unit=month&units=1", nil) // Gets ALL links of any date
		if err != nil {
			log.Fatal(err)
		}
		req.Header.Set("Authorization", "Bearer " + token)
		resp, err := client.Do(req)
		if err != nil {
			log.Fatal(err)
		}
		defer resp.Body.Close()
		bodyText, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}
		links := link{}
		err = json.Unmarshal(bodyText, &links)
		if err != nil {
			log.Panic("error:", err)
		}
		insertDb(linkid, links.LinkClicks)
	}
}

func GetOldestDate () string {
    // Returns the date of two years ago from now
    theTime := time.Now()
    toAdd := -17280 * time.Hour
    newTime := theTime.Add(toAdd)
    log.Println("Oldest date to get links indexes is", newTime)
    log.Println("Formated date is", newTime.Format("2006-01-02T15:04:05-0700"))
    return url.QueryEscape(newTime.Format("2006-01-02T15:04:05-0700"))
}

func GetLinkIds() []string {
    // Outputs the list of bitly link ID's
	var listIds []string
	ids := linkids{}

	client := &http.Client{}
    reqLink := "https://api-ssl.bitly.com/v4/groups/Bj156NPShb0/bitlinks/clicks?" + GetOldestDate()
    log.Printf("link to get list of links id is %s", reqLink)
	req, err := http.NewRequest("GET", reqLink , nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Authorization", "Bearer " + token)
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	bodyText, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	err = json.Unmarshal(bodyText, &ids)
	if err != nil {
		log.Panic("error:", err)
	}
	for _, id := range ids.Links{
		item := strings.Split(id.Id, "/")
		listIds = append(listIds, item[1])
	}
	log.Println("list of link id's", listIds)
	return listIds
}

func insertDb(linkid string, links []struct{Date string "json:\"date\""; Clicks int "json:\"clicks\""}) {
    log.Printf("Going to insert link %s data into Mysql", linkid)
    log.Printf("Data to insert %+v", links)
	db, err := sql.Open("mysql", "userequisens:" + password + "@tcp(db4free.net:3306)/bitlyequisens")
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

func main() {
    //lambda.Start(HandleRequest)
    HandleRequest()
}
