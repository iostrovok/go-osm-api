package osmapi

import (
	"log"
	"os"
	"testing"
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

	req := MyRequest()
	req.UserPass(login, pass)
	if url == "" {
		url = MainURLTest
	}
	req.SetUrl(url)

	req.SetDebug()

	return req
}

func _ChangeSetClose(t *testing.T, c *ChangeSetSt) {
	err := c.Close()
	if err != nil {
		log.Println("_ChangeSetClose. Close.")
		t.Fatal(err)
	}
}
