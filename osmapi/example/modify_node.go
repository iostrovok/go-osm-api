package main

import (
	"flag"
	"github.com/iostrovok/go-osm-api/osmapi"
	"log"
	"os"
	"strconv"
)

func ModifyNode(req *osmapi.MyRequestSt, NodeId int, k, v string) {

	if NodeId < 1 {
		log.Fatal("No set TestNodeId")
	}

	req.SetDebug()

	/* Make object for request */
	ChSet, err := req.Changesets()

	if err != nil {
		log.Fatal(err)
	}

	if err := ChSet.OsmChange("modify"); err != nil {
		log.Fatal(err)
	}

	/* Load data for existing node */
	node, err_n := ChSet.LoadNode(strconv.Itoa(NodeId))
	if err_n != nil {
		log.Fatal(err_n)
	}

	/*  Set new data */
	node.AddTag(k, v)

	/* Upload new data */
	if _, err := ChSet.Upload(); err != nil {
		log.Fatal(err)
	}

	/* Fixed result into OSM */
	if err := ChSet.Close(); err != nil {
		log.Fatal(err)
	}
	/* ChSet is nil now */
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

	ModifyNode(req, *NodeId, "name:ru", "Поселок где-то в ЛОНДОНЕ")
}
