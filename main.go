package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	//"github.com/aws/aws-lambda-go/lambda"
	"os"
    "time"
    "net/url"
	"bitlyequisens/db"
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
    Token = os.Getenv("BITLYTOKEN")
)

func HandleRequest() {
	for _, linkid := range GetLinkIds(){
		client := &http.Client{}
		req, err := http.NewRequest("GET", "https://api-ssl.bitly.com/v4/bitlinks/bit.ly/" + linkid + "/clicks?unit=month&units=1", nil) // Gets ALL links of any date
		if err != nil {
			log.Fatal(err)
		}
		req.Header.Set("Authorization", "Bearer " + Token)
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
		db.InsertDb(linkid, links.LinkClicks)
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
	req.Header.Set("Authorization", "Bearer " + Token)
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

func main() {
	//lambda.Start(HandleRequest)
    HandleRequest()
}
