package upnp

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"encoding/xml"
)

var serviceURNRegexp *regexp.Regexp

func init() {
	var err error
	serviceURNRegexp, err = regexp.Compile(`^urn:schemas-upnp-org:service:(\w+):(\d+)$`)
	if err != nil {
		panic(err)
	}
}

type ServiceURN struct {
	Type    string
	Version uint64
}

func (me ServiceURN) String() string {
	return fmt.Sprintf("urn:schemas-upnp-org:service:%s:%d", me.Type, me.Version)
}

func ParseServiceType(s string) (ret ServiceURN, ok bool) {
	matches := serviceURNRegexp.FindStringSubmatch(s)
	if matches == nil {
		return
	}
	if len(matches) != 3 {
		panic(matches)
	}
	ret.Type = matches[1]
	var err error
	ret.Version, err = strconv.ParseUint(matches[2], 0, 0)
	if err != nil {
		return
	}
	return
}

type SoapAction struct {
	ServiceURN
	Action string
}

func ParseActionHTTPHeader(s string) (ret SoapAction, ok bool) {
	if s[0] != '"' || s[len(s)-1] != '"' {
		return
	}
	s = s[1 : len(s)-1]
	hashIndex := strings.LastIndex(s, "#")
	if hashIndex == -1 {
		return
	}
	ret.Action = s[hashIndex+1:]
	ret.ServiceURN, ok = ParseServiceType(s[:hashIndex])
	ok = true
	return
}

type SpecVersion struct {
	Major int `xml:"major"`
	Minor int `xml:"minor"`
}

type icon struct {
	Mimetype, Width, Height, Depth, URL string
}

type Service struct {
	XMLName     xml.Name `xml:"service"`
	ServiceType string   `xml:"serviceType"`
	ServiceId   string   `xml:"serviceId"`
	SCPDURL     string
	ControlURL  string `xml:"controlURL"`
	EventSubURL string `xml:"eventSubURL"`
}

type Device struct {
	DeviceType   string `xml:"deviceType"`
	FriendlyName string `xml:"friendlyName"`
	Manufacturer string `xml:"manufacturer"`
	ModelName    string `xml:"modelName"`
	UDN          string
	IconList     []icon
	ServiceList  []Service `xml:"serviceList>service"`
}

type DeviceDesc struct {
	XMLName     xml.Name    `xml:"urn:schemas-upnp-org:device-1-0 root"`
	SpecVersion SpecVersion `xml:"specVersion"`
	Device      Device      `xml:"device"`
}