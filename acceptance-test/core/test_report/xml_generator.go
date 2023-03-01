package test_report

import "encoding/xml"

type xmlReport struct {
	XMLName xml.Name    `xml:"testsuites"`
	Suites  []*XmlSuite `xml:"testsuite,omitempty"`
}

type XmlSuite struct {
	Total    int        `xml:"tests,attr"`
	Failures int        `xml:"failures,attr"`
	Name     string     `xml:"name,attr"`
	Time     string     `xml:"time,attr"`
	Testcase []*XmlCase `xml:"testcase,omitempty"`
}

type XmlCase struct {
	Name    string   `xml:"name,attr"`
	Time    string   `xml:"time,attr"`
	Failure *Message `xml:"failure,omitempty"`
	Log     *Message `xml:"system-out,omitempty"`
}

type Message struct {
	Text string `xml:",cdata"`
}
