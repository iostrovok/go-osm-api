package osmapi

import (
	"log"
	"testing"
)

var TestRelationId string = ""

func Test_Relations(t *testing.T) {
	_01_RelationLoad(t)
	_02_RelationCreate(t)
	_03_RelationModify(t)
	_04_RelationDelete(t)
}

func _01_RelationLoad(t *testing.T) {

	t.Log("\n\n--------------------- _01_RelationLoad -----------------------\n\n")

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
	if _, err_n := ChSet.RelationLoad("12993"); err_n != nil {
		log.Println("_01_RelationLoad. RelationLoad")
		t.Fatal(err_n)
	}
	_ChangeSetClose(t, ChSet)

	//t.Fatal("test view")
}

func _02_RelationCreate(t *testing.T) {

	t.Log("\n\n--------------------- _02_RelationCreate -----------------------\n\n")

	req := init_req(t, "CreateSetUpload")
	if req == nil {
		return
	}

	//req.SetDebug()

	ChSet, err := req.Changesets("create")
	if err != nil {
		log.Println("_02_RelationCreate. Create")
		t.Fatal(err)
	}

	// Create new relation
	if _, err := ChSet.RelationNew(); err != nil {
		log.Println("_02_RelationCreate. RelationNew")
		t.Fatal(err)
	}

	if err := ChSet.RelationAddMember("way", "52868", "outer"); err != nil {
		log.Println("_02_RelationCreate. AddMember 52868")
		t.Fatal(err)
	}

	if err := ChSet.RelationAddMember("node", "1282045", "inner"); err != nil {
		log.Println("_02_RelationCreate. AddMember 1282045")
		t.Fatal(err)
	}

	// Now our sequence is: -2 -1 -5 -4 -6
	if TestRelationId, err = ChSet.Upload(); err != nil {
		log.Println("_02_RelationCreate. Upload")
		t.Fatal(err)
	}

	_ChangeSetClose(t, ChSet)

	//t.Fatal("test view")
}

func _03_RelationModify(t *testing.T) {

	t.Log("\n\n--------------------- _03_RelationModify -----------------------\n\n")

	req := init_req(t, "CreateSetUpload")
	if req == nil {
		return
	}

	ChSet, err := req.Changesets("modify")
	if err != nil {
		log.Println("_03_RelationModify. Changesets")
		t.Fatal(err)
	}

	// Load existing relation
	if _, err_n := ChSet.RelationLoad(TestRelationId); err_n != nil {
		log.Println("_03_RelationModify. RelationLoad")
		t.Fatal(err_n)
	}

	if err := ChSet.RelationAddMember("way", "12820", "outer"); err != nil {
		log.Println("_03_RelationModify. AddMember 52868")
		t.Fatal(err)
	}

	if err := ChSet.RelationAddMember("node", "1282046", "inner"); err != nil {
		log.Println("_03_RelationModify. AddMember 1282045")
		t.Fatal(err)
	}

	if err := ChSet.RelationDelMember("node", "1282045"); err != nil {
		log.Println("_03_RelationModify. RelationDelMember 1282045")
		t.Fatal(err)
	}

	if err := ChSet.RelationDelMember("way", "52868"); err != nil {
		log.Println("_03_RelationModify. RelationDelMember 1282045")
		t.Fatal(err)
	}

	// Now our sequence is: -2 -1 -5 -4 -6
	if TestRelationId, err = ChSet.Upload(); err != nil {
		log.Println("_03_RelationModify. Upload")
		t.Fatal(err)
	}

	_ChangeSetClose(t, ChSet)

	//t.Fatal("test view")
}

func _04_RelationDelete(t *testing.T) {

	t.Log("\n\n--------------------- _04_RelationDelete -----------------------\n\n")

	if "" == TestRelationId {
		t.Skip("_04_RelationDelete")
		return
	}

	req := init_req(t, "DeleteRelations")
	if req == nil {
		return
	}

	ChSet, err := req.Changesets("delete")
	if err != nil {
		log.Println("_04_RelationDelete. Create")
		t.Fatal(err)
	}

	// Load existing relation
	if _, err_n := ChSet.RelationLoad(TestRelationId); err_n != nil {
		log.Println("_04_RelationDelete. RelationLoad")
		t.Fatal(err_n)
	}

	if err := ChSet.RelationDelAllMembers(); err != nil {
		log.Println("_04_RelationDelete. RelationDelAllMember")
		t.Fatal(err)
	}

	if TestRelationId, err = ChSet.Upload(); err != nil {
		log.Println("_04_RelationDelete. Upload")
		t.Fatal(err)
	}

	_ChangeSetClose(t, ChSet)

	//t.Fatal("test view")
}
