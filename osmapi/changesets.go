package osmapi

import (
	"encoding/xml"
	"errors"
	"time"
)

/* ===
Changesets
*/

type ChangeSetSt struct {
	Id      string
	Request *MyRequestSt
	OsmCh   *OsmChangeSt
	//Type   string     // node, way rel
}

type TagSt struct {
	Key     string   `xml:"k,attr"`
	Val     string   `xml:"v,attr"`
	XMLName xml.Name `xml:"tag"`
}

type TagListSt struct {
	List []*TagSt
}

func NewTag(k, v string) *TagSt {
	t := TagSt{}
	t.Key = k
	t.Val = v
	return &t
}

type NodeSt struct {
	Tag       []*TagSt `xml:"tag,omitempty"`
	OsmId     string   `xml:"id,attr,omitempty"`
	ReqId     string   `xml:"changeset,attr"`
	Visible   string   `xml:"visible,attr"`
	Lon       string   `xml:"lon,attr,omitempty"`
	Lat       string   `xml:"lat,attr,omitempty"`
	Version   string   `xml:"version,attr,omitempty"`
	User      string   `xml:"user,attr,omitempty"`
	Uid       string   `xml:"uid,attr,omitempty"`
	Timestamp string   `xml:"timestamp,attr,omitempty"`
}

type ChangeSt struct {
	//XMLName xml.Name `xml:",omitempty"`
	Node []*NodeSt `xml:"node,omitempty"`
	Way  *WaySt    `xml:"way,omitempty"`
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

func (r *MyRequestSt) Changesets(t string) (*ChangeSetSt, error) {
	c := ChangeSetSt{}

	c.Id = ""
	c.Request = r
	if err := c.Create(); err != nil {
		return nil, err
	}

	if err := c.OsmChange(t); err != nil {
		return nil, err
	}

	return &c, nil
}

/*   */
func (ChSet *ChangeSetSt) OsmChange(t string) error {
	OsmCh := OsmChangeSt{}

	if t != "create" && t != "modify" && t != "delete" && t != "changeset" {
		return errors.New("OsmChange. You have to use create|modify|delete as OsmChange type. Now it is " + t)
	}

	OsmCh.Type = t
	OsmCh.Version = ProtocolVersion
	OsmCh.Generator = UserAgent
	ch := ChangeSt{[]*NodeSt{}, nil}

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
	n := []*TagSt{NewTag(k, v)}
	for _, one := range c.Tag {
		if one.Key != k {
			n = append(n, one)
		}
	}
	c.Tag = n
}

/*
When we want to modify or delete node we have get infomation from api.site
*/
func (Request *MyRequestSt) LoadNodeDate(OsmId string) (*NodeSt, error) {

	/* Answer has to be empty */
	data, err := Request.GetXML("/api/0.6/node/" + OsmId)
	if err != nil {
		return nil, err
	}

	n := NodeSt{}
	n.Tag = []*TagSt{}
	n.Lat = xml_str(data, "/osm/node/@lat")
	n.Lon = xml_str(data, "/osm/node/@lon")
	n.OsmId = OsmId
	n.ReqId = xml_str(data, "/osm/node/@changeset")
	n.Timestamp = xml_str(data, "/osm/node/@timestamp")
	n.Version = xml_str(data, "/osm/node/@version")
	n.Visible = xml_str(data, "/osm/node/@visible")

	if n.Lon == "" || n.Lat == "" {
		return nil, errors.New("Note " + OsmId + " not found")
	}

	for _, v := range xml_slice(data, "/osm/node/tag", []string{"k", "v"}) {
		if v["k"] == "" || v["v"] == "" {
			continue
		}
		t := TagSt{}
		t.Key = v["k"]
		t.Val = v["v"]
		n.Tag = append(n.Tag, &t)
	}

	return &n, nil
}

func (ChSet *ChangeSetSt) LoadNode(OsmId string) (*NodeSt, error) {

	/* Answer has to be empty */
	n, err := ChSet.Request.LoadNodeDate(OsmId)
	if err != nil {
		return nil, err
	}

	n.ReqId = ChSet.Id
	tm := time.Now()
	n.Timestamp = tm.Format(TimeFormatLayout)

	ChSet.OsmCh._addNode(n)

	return n, nil
}

func (OsmCh *OsmChangeSt) _addNode(node *NodeSt) error {

	switch OsmCh.Type {
	case "modify":
		OsmCh.Modify.Node = append(OsmCh.Modify.Node, node)
	case "create":
		OsmCh.Create.Node = append(OsmCh.Create.Node, node)
	case "delete":
		OsmCh.Delete.Node = append(OsmCh.Delete.Node, node)
	}

	return nil
}

/*
When we creat new node
*/
func (ChSet *ChangeSetSt) NewNode(Lat, Lon string) (*NodeSt, error) {

	n := NodeSt{}
	n.Tag = []*TagSt{}
	n.ReqId = ChSet.Id
	n.OsmId = ""
	n.Lon = Lon
	n.Lat = Lat
	n.Version = "1"
	n.Visible = "true"

	tm := time.Now()
	n.Timestamp = tm.Format(TimeFormatLayout)

	ChSet.OsmCh._addNode(&n)

	return &n, nil
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
	t.Changeset = &TagListSt{[]*TagSt{NewTag("comment", "changeset comment"), NewTag("created_by", UserAgent)}}
	body2, err2 := xml.MarshalIndent(t, "", "")
	if err2 != nil {
		return err2
	}

	ChSet.Id = ""
	id, err := ChSet.Request.Put("/api/0.6/changeset/create", string(body2))
	if err == nil {
		ChSet.Id = id
	}

	if id == "" {
		return errors.New("Bad answer. Data from " + ChSet.Request.Url + " does not contain changeset's id.")
	}

	return err
}

/* ===
Changesets: Diff upload: POST /api/0.6/changeset/#id/upload
*/
func (ChSet *ChangeSetSt) Upload() (string, error) {

	//(c *ChangeSetSt)
	if ChSet.Id == "" {
		errors.New("Cann't use uninitialize")
	}

	body, err_m := xml.MarshalIndent(ChSet.OsmCh, "", "")
	if err_m != nil {
		return "", err_m
	}

	data, err := ChSet.Request.PostXML("/api/0.6/changeset/"+ChSet.Id+"/upload", string(body))
	if err != nil {
		return "", err
	}

	node_old_id := xml_str(data, "/diffResult/node/@old_id")
	node_new_id := xml_str(data, "/diffResult/node/@new_id")
	way_old_id := xml_str(data, "/diffResult/way/@old_id")
	way_new_id := xml_str(data, "/diffResult/way/@new_id")

	err_line := "Bad result ChangeSetSt upload."

	if ChSet.OsmCh.Type == "modify" {
		if node_old_id != "" && node_old_id != node_new_id {
			return "", errors.New(err_line + "Old node id equals new.")
		} else if way_old_id != "" && way_old_id != way_new_id {
			return "", errors.New(err_line + " Modify. Old way id equals new.")
		}
	}

	if ChSet.OsmCh.Type == "delete" && "0" != node_new_id && "" != way_new_id {
		return "", errors.New(err_line + " Delete. Bad new id " + ChSet.OsmCh.Type)
	}

	if ChSet.OsmCh.Type == "create" && "" == node_new_id && "" == way_new_id {
		return "", errors.New(err_line + " Create. New  node/way id empty.")
	}

	if way_new_id != "" {
		return way_new_id, nil
	}
	return node_new_id, nil
}

/*

Access functions to nodes

*/
func (OsmCh *OsmChangeSt) Nodes() []*NodeSt {

	switch OsmCh.Type {
	case "modify":
		return OsmCh.Modify.Node
	case "create":
		return OsmCh.Create.Node
	case "delete":
		return OsmCh.Delete.Node
	}

	return []*NodeSt{}
}

func (OsmCh *OsmChangeSt) Node(NodeId string) *NodeSt {
	list := OsmCh.Nodes()
	for _, v := range list {
		if v.OsmId == NodeId {
			return v
		}
	}
	return nil
}

func (ChSet *ChangeSetSt) DelAllNodes() error {
	return ChSet.OsmCh.DelAllNodes()
}

func (OsmCh *OsmChangeSt) DelAllNodes() error {

	if OsmCh.Modify != nil {
		OsmCh.Modify.Node = nil
	}
	if OsmCh.Create != nil {
		OsmCh.Create.Node = nil
	}
	if OsmCh.Delete != nil {
		OsmCh.Delete.Node = nil
	}

	return nil
}

func (ChSet *ChangeSetSt) DelNode(NodeId string) error {
	return ChSet.OsmCh.DelNode(NodeId)
}

func (OsmCh *OsmChangeSt) DelNode(NodeId string) error {
	list := OsmCh.Nodes()
	OsmCh.DelAllNodes()

	for _, v := range list {
		if v.OsmId == NodeId {
			continue
		}
		OsmCh._addNode(v)
	}

	return nil
}
