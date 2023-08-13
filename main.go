package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"github.com/aws/aws-lambda-go/lambda"
	"database/sql"
	"fmt"
	"strconv"
	"os"
	_ "github.com/go-sql-driver/mysql"
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

func HandleRequest() {
	//linkid := "3N7jxxl"
	for _, linkid := range GetLinkIds(){
		token := os.Getenv("BITLYTOKEN")
		client := &http.Client{}
		req, err := http.NewRequest("GET", "https://api-ssl.bitly.com/v4/bitlinks/bit.ly/" + linkid + "/clicks?unit=day&units=2", nil) // TODO, he d'agafar els clicks de tot el mes
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

func GetLinkIds() []string {
	token := os.Getenv("BITLYTOKEN")
	var listIds []string
	ids := linkids{}

	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://api-ssl.bitly.com/v4/groups/Bj156NPShb0/bitlinks/clicks?unit=month&units=1&unit_reference=2023-08-02T15%3A04%3A05-0700&size=10", nil)
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
	fmt.Println(listIds)
	return listIds
}

func insertDb(linkid string, links []struct{Date string "json:\"date\""; Clicks int "json:\"clicks\""}) {
	password := os.Getenv("DBPASSWORD")
	db, err := sql.Open("mysql", "userequisens:" + password + "@tcp(db4free.net:3306)/bitlyequisens")
    if err != nil {
		log.Panic("impossible to create the connection: %s", err)
    }
    defer db.Close()
	fmt.Printf("%+v\n", links[1].Date)
	date := strings.Split(links[1].Date, "T")
	clicks := strconv.Itoa(links[1].Clicks)

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
	lambda.Start(HandleRequest)
}