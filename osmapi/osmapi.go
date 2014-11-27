package osmapi

/*
	API v0.6
	For more infomation see http://wiki.openstreetmap.org/wiki/API_v0.6
*/

import (
	"bytes"
	"encoding/xml"
	"errors"
	"gopkg.in/xmlpath.v2"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
)

const MainURL = "http://api.openstreetmap.org"
const MainURLTest = "http://api06.dev.openstreetmap.org"
const ProtocolVersion = "0.6"
const TimeFormatLayout = "2006-01-02T15:04:05-07:00"
const UserAgent = "Sputnik.Ru.Adminka" // Default user agent :)
const User = "TestUser"                // Default user agent :)
const UserId = "TestUser"              // Default user agent :)

type myjar struct {
	jar map[string][]*http.Cookie
}

func (p *myjar) SetCookies(u *url.URL, cookies []*http.Cookie) {
	p.jar[u.Host] = cookies
}

func (p *myjar) Cookies(u *url.URL) []*http.Cookie {
	return p.jar[u.Host]
}

type MyRequestSt struct {
	User      string
	Pass      string
	Url       string
	UserAgent string
	Debug     bool
	Jar       *myjar
}

type Miscellaneous struct {
	Version     map[string]float64
	Area        float64
	Tracepoints int
	Waynodes    int
	Changesets  int
	Timeout     int
	Status      map[string]string
}

func (m *MyRequestSt) SetDebug(v ...bool) {
	if len(v) > 0 {
		m.Debug = v[0]
	} else {
		m.Debug = true
	}
}

func (m *MyRequestSt) UserPass(user, pass string) {
	m.User = user
	m.Pass = pass
}

func MyRequest(auths ...string) *MyRequestSt {

	m := MyRequestSt{}

	if len(auths) > 0 {
		m.User = auths[0]
	}
	if len(auths) > 1 {
		m.Pass = auths[1]
	}

	m.Jar = &myjar{}
	m.Jar.jar = make(map[string][]*http.Cookie)
	m.Url = MainURL

	return &m
}

func (m *MyRequestSt) SetUrl(url string) {
	m.Url = url
}

func (m *MyRequestSt) Generator(user_agent string) {
	m.UserAgent = user_agent
}

func (m *MyRequestSt) makeSendRequest(Type, Url string, Content ...string) (string, error) {

	var err error
	var req *http.Request
	Url = m.Url + Url

	if m.Debug {
		log.Printf("Type = %s,  Url = %s \n", Type, Url)
	}

	if len(Content) > 0 {
		if m.Debug {
			log.Printf("\nsendContent = %s\n", Content[0])
		}
		str := []byte(Content[0])
		req, err = http.NewRequest(Type, Url, bytes.NewBuffer(str))
		req.Header.Set("Content-Type", "text/xml")
	} else {
		req, err = http.NewRequest(Type, Url, nil)
	}

	if err != nil {
		if m.Debug {
			log.Printf("Type: %s, Url: %s, error: %s \n", Type, Url, err)
		}
		return "", err
	}

	if m.User != "" {
		if m.Debug {
			log.Printf("Set SetBasicAuth. User: %s, Pass %s\n", m.User, m.Pass)
		}
		req.SetBasicAuth(m.User, m.Pass)
	}

	req.Header.Set("User-Agent", m.UserAgent)

	client := &http.Client{}
	client.Jar = m.Jar

	res, err_d := client.Do(req)
	if err_d != nil {
		if m.Debug {
			log.Printf("Type: %s, Url: %s, error: %s \n", Type, Url, err_d)
		}
		return "", err_d
	}

	body, err_r := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err_r != nil {
		if m.Debug {
			log.Printf("Type: %s, Url: %s, error: %s \n", Type, Url, err_r)
		}
		return "", err_r
	}

	if res.StatusCode != 200 {
		if m.Debug {
			log.Printf("Type: %s, Url: %s, StatusCode: %d, error: %s \n", Type, Url, res.StatusCode, string(body))
		}
		return "", errors.New(string(body))
	}

	if m.Debug {
		log.Printf("\n------> Read Content\n" + string(body) + "\nRead Content <------\n")
	}

	return string(body), nil
}

func (m *MyRequestSt) Get(url string, type_req ...string) (string, error) {
	return m.makeSendRequest("GET", url)
}

func (m *MyRequestSt) GetXML(url string, type_req ...string) (*xmlpath.Node, error) {

	body, err := m.Get(url, type_req...)

	if err != nil {
		return nil, err
	}

	return xmlpath.ParseDecoder(xml.NewDecoder(strings.NewReader(body)))
}

func (m *MyRequestSt) Post(url string, type_req ...string) (string, error) {
	return m.makeSendRequest("POST", url, type_req...)
}

func (m *MyRequestSt) PostXML(url string, type_req ...string) (*xmlpath.Node, error) {

	body, err := m.Post(url, type_req...)

	if err != nil {
		return nil, err
	}

	return xmlpath.ParseDecoder(xml.NewDecoder(strings.NewReader(body)))
}

func (m *MyRequestSt) Put(url string, type_req ...string) (string, error) {
	return m.makeSendRequest("PUT", url, type_req...)
}

func (m *MyRequestSt) PutXML(url string, type_req ...string) (*xmlpath.Node, error) {

	body, err := m.Put(url, type_req...)

	if err != nil {
		return nil, err
	}

	return xmlpath.ParseDecoder(xml.NewDecoder(strings.NewReader(body)))
}

func xml_str(data *xmlpath.Node, where string) string {

	path := xmlpath.MustCompile(where)
	val, ok := path.String(data)

	if !ok {
		return ""
	}

	return val
}

func xml_slice(data *xmlpath.Node, where string, what []string) []map[string]string {

	path := xmlpath.MustCompile(where)
	iters := path.Iter(data)
	out := []map[string]string{}

	for iters.Next() {
		one := iters.Node()

		ph := map[string]string{}
		for _, tag := range what {
			ph[tag] = xml_str(one, "@"+tag)
		}

		out = append(out, ph)
	}

	return out
}
