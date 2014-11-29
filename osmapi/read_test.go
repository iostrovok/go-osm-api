package osmapi

import (
	"log"
	//"math/rand"
	//"strconv"
	"testing"
)

var TestReadNodeId string = "1996249"
var TestReadWayId string = "6667"

func Test_Read(t *testing.T) {
	_01_ReadNode(t)
	//_02_ReadWay(t)
}

func _01_ReadNode(t *testing.T) {

	t.Log("\n\n--------------------- _01_ReadNode -----------------------\n\n")

	req := init_req(t, "CreateSetUpload")
	if req == nil {
		return
	}

	req.SetDebug()

	/* Create new node */
	node, err_n := req.LoadNodeDate(TestReadNodeId)
	if err_n != nil {
		log.Println("_01_ReadNode. LoadNodeDate")
		t.Fatal(err_n)
	}

	/* Upload new data */
	if source, find := node.Tag("source"); !find {
		log.Println()
		t.Fatal("_01_ReadNode. LoadNodeDate -> Tag for " + TestReadNodeId)
	} else if source != "ourfootprints" {
		t.Fatal("_01_ReadNode. LoadNodeDate -> Tag for " + TestReadNodeId)
	}

	//	t.Fatal("test view")
}
