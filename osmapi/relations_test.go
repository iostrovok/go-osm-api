package osmapi

import (
	"github.com/davecgh/go-spew/spew"
	"log"
	"math/rand"
	"os"
	//"strconv"
	"testing"
)

//export OSM_URL="http://api06.dev.openstreetmap.org"; export OSM_USER="lashko@corp.sputnik.ru"; export OSM_PASSWD="gdezhivetsputnik2";cd /kmsearch/go-osm-api/; set GOPATH="/kmsearch/go-osm-api/"; go test ./osmapi/relations.go ./osmapi/osmapi_relations_test.go ./osmapi/osmapi.go ./osmapi/changesets.go ./osmapi/capabilities.go

const minimumVersion = "0.6"
const maximumVersion = "0.6"

var (
	rnd            *rand.Rand = rand.New(rand.NewSource(99))
	TestRelationId string     = ""
)

func Test_MyRequestSt(t *testing.T) {
	req := MyRequest()
	body, err := req.Get("")

	if err != nil || body == "" {
		log.Println("Test_MyRequestSt")
		t.Fatal(err)
	}
}

func init_req(t *testing.T, title string) *MyRequestSt {
	login := os.Getenv("OSM_USER")
	pass := os.Getenv("OSM_PASSWD")
	url := os.Getenv("OSM_URL")

	if login == "" || pass == "" {
		log.Println("Skip " + title + ". login and password are not found.")
		return nil
	}

	/*
		if "CreateSetUpload" != title && "DeleteSetUpload" != title && "ChangeSetUpload" != title {
			return nil
		}
	*/

	req := MyRequest()
	req.UserPass(login, pass)
	if url == "" {
		url = MainURLTest
	}
	req.SetUrl(url)

	return req
}

func _ChangeSetClose(t *testing.T, c *ChangeSetSt) {
	err := c.Close()
	if err != nil {
		log.Println("Test_ChangeSetCreate. Close.")
		t.Fatal(err)
	}
}

func Test_Relations(t *testing.T) {
	//_01_RelationLoad(t)
	_02_CreateRelations(t)
	_03_ModifyRelations(t)
	_04_DeleteRelations(t)
}

func _01_RelationLoad(t *testing.T) {

	req := init_req(t, "CreateSetUpload")
	if req == nil {
		return
	}

	//req.SetDebug()

	ChSet, err := req.Changesets("modify")
	if err != nil {
		log.Println("_01_RelationLoad. Changesets")
		t.Fatal(err)
	}

	/* Load existing relation  2996187 */
	if r, err_n := ChSet.RelationLoad("12993"); err_n != nil {
		log.Println("_01_RelationLoad. RelationLoad")
		t.Fatal(err_n)
	} else {
		spew.Dump(r)
	}

	_ChangeSetClose(t, ChSet)

	//t.Fatal("test view")
}

func _02_CreateRelations(t *testing.T) {

	req := init_req(t, "CreateSetUpload")
	if req == nil {
		return
	}

	//req.SetDebug()

	ChSet, err := req.Changesets("create")
	if err != nil {
		log.Println("_02_CreateRelations. Create")
		t.Fatal(err)
	}

	// Create new relation
	if _, err := ChSet.RelationNew(); err != nil {
		log.Println("_02_CreateRelations. RelationNew")
		t.Fatal(err)
	}

	if err := ChSet.RelationAddMember("way", "52868", "outer"); err != nil {
		log.Println("_02_CreateRelations. AddMember 52868")
		t.Fatal(err)
	}

	if err := ChSet.RelationAddMember("node", "1282045", "inner"); err != nil {
		log.Println("_02_CreateRelations. AddMember 1282045")
		t.Fatal(err)
	}

	// Now our sequence is: -2 -1 -5 -4 -6
	if TestRelationId, err = ChSet.Upload(); err != nil {
		log.Println("_02_CreateRelations. Upload")
		t.Fatal(err)
	}

	_ChangeSetClose(t, ChSet)

	//t.Fatal("test view")
}

func _03_ModifyRelations(t *testing.T) {

	req := init_req(t, "CreateSetUpload")
	if req == nil {
		return
	}

	req.SetDebug()

	ChSet, err := req.Changesets("modify")
	if err != nil {
		log.Println("_03_ModifyRelations. Changesets")
		t.Fatal(err)
	}

	// Load existing relation
	if _, err_n := ChSet.RelationLoad(TestRelationId); err_n != nil {
		log.Println("_04_DeleteRelations. RelationLoad")
		t.Fatal(err_n)
	}

	if err := ChSet.RelationAddMember("way", "12820", "outer"); err != nil {
		log.Println("_03_ModifyRelations. AddMember 52868")
		t.Fatal(err)
	}

	if err := ChSet.RelationAddMember("node", "1282046", "inner"); err != nil {
		log.Println("_03_ModifyRelations. AddMember 1282045")
		t.Fatal(err)
	}

	if err := ChSet.RelationDelMember("node", "1282045"); err != nil {
		log.Println("_03_ModifyRelations. RelationDelMember 1282045")
		t.Fatal(err)
	}

	if err := ChSet.RelationDelMember("way", "52868"); err != nil {
		log.Println("_03_ModifyRelations. RelationDelMember 1282045")
		t.Fatal(err)
	}

	// Now our sequence is: -2 -1 -5 -4 -6
	if TestRelationId, err = ChSet.Upload(); err != nil {
		log.Println("_03_ModifyRelations. Upload")
		t.Fatal(err)
	}

	//req.SetDebug(false)

	_ChangeSetClose(t, ChSet)

	//t.Fatal("test view")
}

func _04_DeleteRelations(t *testing.T) {

	if "" == TestRelationId {
		t.Skip("_04_DeleteRelations")
		return
	}

	req := init_req(t, "DeleteRelations")
	if req == nil {
		return
	}

	req.SetDebug()

	ChSet, err := req.Changesets("delete")
	if err != nil {
		log.Println("_04_DeleteRelations. Create")
		t.Fatal(err)
	}

	// Load existing relation
	if _, err_n := ChSet.RelationLoad(TestRelationId); err_n != nil {
		log.Println("_04_DeleteRelations. RelationLoad")
		t.Fatal(err_n)
	}

	if err := ChSet.RelationDelAllMembers(); err != nil {
		log.Println("_04_DeleteRelations. RelationDelAllMember")
		t.Fatal(err)
	}

	if TestRelationId, err = ChSet.Upload(); err != nil {
		log.Println("_04_DeleteRelations. Upload")
		t.Fatal(err)
	}

	_ChangeSetClose(t, ChSet)

	//t.Fatal("test view")
}
