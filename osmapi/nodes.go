package osmapi

import (
	"encoding/xml"
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
	Tags      []*TagSt `xml:"tag,omitempty"`
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

func (c *NodeSt) Tag(k string) (string, bool) {
	for _, one := range c.Tags {
		if one.Key == k {
			return one.Val, true
		}
	}
	return "", false
}

func (c *NodeSt) AddTag(k, v string) {
	n := []*TagSt{NewTag(k, v)}
	for _, one := range c.Tags {
		if one.Key != k {
			n = append(n, one)
		}
	}
	c.Tags = n
}

func (c *NodeSt) DelTag(k string) {
	n := []*TagSt{}
	for _, one := range c.Tags {
		if one.Key != k {
			n = append(n, one)
		}
	}
	c.Tags = n
}

/*
When we want to modify or delete node we have get infomation from api.site
*/
func (ChSet *ChangeSetSt) LoadNode(OsmId string) (*NodeSt, error) {

	/* Answer has to be empty */
	n, err := ChSet.Request.LoadNodeDate(OsmId)
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
	n.Tags = []*TagSt{}
	n.ReqId = ChSet.Id
	n.OsmId = ""
	n.Lon = Lon
	n.Lat = Lat
	n.Version = "1"
	n.Visible = "true"

	tm := time.Now()
	n.Timestamp = tm.Format(TimeFormatLayout)

	ChSet._addNode(&n)

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
