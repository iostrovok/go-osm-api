package osmapi

import (
	"log"
	"math/rand"
	"strconv"
	"testing"
)

var TestWayId string = ""

func Test_Ways(t *testing.T) {
	_01_WayLoad(t)
	_02_CreateWays(t)
	_03_DeleteWays(t)
}

func _01_WayLoad(t *testing.T) {

	t.Log("\n\n--------------------- _01_WayLoad -----------------------\n\n")

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

	t.Log("\n\n--------------------- _02_CreateWays -----------------------\n\n")

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

	rnd := rand.New(rand.NewSource(99))

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
			t.Log("_02_CreateWays. WayAddNode\n")
			t.Fatal(err)
		} else {
			log.Printf("Adds ref = %s\n", id)
		}
	}

	/* Now our sequence is: -2 -1 -5 -4 -6 */
	if err = ChSet.WayDelNode("-3"); err != nil {
		log.Println("_02_CreateWays. WayDelNode")
		t.Fatal(err)
	}

	if TestWayId, err = ChSet.Upload(); err != nil {
		log.Println("_02_CreateWays. Upload")
		t.Fatal(err)
	}

	if TestWayId == "" {
		log.Println("_02_CreateWays. TestWayId is empty")
		t.Fatal(err)
	}

	t.Log("_02_CreateWays. TestWayId is " + TestWayId + "\n")

	_ChangeSetClose(t, ChSet)

	//t.Fatal("test view")
}

func _03_DeleteWays(t *testing.T) {

	t.Log("\n\n--------------------- _03_DeleteWays -----------------------\n\n")

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
		log.Println("_03_DeleteWays. WayLoad")
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
