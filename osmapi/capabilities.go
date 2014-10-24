package osmapi

import (
	//"github.com/iostrovok/go-iutils/iutils"
	"fmt"
	"gopkg.in/xmlpath.v2"
)

/* ===
Miscellaneous

Capabilities: GET /api/capabilities
This API call is meant to provide information about the capabilities and limitations of the current API.

*/

func (mr *MyRequestSt) MiscellCap() (*xmlpath.Node, error) {
	return mr.GetXML("/api/capabilities/")
}

/* ===

Retrieving permissions: GET /api/0.6/permissions
This API call is meant to provide information about the capabilities and limitations of the current API.

*/

func (mr *MyRequestSt) MiscellPermissions() (*xmlpath.Node, error) {
	return mr.GetXML("/api/0.6/permissions")
}

/* ===
Retrieving map data by bounding box: GET /api/0.6/map
*/

func (mr *MyRequestSt) MiscellMap(left string, bottom string, right string, top string) (*xmlpath.Node, error) {
	url := fmt.Sprintf("/api/0.6/map?bbox=%s,%s,%s,%s", left, bottom, right, top)
	return mr.GetXML(url)
}
