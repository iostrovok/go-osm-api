package main

import (
	"flag"
	"github.com/iostrovok/go-osm-api/osmapi"
	"log"
	"os"
)

func CreateNode(req *osmapi.MyRequestSt, Lat, Lon string) {

	if Lat == "" || Lon == "" {
		log.Fatal("No setup Lat or Lon param")
	}

	ChSet, err := req.Changesets()
	if err != nil {
		log.Fatal(err)
	}

	if err := ChSet.OsmChange("create"); err != nil {
		log.Fatal(err)
	}

	/* Create new node */
	node, err_n := ChSet.NewNode(Lat, Lon)
	if err_n != nil {
		log.Fatal(err_n)
	}

	/*  Set new data */
	node.AddTag("name:en", "Anywhere in London")
	node.AddTag("name:ru", "Поселок где-то в лондоне")
	node.AddTag("name:uk", "Поселок где-то в лондоне")
	node.AddTag("place", "street")

	/* Upload new data */
	if newId, err := ChSet.Upload(); err != nil {
		log.Fatal(err)
	} else {
		log.Printf("Create node. Id; %s\n", newId)
	}

	if err := ChSet.Close(); err != nil {
		log.Fatal(err)
	}
}

func main() {
	login := os.Getenv("OSM_USER")
	pass := os.Getenv("OSM_PASSWD")
	url := os.Getenv("OSM_URL")

	lat := flag.String("lat", "", "Latitude +90 is the North Pole; latitude -90 is the South Pole.")
	lon := flag.String("lon", "", "The maximum and minimum longitudes (+180 and -180) are along the same north-south line through the middle of the Pacific Ocean.")

	flag.Parse()

	if login == "" || pass == "" {
		log.Fatalln("Login and password are not found.")
	}

	req := osmapi.MyRequest(login, pass)
	if url == "" {
		url = osmapi.MainURLTest
	}
	req.SetUrl(url)

	CreateNode(req, *lat, *lon)
}
