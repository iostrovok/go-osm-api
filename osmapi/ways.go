package osmapi

import (
	"encoding/xml"
	"errors"
	"gopkg.in/xmlpath.v2"
	"strconv"
	"time"
)

/*

Way create

<osm>
 <way changeset="12">
   <tag k="note" v="Just a way"/>
   ...
   <nd ref="123"/>
   <nd ref="4345"/>
   ...
 </way>
</osm>

*/

type WayCreate struct {
	ChangesetId string `xml:"changeset,attr"`
}

type OsmWayCreate struct {
	Way     *WayCreate `xml:"way"`
	XMLName xml.Name   `xml:"osm"`
}

/* ===
Ways in changesets
*/

type WayNdSt struct {
	Ref string `xml:"ref,attr,omitempty"`
}

type WaySt struct {
	Tags      []*TagSt   `xml:"tag"`
	Nodes     []*WayNdSt `xml:"nd"`
	OsmId     string     `xml:"id,attr"`
	ReqId     string     `xml:"changeset,attr"`
	Timestamp string     `xml:"timestamp,attr"`
	Uid       string     `xml:"uid,attr,omitempty"`
	User      string     `xml:"user,attr,omitempty"`
	Version   string     `xml:"version,attr,omitempty"`
	Visible   string     `xml:"visible,attr,omitempty"`
	IsNew     bool       `xml:"-"`
}

/*
When we want to modify or delete node we have get infomation from api.site
*/
func (ChSet *ChangeSetSt) WayNew() (*WaySt, error) {
	w := WaySt{}
	w.Tags = nil
	w.Nodes = nil
	w.OsmId = "-1"
	w.ReqId = ChSet.Id
	w.Version = "1"
	w.Visible = "true"
	w.IsNew = true

	tm := time.Now()
	w.Timestamp = tm.Format(TimeFormatLayout)

	ChSet.OsmCh._setWay(&w)

	return &w, nil
}

func (ChSet *ChangeSetSt) WayLoadData(OsmId string) (*xmlpath.Node, error) {

	/* Answer has to be empty */
	data, err := ChSet.Request.GetXML("/api/0.6/way/" + OsmId)
	if err != nil {
		return nil, err
	}

	if "" == xml_str(data, "/osm/way/@id") {
		return nil, errors.New("WayLoadData. Way [" + OsmId + "]no found.")
	}

	return data, nil
}

/*
When we want to modify or delete node we have get infomation from api.site
*/
func (ChSet *ChangeSetSt) WayLoad(OsmId string) (*WaySt, error) {

	/* Answer has to be empty */
	data, err := ChSet.WayLoadData(OsmId)
	if err != nil {
		return nil, err
	}

	w, err_w := ChSet.WayNew()
	if err_w != nil {
		return nil, err_w
	}

	w.OsmId = OsmId
	w.Version = xml_str(data, "/osm/way/@version")
	w.Visible = xml_str(data, "/osm/way/@visible")
	w.Uid = xml_str(data, "/osm/way/@uid")
	w.User = xml_str(data, "/osm/way/@user")
	w.IsNew = false

	ChSet.OsmCh._setWay(w)

	for _, v := range xml_slice(data, "/osm/way/nd", []string{"ref"}) {
		if v["ref"] == "" {
			continue
		}
		if err := ChSet.LoadRef(v["ref"]); err != nil {
			return nil, err
		}
	}

	for _, v := range xml_slice(data, "/osm/way/tag", []string{"k", "v"}) {
		w._add_tag(v["k"], v["v"])
	}

	return w, nil
}

func (w *WaySt) _add_tag(k, v string) {
	if k == "" || v == "" {
		return
	}
	if w.Tags == nil {
		w.Tags = []*TagSt{}
	}

	t := TagSt{}
	t.Key = k
	t.Val = v

	w.Tags = append(w.Tags, &t)
}

func (ChSet *ChangeSetSt) LoadRef(ref string) error {
	/* Answer has to be empty */

	if !ChSet._checkWay() {
		return errors.New("LoadRef. No way changeset")
	}

	n, err := ChSet.LoadNodeDate(ref)
	if err != nil {
		return err
	}

	tm := time.Now()
	n.Timestamp = tm.Format(TimeFormatLayout)

	ChSet.OsmCh._addNode(n)
	if err := ChSet._put_ref_to_way(n.OsmId); err != nil {
		return err
	}

	return nil
}

func (ChSet *ChangeSetSt) WayAddNode(n *NodeSt, after_node_id ...string) (string, error) {
	/* Answer has to be empty */
	if !ChSet._checkWay() {
		//return "", errors.New("WayAddNode. No way changeset")
	}

	if n.OsmId == "" {
		id, err := ChSet._next_ref_id()
		if err != nil {
			return "", err
		}
		n.OsmId = id
	}

	tm := time.Now()
	n.Timestamp = tm.Format(TimeFormatLayout)

	// TODO: Node already added into ChangeSet
	// ChSet.OsmCh._addNode(n)
	if err := ChSet._put_ref_to_way(n.OsmId, after_node_id...); err != nil {
		return "", err
	}

	ChSet._update_way_id()

	return n.OsmId, nil
}

func (ChSet *ChangeSetSt) WayDelAllNodes() error {
	/* Answer has to be empty */
	if !ChSet._checkWay() {
		return errors.New("WayDelNode. No way changeset")
	}

	ChSet.DelAllNodes()
	if err := ChSet._del_all_ref_from_way(); err != nil {
		return err
	}

	return nil
}

func (ChSet *ChangeSetSt) WayDelNode(OsmId string) error {
	/* Answer has to be empty */
	if !ChSet._checkWay() {
		return errors.New("WayDelNode. No way changeset")
	}

	ChSet.OsmCh.DelNode(OsmId)
	if err := ChSet._del_ref_from_way(OsmId); err != nil {
		return err
	}

	return nil
}

func (ChSet *ChangeSetSt) _checkWay() bool {
	return ChSet.OsmCh._checkWay()
}

func (OsmCh *OsmChangeSt) _checkWay() bool {
	switch OsmCh.Type {
	case "modify":
		if OsmCh.Modify.Way != nil {
			return true
		}
	case "create":
		if OsmCh.Create.Way != nil {
			return true
		}
	case "delete":
		if OsmCh.Delete.Way != nil {
			return true
		}
	}
	return false
}

func (OsmCh *OsmChangeSt) _setWay(way *WaySt) error {
	switch OsmCh.Type {
	case "modify":
		OsmCh.Modify.Way = way
	case "create":
		OsmCh.Create.Way = way
	case "delete":
		OsmCh.Delete.Way = way
	}
	return errors.New("_setWay. No way changeset")

}

func (ChSet *ChangeSetSt) _next_ref_id() (string, error) {
	switch ChSet.OsmCh.Type {
	case "modify":
		return ChSet.OsmCh.Modify.Way._next_ref_id()
	case "create":
		return ChSet.OsmCh.Create.Way._next_ref_id()
	case "delete":
		return ChSet.OsmCh.Delete.Way._next_ref_id()
	}
	return "", errors.New("_next_ref_id. No way changeset")
}

func (w *WaySt) _next_ref_id() (string, error) {
	if w.Nodes == nil {
		return "-1", nil
	}

	i := -1
	for _, v := range w.Nodes {
		id, err := strconv.Atoi(v.Ref)
		if err != nil {
			return "", err
		}
		if i >= id {
			i = id - 1
		}
	}

	return strconv.Itoa(i), nil
}

func (ChSet *ChangeSetSt) _del_all_ref_from_way() error {
	switch ChSet.OsmCh.Type {
	case "modify":
		ChSet.OsmCh.Modify.Way = nil
		return nil
	case "create":
		ChSet.OsmCh.Create.Way = nil
		return nil
	case "delete":
		ChSet.OsmCh.Delete.Way = nil
		return nil
	}

	return errors.New("_del_all_ref_from_way. No way changeset")
}

func (ChSet *ChangeSetSt) _del_ref_from_way(ref string) error {
	switch ChSet.OsmCh.Type {
	case "modify":
		return ChSet.OsmCh.Modify.Way._del_ref_from_way(ref)
	case "create":
		return ChSet.OsmCh.Create.Way._del_ref_from_way(ref)
	case "delete":
		return ChSet.OsmCh.Delete.Way._del_ref_from_way(ref)
	}

	return errors.New("_del_ref_from_way. No way changeset")
}

func (w *WaySt) _check_nodes() {
	if w.Nodes == nil {
		w.Nodes = []*WayNdSt{}
	}
}

func (w *WaySt) _del_ref_from_way(NodeId string) error {

	w._check_nodes()

	nds := []*WayNdSt{}
	for _, v := range w.Nodes {

		if v.Ref != NodeId {
			nds = append(nds, v)
		}
	}

	w.Nodes = nds
	return nil
}

func (ChSet *ChangeSetSt) _update_way_id() error {

	id, err := ChSet._next_ref_id()
	if err != nil {
		return err
	}

	switch ChSet.OsmCh.Type {
	case "modify":
		return ChSet.OsmCh.Modify.Way._update_way_id(id)
	case "create":
		return ChSet.OsmCh.Create.Way._update_way_id(id)
	case "delete":
		return ChSet.OsmCh.Delete.Way._update_way_id(id)
	}

	return errors.New("_update_way_id. No way changeset")
}

func (w *WaySt) _update_way_id(NextId string) error {
	if w.IsNew {
		w.OsmId = NextId
	}

	return nil
}

func (ChSet *ChangeSetSt) _put_ref_to_way(ref string, after_node_id ...string) error {
	switch ChSet.OsmCh.Type {
	case "modify":
		return ChSet.OsmCh.Modify.Way._put_ref_to_way(ref, after_node_id...)
	case "create":
		return ChSet.OsmCh.Create.Way._put_ref_to_way(ref, after_node_id...)
	case "delete":
		return ChSet.OsmCh.Delete.Way._put_ref_to_way(ref, after_node_id...)
	}

	return errors.New("_put_ref_to_way. No way changeset")
}

func (w *WaySt) _put_ref_to_way(NodeId string, after_node_id ...string) error {

	w._check_nodes()

	if len(after_node_id) == 0 {
		w.Nodes = append(w.Nodes, &WayNdSt{NodeId})
		return nil
	}

	nds := []*WayNdSt{}
	if after_node_id[0] == "0" {
		nds = append(nds, &WayNdSt{NodeId})
		w.Nodes = append(nds, w.Nodes...)
		return nil
	}

	for _, v := range w.Nodes {
		nds = append(nds, v)
		if v.Ref == after_node_id[0] {
			nds = append(nds, &WayNdSt{NodeId})
		}
	}

	if len(nds) == len(w.Nodes) {
		return errors.New("_put_ref_to_way. Node not found")
	}
	w.Nodes = nds
	return nil
}

//======================
