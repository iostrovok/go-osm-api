GOLANG interface for OpenStreetMap API v0.6

	NOTICE! It is developer version.

Please see more inforamtion in OSM wiki: http://wiki.openstreetmap.org/wiki/API_v0.6

Installing 

	go get github.com/iostrovok/go-osm-api/osmapi

Using.
Part 1. Modify existing node (point)

	import "github.com/iostrovok/go-osm-api/osmapi"
	import "log"

	req := osmapi.MyRequest()
	req.UserPass(login, pass)
	if req == nil {
		log.Fatal("Request create")
	}

	ChSet, err := req.Changesets()
	if err != nil || ChSet.Id == "" {
		log.Fatal ("Changesets create")
	}

	/* Create ChangeSet */
	if err := ChSet.OsmChange("modify"); err != nil {
        log.Println("Test_ChangeSetUpload. OsmChange")
		t.Fatal(err)
	}

	/* 
		Create node for modification
		"1442930428" is real point
	*/
	if node, err_n := ChSet.LoadNode("1442930428"); err != nil {
		log.Fatal("Node сreate")
	} else {
		/*  Set new tag into data */
		node.AddTag("name:en", "Akultsevo")
		node.AddTag("name:uk", "Акульцево")
	}

	/* 
		Create node for modification
		"1442930461" is real point near "1442930428" :)
	*/
	if node, err := ChSet.LoadNode( "1442930461"); err != nil {
		log.Fatal("Node сreate")
	} else {
		/*  Set new tag into data */
		node.AddTag("name:en", "Petrushikha")
	}

	/* Upload new data to server */
	if err := ChSet.Upload(); err != nil {
		log.Fatal(err)
	}

	/*  Fixed result */
	if err := ChSet.Close(); err != nil {
		log.Fatal("Changesets close")
	}


