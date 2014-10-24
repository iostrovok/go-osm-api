package osmapi

import (
	"errors"
	//"github.com/davecgh/go-spew/spew"
	"gopkg.in/xmlpath.v2"
	"log"
	"os"
	"testing"
)

const minimumVersion = "0.6"
const maximumVersion = "0.6"

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

	if login == "" || pass == "" {
		log.Println("Skip " + title + ". login and password are not found.")
		return nil
	}

	req := MyRequest()
	req.UserPass(login, pass)

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

	c, err := req.Changesets()
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
	_, err_n := req.MiscellMap("41.06221", "57.24570", "41.07218", "57.24942")
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

func Test_ChangeSetUpload(t *testing.T) {

	req := init_req(t, "ChangeSetUpload")
	if req == nil {
		return
	}

	//req.SetDebug()

	ChSet, err := req.Changesets()
	if err != nil {
		log.Println("Test_ChangeSetUpload. Create")
		t.Fatal(err)
	}

	if ChSet.Id == "" {
		log.Println("Test_ChangeSetUpload")
		t.Fatal(err)
	}

	if err := ChSet.OsmChange("modify"); err != nil {
		log.Println("Test_ChangeSetUpload. OsmChange")
		t.Fatal(err)
	}

	/* Create new node */
	node, err_n := ChSet.LoadNode("1442930428")
	if err_n != nil {
		log.Println("Test_ChangeSetUpload. Create")
		t.Fatal(err_n)
	}

	/*  Set new data */
	node.AddTag("name:en", "Akultsevo")
	node.AddTag("name:uk", "Акульцево")
	node.AddTag("place", "hamlet")

	/* Create new node */
	node2, err_n2 := ChSet.LoadNode("1442930461")
	if err_n2 != nil {
		log.Println("Test_ChangeSetUpload. Create")
		t.Fatal(err_n2)
	}
	/*  Set new data */
	node2.AddTag("name:en", "Petrushikha")

	//req.SetDebug()

	/* Upload new data */
	if err := ChSet.Upload(); err != nil {
		log.Println("Test_ChangeSetUpload. Upload")
		t.Fatal(err)
	}

	_ChangeSetClose(t, ChSet)

	//t.Fatal("test view")
}
