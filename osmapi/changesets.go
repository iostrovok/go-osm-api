package osmapi

import (
	"encoding/xml"
	"errors"
	"log"
	"time"
)

/* ===
Changesets
*/

type ChangeSetSt struct {
	Id      string
	Request *MyRequestSt
	OsmCh   *OsmChangeSt
}

type TagSt struct {
	Key     string   `xml:"k,attr"`
	Val     string   `xml:"v,attr"`
	XMLName xml.Name `xml:"tag"`
}

type TagListSt struct {
	List []TagSt
}

func NewTag(k, v string) TagSt {
	t := TagSt{}
	t.Key = k
	t.Val = v
	return t
}

type WaySt struct {
	Tag       []TagSt `xml:"tag"`
	OsmId     string  `xml:"id,attr"`
	ReqId     string  `xml:"changeset,attr"`
	Visible   string  `xml:"visible,attr"`
	Lon       string  `xml:"lon,attr"`
	Lat       string  `xml:"lat,attr"`
	Version   string  `xml:"version,attr"`
	User      string  `xml:"user,attr"`
	Uid       string  `xml:"uid,attr"`
	Timestamp string  `xml:"timestamp,attr"`
}

type NodeSt struct {
	Tag       []TagSt `xml:"tag,omitempty"`
	OsmId     string  `xml:"id,attr,omitempty"`
	ReqId     string  `xml:"changeset,attr"`
	Visible   string  `xml:"visible,attr"`
	Lon       string  `xml:"lon,attr,omitempty"`
	Lat       string  `xml:"lat,attr,omitempty"`
	Version   string  `xml:"version,attr,omitempty"`
	User      string  `xml:"user,attr,omitempty"`
	Uid       string  `xml:"uid,attr,omitempty"`
	Timestamp string  `xml:"timestamp,attr,omitempty"`
}

type ChangeSt struct {
	//XMLName xml.Name `xml:",omitempty"`
	Node []*NodeSt `xml:"node"`
	Way  []*WaySt  `xml:"way"`
}

type OsmChangeSt struct {
	XMLName   xml.Name   `xml:"osmChange"`
	Version   string     `xml:"version,attr"`
	Generator string     `xml:"generator,attr"`
	Modify    *ChangeSt  `xml:"modify,omitempty"`
	Create    *ChangeSt  `xml:"create,omitempty"`
	Delete    *ChangeSt  `xml:"delete,omitempty"`
	Changeset *TagListSt `xml:"changeset,omitempty"`
	Type      string     `xml:"-"`
}

type OsmSt struct {
	XMLName   xml.Name   `xml:"osm"`
	Version   string     `xml:"version,attr"`
	Generator string     `xml:"generator,attr"`
	Changeset *TagListSt `xml:"changeset,omitempty"`
}

func (r *MyRequestSt) Changesets() (*ChangeSetSt, error) {
	c := ChangeSetSt{}

	c.Id = ""
	c.Request = r
	err := c.Create()
	return &c, err
}

/*   */
func (ChSet *ChangeSetSt) OsmChange(t string) error {
	OsmCh := OsmChangeSt{}

	if t != "create" && t != "modify" && t != "delete" && t != "changeset" {
		log.Fatalf("OsmChange. You have to use create|modify|delete as OsmChange type. Now it is `%s`", t)
	}

	OsmCh.Type = t
	OsmCh.Version = ProtocolVersion
	OsmCh.Generator = UserAgent
	nodes := []*NodeSt{}
	ways := []*WaySt{}
	ch := ChangeSt{nodes, ways}

	switch OsmCh.Type {
	case "modify":
		OsmCh.Modify = &ch
	case "create":
		OsmCh.Create = &ch
	case "delete":
		OsmCh.Delete = &ch
	}

	ChSet.OsmCh = &OsmCh

	return nil
}

func (c *NodeSt) AddTag(k, v string) {
	c.Tag = append(c.Tag, NewTag(k, v))
}

/*
Whan we creat new node we have get infomation from api.site
*/
func (ChSet *ChangeSetSt) LoadNode(OsmId string) (*NodeSt, error) {

	/* Answer has to be empty */
	data, err := ChSet.Request.GetXML("/api/0.6/node/" + OsmId)
	if err != nil {
		return nil, err
	}

	n := NodeSt{}
	n.Tag = []TagSt{}
	n.ReqId = ChSet.Id
	n.OsmId = OsmId
	n.Lon = xml_str(data, "/osm/node/@lon")
	n.Lat = xml_str(data, "/osm/node/@lat")
	n.Version = xml_str(data, "/osm/node/@version")
	n.Visible = xml_str(data, "/osm/node/@visible")

	tm := time.Now()
	n.Timestamp = tm.Format(TimeFormatLayout)

	ChSet.OsmCh._addNode(&n)

	return &n, nil
}

func (OsmCh *OsmChangeSt) _addNode(node *NodeSt, ways ...*WaySt) error {

	if len(ways) > 0 {
		//ways[0].
	}

	switch OsmCh.Type {
	case "modify":
		OsmCh.Modify.Node = append(OsmCh.Modify.Node, node)
	case "create":
		OsmCh.Create.Node = append(OsmCh.Modify.Node, node)
	case "delete":
		OsmCh.Delete.Node = append(OsmCh.Modify.Node, node)
	}

	return nil
}

/* ===
Changesets: Close: PUT /api/0.6/changeset/#id/close
*/
func (ChSet *ChangeSetSt) Close() error {
	/*  Changesets doesn't open. May by :) */
	if ChSet.Id == "" {
		return nil
	}

	/* Answer has to be empty */
	_, err := ChSet.Request.Put("/api/0.6/changeset/" + ChSet.Id + "/close")

	/* Clean memory. For any case */
	ChSet.OsmCh = nil
	ChSet = nil

	return err
}

/* ===
Changesets: Create: PUT /api/0.6/changeset/create
*/
func (ChSet *ChangeSetSt) Create() error {

	t := OsmSt{}
	t.Version = "0.6"
	t.Generator = UserAgent
	t.Changeset = &TagListSt{[]TagSt{NewTag("comment", "changeset comment"), NewTag("created_by", UserAgent)}}
	body2, err2 := xml.MarshalIndent(t, "", "")
	if err2 != nil {
		return err2
	}

	ChSet.Id = ""
	id, err := ChSet.Request.Put("/api/0.6/changeset/create", string(body2))
	if err == nil {
		ChSet.Id = id
	}

	return err
}

/* ===
Changesets: Diff upload: POST /api/0.6/changeset/#id/upload
*/
func (ChSet *ChangeSetSt) Upload() error {

	//(c *ChangeSetSt)
	if ChSet.Id == "" {
		errors.New("Cann't use uninitialize")
	}

	body, err_m := xml.MarshalIndent(ChSet.OsmCh, "", "")
	if err_m != nil {
		return err_m
	}

	data, err := ChSet.Request.PostXML("/api/0.6/changeset/"+ChSet.Id+"/upload", string(body))
	if err != nil {
		return err
	}

	old_id := xml_str(data, "/diffResult/node/@old_id")
	new_id := xml_str(data, "/diffResult/node/@new_id")

	if ChSet.OsmCh.Type == "modify" && old_id != new_id {
		return errors.New("Bad result")
	}

	if (ChSet.OsmCh.Type == "modify" || ChSet.OsmCh.Type == "create") && "" == new_id {
		return errors.New("Bad result")
	}

	return err
}
