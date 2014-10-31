package osmapi

import (
	"encoding/xml"
	"errors"
	"time"
)

/* ===
Nodes
*/

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
func (ChSet *ChangeSetSt) LoadNodeDate(OsmId string) (*NodeSt, error) {

	/* Answer has to be empty */
	data, err := ChSet.Request.GetXML("/api/0.6/node/" + OsmId)
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
	n, err := ChSet.LoadNodeDate(OsmId)
	if err != nil {
		return nil, err
	}

	n.ReqId = ChSet.Id
	tm := time.Now()
	n.Timestamp = tm.Format(TimeFormatLayout)

	ChSet._addNode(n)

	return n, nil
}

func (ChSet *ChangeSetSt) _addNode(node *NodeSt) error {

	switch ChSet.OsmCh.Type {
	case "modify":
		ChSet.OsmCh.Modify.Node = append(ChSet.OsmCh.Modify.Node, node)
	case "create":
		ChSet.OsmCh.Create.Node = append(ChSet.OsmCh.Create.Node, node)
	case "delete":
		ChSet.OsmCh.Delete.Node = append(ChSet.OsmCh.Delete.Node, node)
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

	ChSet._addNode(&n)
	ChSet.OsmCh.ChangeType = "node"

	return &n, nil
}

/*

Access functions to nodes

*/
func (ChSet *ChangeSetSt) Nodes() []*NodeSt {

	switch ChSet.OsmCh.Type {
	case "modify":
		return ChSet.OsmCh.Modify.Node
	case "create":
		return ChSet.OsmCh.Create.Node
	case "delete":
		return ChSet.OsmCh.Delete.Node
	}

	return []*NodeSt{}
}

func (ChSet *ChangeSetSt) Node(NodeId string) *NodeSt {
	list := ChSet.Nodes()
	for _, v := range list {
		if v.OsmId == NodeId {
			return v
		}
	}
	return nil
}

func (ChSet *ChangeSetSt) DelAllNodes() error {

	if ChSet.OsmCh.Modify != nil {
		ChSet.OsmCh.Modify.Node = nil
	}
	if ChSet.OsmCh.Create != nil {
		ChSet.OsmCh.Create.Node = nil
	}
	if ChSet.OsmCh.Delete != nil {
		ChSet.OsmCh.Delete.Node = nil
	}

	return nil
}

func (ChSet *ChangeSetSt) DelNode(NodeId string) error {

	list := ChSet.Nodes()

	ChSet.DelAllNodes()

	for _, v := range list {
		if v.OsmId == NodeId {
			continue
		}
		ChSet._addNode(v)
	}

	return nil
}
