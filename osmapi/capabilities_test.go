package osmapi

import (
	"errors"
	"gopkg.in/xmlpath.v2"
	"log"
	"testing"
)

const minimumVersion = "0.6"
const maximumVersion = "0.6"

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
