package osmapi

import (
	//"github.com/davecgh/go-spew/spew"
	"log"
	"math/rand"
	"os"
	"strconv"
	"testing"
)

//export OSM_URL="http://api06.dev.openstreetmap.org"; export OSM_USER="lashko@corp.sputnik.ru"; export OSM_PASSWD="gdezhivetsputnik2";cd /kmsearch/go-osm-api/; set GOPATH="/kmsearch/go-osm-api/"; go test ./osmapi/ways.go ./osmapi/osmapi_ways_test.go ./osmapi/osmapi.go ./osmapi/changesets.go ./osmapi/capabilities.go

const minimumVersion = "0.6"
const maximumVersion = "0.6"

var (
	rnd       *rand.Rand = rand.New(rand.NewSource(99))
	TestWayId string     = ""
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

func Test_Ways(t *testing.T) {
	_01_WayLoad(t)
	_02_CreateWays(t)
	_03_DeleteWays(t)
}

func _01_WayLoad(t *testing.T) {

	req := init_req(t, "CreateSetUpload")
	if req == nil {
		return
	}

	req.SetDebug()

	ChSet, err := req.Changesets("modify")
	if err != nil {
		log.Println("_01_WayLoad. Create")
		t.Fatal(err)
	}

	/* Load existing way */
	if _, err_n := ChSet.WayLoad("52868"); err_n != nil {
		log.Println("_01_WayLoad. NewNode")
		t.Fatal(err_n)
	}

	_ChangeSetClose(t, ChSet)

	//t.Fatal("test view")
}

func _02_CreateWays(t *testing.T) {

	req := init_req(t, "CreateSetUpload")
	if req == nil {
		return
	}

	req.SetDebug()

	ChSet, err := req.Changesets("create")
	if err != nil {
		log.Println("_02_CreateWays. Create")
		t.Fatal(err)
	}

	/* Create new way */
	if _, err := ChSet.WayNew(); err != nil {
		log.Println("_02_CreateWays. NewNode")
		t.Fatal(err)
	}

	i := 6
	for i > 0 {
		i--

		TestLat := strconv.FormatFloat(-0.05+0.1*rnd.Float64(), 'f', 6, 64)
		TestLon := strconv.FormatFloat(51.50+0.01*rnd.Float64(), 'f', 6, 64)

		node, err := ChSet.NewNode(TestLat, TestLon)
		if err != nil {
			log.Println("_02_CreateWays. NewNode")
			t.Fatal(err)
		}

		var id string
		switch i {
		case 2:
			id, err = ChSet.WayAddNode(node, "0")
		case 1:
			id, err = ChSet.WayAddNode(node, "-3")
		default:
			id, err = ChSet.WayAddNode(node)
		}

		if err != nil {
			log.Println("_02_CreateWays. WayAddNode")
			t.Fatal(err)
		} else {
			log.Printf("Adds ref = %s\n", id)
		}
	}

	//
	if err = ChSet.WayDelNode("-3"); err != nil {
		log.Println("_02_CreateWays. WayDelNode")
		t.Fatal(err)
	}

	/* Now our sequence is: -2 -1 -5 -4 -6 */
	if TestWayId, err = ChSet.Upload(); err != nil {
		log.Println("_02_CreateWays. Upload")
		t.Fatal(err)
	}

	_ChangeSetClose(t, ChSet)

	//t.Fatal("test view")
}

func _03_DeleteWays(t *testing.T) {

	if "" == TestWayId {
		t.Skip("_03_DeleteWays")
		return
	}

	req := init_req(t, "DeleteWays")
	if req == nil {
		return
	}

	req.SetDebug()

	ChSet, err := req.Changesets("delete")
	if err != nil {
		log.Println("_03_DeleteWays. Create")
		t.Fatal(err)
	}

	/* Load existing way */
	if _, err_n := ChSet.WayLoad(TestWayId); err_n != nil {
		log.Println("v. NewNode")
		t.Fatal(err_n)
	}
	if err := ChSet.DelAllNodes(); err != nil {
		log.Println("_03_DeleteWays. Upload")
		t.Fatal(err)
	}
	/* Now our node's sequence is: -2 -1 -5 -4 -6 */
	if TestWayId, err = ChSet.Upload(); err != nil {
		log.Println("_03_DeleteWays. Upload")
		t.Fatal(err)
	}

	_ChangeSetClose(t, ChSet)

	//t.Fatal("test view")
}
