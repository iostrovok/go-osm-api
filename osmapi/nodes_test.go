package osmapi

import (
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"testing"
)

var TestNodeId string = ""

func Test_Node(t *testing.T) {
	_01_NodeCreate(t)
	_02_NodeChange(t)
	_03_NodeDelete(t)
	_04_NodeChangeMassive(t)
}

func _01_NodeCreate(t *testing.T) {

	t.Log("\n\n--------------------- _01_NodeCreate -----------------------\n\n")

	req := init_req(t, "CreateSetUpload")
	if req == nil {
		return
	}

	req.SetDebug()

	ChSet, err := req.Changesets("create")
	if err != nil {
		log.Println("_01_NodeCreate. Create")
		t.Fatal(err)
	}

	rnd := rand.New(rand.NewSource(99))
	TestLat := strconv.FormatFloat(-0.05+0.1*rnd.Float64(), 'f', 6, 64)
	TestLon := strconv.FormatFloat(51.50+0.01*rnd.Float64(), 'f', 6, 64)

	/* Create new node */
	node, err_n := ChSet.NewNode(TestLat, TestLon)
	if err_n != nil {
		log.Println("_01_NodeCreate. NewNode")
		t.Fatal(err_n)
	}

	/*  Set new data */
	node.AddTag("name:en", "Anywhere in London")
	node.AddTag("name:ru", "Поселок где-то в лондоне")
	node.AddTag("name:uk", "Поселок где-то в лондоне")
	node.AddTag("place", "street")

	//req.SetDebug()

	/* Upload new data */
	if TestNodeId, err = ChSet.Upload(); err != nil {
		log.Println("_01_NodeCreate. Upload")
		t.Fatal(err)
	}

	if TestNodeId == "" {
		log.Println("_01_NodeCreate. TestNodeId is empty")
		t.Fatal(err)
	}

	t.Log("_01_NodeCreate. TestNodeId is " + TestNodeId + "\n")

	_ChangeSetClose(t, ChSet)

	//t.Fatal("test view")
}

func _04_NodeChangeMassive(t *testing.T) {

	t.Log("\n\n--------------------- _04_NodeChangeMassive -----------------------\n\n")
	ids := []string{}
	i := 1
	for i < 4 {
		i++
		_01_NodeCreate(t)
		if TestNodeId == "" {
			log.Fatal("_04_NodeChangeMassive. No set TestNodeId")
		}
		ids = append(ids, TestNodeId)
	}

	req := init_req(t, "ChangeSetUpload")
	if req == nil {
		return
	}

	req.SetDebug()

	ChSet, err := req.Changesets("modify")
	if err != nil {
		log.Println("_04_NodeChangeMassive. Create")
		t.Fatal(err)
	}

	for _, id := range ids {
		node, err_n := ChSet.LoadNode(id)
		if err_n != nil {
			log.Println("_04_NodeChangeMassive. LoadNode")
			t.Fatal(err_n)
		}

		/*  Set new data */
		node.AddTag("name:en", fmt.Sprintf("Anywhere in London - %s", id))
		node.AddTag("name:ru", fmt.Sprintf("Поселок где-то в лондоне - %s", id))
		node.AddTag("name:uk", fmt.Sprintf("Поселок где-то в лондоне - %s", id))
		node.AddTag("place", "hamlet")
	}

	req.SetDebug()

	/* Upload new data */
	if _, err := ChSet.Upload(); err != nil {
		log.Println("_04_NodeChangeMassive. Upload")
		t.Fatal(err)
	}

	_ChangeSetClose(t, ChSet)

	t.Fatal("test view")
}

func _02_NodeChange(t *testing.T) {

	t.Log("\n\n--------------------- _02_NodeChange -----------------------\n\n")

	if TestNodeId == "" {
		log.Fatal("_02_NodeChange. No set TestNodeId")
	}

	req := init_req(t, "ChangeSetUpload")
	if req == nil {
		return
	}

	req.SetDebug()

	ChSet, err := req.Changesets("modify")
	if err != nil {
		log.Println("_02_NodeChange. Create")
		t.Fatal(err)
	}

	/* Create new node */
	node, err_n := ChSet.LoadNode(TestNodeId)
	if err_n != nil {
		log.Println("_02_NodeChange. LoadNode")
		t.Fatal(err_n)
	}

	/*  Set new data */
	node.AddTag("name:en", "Hamlet in LONDON")
	node.AddTag("name:ru", "Поселок где-то в ЛОНДОНЕ")
	node.AddTag("name:uk", "Поселок где-то в ЛОНДОНЕ")
	node.AddTag("place", "hamlet")

	//req.SetDebug()

	/* Upload new data */
	if _, err := ChSet.Upload(); err != nil {
		log.Println("_02_NodeChange. Upload")
		t.Fatal(err)
	}

	_ChangeSetClose(t, ChSet)

	//t.Fatal("test view")
}

func _03_NodeDelete(t *testing.T) {

	t.Log("\n\n--------------------- _03_NodeDelete -----------------------\n\n")

	if TestNodeId == "" {
		log.Fatal("_03_NodeDelete. No set TestNodeId")
	}

	req := init_req(t, "DeleteSetUpload")
	if req == nil {
		return
	}

	req.SetDebug()

	ChSet, err := req.Changesets("delete")
	if err != nil {
		log.Println("_03_NodeDelete. Create")
		t.Fatal(err)
	}

	/* Create new node */
	_, err_n := ChSet.LoadNode(TestNodeId)
	if err_n != nil {
		log.Println("_03_NodeDelete. LoadNode")
		t.Fatal(err_n)
	}

	/* Upload new data */
	if _, err := ChSet.Upload(); err != nil {
		log.Println("_03_NodeDelete. Upload")
		t.Fatal(err)
	}

	TestNodeId = ""

	_ChangeSetClose(t, ChSet)

	//t.Fatal("test view")
}
