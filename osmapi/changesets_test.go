package osmapi

import (
	"testing"
)

func Test_ChangeSetCreate(t *testing.T) {

	req := init_req(t, "ChangeSetCreate")
	if req == nil {
		return
	}

	c, err := req.Changesets("delete")
	if err != nil {
		t.Log("Test_ChangeSetCreate. Changesets")
		t.Fatal(err)
	}

	if c.Id == "" {
		t.Log("Test_ChangeSetCreate. c.Id is empty.")
		t.Fatal(err)
	}

	_ChangeSetClose(t, c)

}
