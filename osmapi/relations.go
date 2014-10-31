package osmapi

import (
	"errors"
	"time"
)

/*

Relation create

<?xml version="1.0" encoding="UTF-8"?>
<osm version="0.6" generator="CGImap 0.3.3 (24667 thorn-03.openstreetmap.org)" \
		copyright="OpenStreetMap and contributors" attribution="http://www.openstreetmap.org/copyright" \
		license="http://opendatacommons.org/licenses/odbl/1-0/">
 <relation id="2996185" visible="true" version="1" changeset="16527892" \
 	timestamp="2013-06-12T18:49:10Z" user="landfahrer" uid="1069176">
  <member type="relation" ref="224645265" role="outer"/>
  <member type="relation" ref="225456073" role="inner"/>
  <tag k="type" v="multipolygon"/>
 </relation>
</osm>

*/

/* ===
Relations in changesets
*/

type MemberSt struct {
	Ref  string `xml:"ref,attr"`
	Type string `xml:"type,attr"`
	Role string `xml:"role,attr"`
}

type RelationSt struct {
	Tags      []*TagSt    `xml:"tag"`
	Members   []*MemberSt `xml:"member"`
	OsmId     string      `xml:"id,attr"`
	ReqId     string      `xml:"changeset,attr"`
	Timestamp string      `xml:"timestamp,attr"`
	Uid       string      `xml:"uid,attr,omitempty"`
	User      string      `xml:"user,attr,omitempty"`
	Version   string      `xml:"version,attr,omitempty"`
	Visible   string      `xml:"visible,attr,omitempty"`
	IsNew     bool        `xml:"-"`
}

/*
When we want to modify or delete node we have get infomation from api.site
*/
func (ChSet *ChangeSetSt) RelationNew() (*RelationSt, error) {

	r := RelationSt{}
	r.Tags = nil
	r.Members = nil
	r.ReqId = ChSet.Id
	r.Version = "1"
	r.OsmId = "-1"
	r.Visible = "true"
	r.IsNew = true

	tm := time.Now()
	r.Timestamp = tm.Format(TimeFormatLayout)

	ChSet.OsmCh._setRelation(&r)
	ChSet.OsmCh.ChangeType = "relation"

	return &r, nil
}

/*
When we want to modify or delete node we have get infomation from api.site
*/
func (ChSet *ChangeSetSt) RelationLoad(OsmId string) (*RelationSt, error) {

	/* Answer has to be empty */
	data, err := ChSet.Request.GetXML("/api/0.6/relation/" + OsmId)
	if err != nil {
		return nil, err
	}

	r, err_r := ChSet.RelationNew()
	if err_r != nil {
		return nil, err_r
	}
	r.OsmId = OsmId
	r.Version = xml_str(data, "/osm/relation/@version")
	r.Visible = xml_str(data, "/osm/relation/@visible")
	r.Uid = xml_str(data, "/osm/relation/@uid")
	r.User = xml_str(data, "/osm/relation/@user")
	r.IsNew = false

	for _, v := range xml_slice(data, "/osm/relation/member", []string{"type", "ref", "role"}) {
		if v["ref"] == "" {
			continue
		}
		r._add_member(v["type"], v["ref"], v["role"])
	}

	for _, v := range xml_slice(data, "/osm/relation/tag", []string{"k", "v"}) {
		r._add_tag(v["k"], v["v"])
	}

	return r, nil
}

func (ChSet *ChangeSetSt) RelationAddMember(t, ref, role string) error {

	err := errors.New("RelationAddMember. Bad member type. Must be [way|node]. Now is " + t)
	switch t {
	case "node":
		_, err = ChSet.LoadNodeDate(ref)
	case "way":
		_, err = ChSet.WayLoadData(ref)
	}

	if err != nil {
		return err
	}

	return ChSet._add_member(t, ref, role)
}

func (ChSet *ChangeSetSt) _add_member(t, ref, role string) error {
	switch ChSet.OsmCh.Type {
	case "modify":
		return ChSet.OsmCh.Modify.Relation._add_member(t, ref, role)
	case "create":
		return ChSet.OsmCh.Create.Relation._add_member(t, ref, role)
	case "delete":
		return ChSet.OsmCh.Delete.Relation._add_member(t, ref, role)
	}

	return errors.New("_setRelation. No relation changeset")
}

func (r *RelationSt) _add_member(t, ref, role string) error {
	if t == "" || ref == "" {
		return errors.New("Relation. _add_member. Empty type or ref for member")
	}
	if r.Members == nil {
		r.Members = []*MemberSt{}
	}

	m := MemberSt{}
	m.Ref = ref
	m.Type = t
	m.Role = role
	r.Members = append(r.Members, &m)

	return nil
}

func (w *RelationSt) _add_tag(k, v string) {
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

func (OsmCh *OsmChangeSt) _setRelation(relation *RelationSt) error {
	switch OsmCh.Type {
	case "modify":
		OsmCh.Modify.Relation = relation
	case "create":
		OsmCh.Create.Relation = relation
	case "delete":
		OsmCh.Delete.Relation = relation
	}
	return errors.New("_setRelation. No relation changeset")

}

func (ChSet *ChangeSetSt) RelationDelMember(t, ref string) error {
	// Answer has to be empty
	switch ChSet.OsmCh.Type {
	case "modify":
		return ChSet.OsmCh.Modify.Relation._del_member(t, ref)
	case "create":
		return ChSet.OsmCh.Create.Relation._del_member(t, ref)
	case "delete":
		return ChSet.OsmCh.Delete.Relation._del_member(t, ref)
	}

	return errors.New("_del_ref_from_relation. No relation changeset")
}

func (ChSet *ChangeSetSt) RelationDelAllMembers() error {

	switch ChSet.OsmCh.Type {
	case "modify":
		ChSet.OsmCh.Modify.Relation.Members = nil
	case "create":
		ChSet.OsmCh.Create.Relation.Members = nil
	case "delete":
		ChSet.OsmCh.Delete.Relation.Members = nil
	default:
		return errors.New("RelationDelAllMembers. No relation changeset")
	}

	return nil
}

func (r *RelationSt) _del_member(t, ref string) error {

	if r.Members != nil {
		nds := []*MemberSt{}
		for _, v := range r.Members {
			if v.Ref != ref || v.Type != t {
				nds = append(nds, v)
			}
		}
		r.Members = nds
	}

	return nil
}
