package main

import (
	"flag"
	"github.com/iostrovok/go-osm-api/osmapi"
	"log"
	"os"
	"strconv"
)

func DeleteNode(req *osmapi.MyRequestSt, NodeId int) {

	if NodeId < 1 {
		log.Fatal("No set TestNodeId")
	}

	req.SetDebug()

	ChSet, err := req.Changesets()

	if err != nil {
		log.Fatal(err)
	}

	if err := ChSet.OsmChange("delete"); err != nil {
		log.Fatal(err)
	}

	/* Load data for existing node */
	if _, err := ChSet.LoadNode(strconv.Itoa(NodeId)); err != nil {
		log.Fatal(err)
	}

	/* Upload new data */
	if _, err := ChSet.Upload(); err != nil {
		log.Fatal(err)
	}

	if err := ChSet.Close(); err != nil {
		log.Fatal(err)
	}
}

func main() {
	login := os.Getenv("OSM_USER")
	pass := os.Getenv("OSM_PASSWD")
	url := os.Getenv("OSM_URL")
	NodeId := flag.Int("nodeid", 0, "nodeid have to be integer.")
	flag.Parse()

	if login == "" || pass == "" {
		log.Fatalln("Login and password are not found.")
	}

	req := osmapi.MyRequest(login, pass)
	if url == "" {
		url = osmapi.MainURLTest
	}
	req.SetUrl(url)

	DeleteNode(req, *NodeId)
}
