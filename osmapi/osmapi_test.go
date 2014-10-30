package osmapi

import (
	"errors"
	//"github.com/davecgh/go-spew/spew"
	"gopkg.in/xmlpath.v2"
	"log"
	"math/rand"
	"os"
	"strconv"
	"testing"
)

const minimumVersion = "0.6"
const maximumVersion = "0.6"

var (
	rnd        *rand.Rand = rand.New(rand.NewSource(99))
	TestNodeId string     = ""
	TestLat    string     = strconv.FormatFloat(-0.05+0.1*rnd.Float64(), 'f', 6, 64)
	TestLon    string     = strconv.FormatFloat(51.50+0.01*rnd.Float64(), 'f', 6, 64)
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

func Test_ChangeSetCreate(t *testing.T) {

	req := init_req(t, "ChangeSetCreate")
	if req == nil {
		return
	}

	c, err := req.Changesets("delete")
	if err != nil {
		log.Println("Test_ChangeSetCreate. Create")
		t.Fatal(err)
	}

	if c.Id == "" {
		log.Println("Test_ChangeSetCreate")
		t.Fatal(err)
	}

	_ChangeSetClose(t, c)

}

func Test_MiscellMap(t *testing.T) {
	req := init_req(t, "MiscellMap")
	if req == nil {
		return
	}

	req.SetDebug()

	_, err_n := req.MiscellMap("-0.1275", "51.497222", "0.1275", "51.517222")
	if err_n != nil {
		log.Println("Test_MiscellMap")
		t.Fatal(err_n)
	}
}

func Test_MiscellPermissions(t *testing.T) {
	req := init_req(t, "MiscellPermissions")
	if req == nil {
		return
	}

	node, err_n := req.MiscellPermissions()
	if err_n != nil {
		log.Println("Test_MiscellPermissions")
		t.Fatal(err_n)
	}

	permission_name := xml_str(node, "/osm/permissions/permission/@name")
	if permission_name == "" {
		t.Fatal(errors.New("Test_MiscellPermissions. Not found '/osm/permissions/permission/@name'"))
	}
}

func Test_MiscellCap(t *testing.T) {

	req := init_req(t, "MiscellCap")
	if req == nil {
		return
	}

	node, err_n := req.MiscellCap()
	if err_n != nil {
		log.Println("Test_MiscellCap")
		t.Fatal(err_n)
	}

	path := xmlpath.MustCompile("/osm/api/version/@minimum")
	minimum, ok := path.String(node)
	if !ok {
		t.Fatal(errors.New("Test_MiscellCap. Not found 'osm/api/version/@minimum'"))
	}

	path = xmlpath.MustCompile("/osm/api/version/@maximum")
	maximum, ok := path.String(node)
	if !ok {
		t.Fatal(errors.New("Test_MiscellCap. Not found 'osm/api/version/@maximum'"))
	}

	if minimum != minimumVersion {
		t.Fatal(errors.New("Test_MiscellCap. Bad minimum version"))
	}
	if maximum != maximumVersion {
		t.Fatal(errors.New("Test_MiscellCap. Bad maximum version"))
	}
}

func Test_Node(t *testing.T) {
	_01_CreateSetUpload(t)
	_02_ChangeSetUpload(t)
	_03_DeleteSetUpload(t)
}

func _01_CreateSetUpload(t *testing.T) {

	req := init_req(t, "CreateSetUpload")
	if req == nil {
		return
	}

	req.SetDebug()

	ChSet, err := req.Changesets("create")
	if err != nil {
		log.Println("_01_CreateSetUpload. Create")
		t.Fatal(err)
	}

	/* Create new node */
	node, err_n := ChSet.NewNode(TestLat, TestLon)
	if err_n != nil {
		log.Println("_01_CreateSetUpload. NewNode")
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
		log.Println("_01_CreateSetUpload. Upload")
		t.Fatal(err)
	}

	_ChangeSetClose(t, ChSet)

	//t.Fatal("test view")
}

func _02_ChangeSetUpload(t *testing.T) {

	if TestNodeId == "" {
		log.Fatal("_02_ChangeSetUpload. No set TestNodeId")
	}

	req := init_req(t, "ChangeSetUpload")
	if req == nil {
		return
	}

	req.SetDebug()

	ChSet, err := req.Changesets("modify")
	if err != nil {
		log.Println("_02_ChangeSetUpload. Create")
		t.Fatal(err)
	}

	/* Create new node */
	node, err_n := ChSet.LoadNode(TestNodeId)
	if err_n != nil {
		log.Println("_02_ChangeSetUpload. LoadNode")
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
		log.Println("_02_ChangeSetUpload. Upload")
		t.Fatal(err)
	}

	_ChangeSetClose(t, ChSet)

	//t.Fatal("test view")
}

func _03_DeleteSetUpload(t *testing.T) {

	if TestNodeId == "" {
		log.Fatal("_03_DeleteSetUpload. No set TestNodeId")
	}

	req := init_req(t, "DeleteSetUpload")
	if req == nil {
		return
	}

	req.SetDebug()

	ChSet, err := req.Changesets("delete")
	if err != nil {
		log.Println("_03_DeleteSetUpload. Create")
		t.Fatal(err)
	}

	/* Create new node */
	_, err_n := ChSet.LoadNode(TestNodeId)
	if err_n != nil {
		log.Println("_03_DeleteSetUpload. LoadNode")
		t.Fatal(err_n)
	}

	/* Upload new data */
	if _, err := ChSet.Upload(); err != nil {
		log.Println("_03_DeleteSetUpload. Upload")
		t.Fatal(err)
	}

	TestNodeId = ""

	_ChangeSetClose(t, ChSet)

	//t.Fatal("test view")
}
