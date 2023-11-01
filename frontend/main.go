package main

import (
	"encoding/xml"
	"fmt"
	"os"
)

type Text struct {
	XMLName xml.Name `xml:"Text"`
	Value   string   `xml:",chardata"`
}

type View struct {
	XMLName xml.Name `xml:"View"`
	Text    Text     `xml:"Text"`
}

type Version struct {
	XMLName xml.Name `xml:"Version"`
	Value   string   `xml:",chardata"`
}

type CInter struct {
	XMLName xml.Name `xml:"CInter"`
	View    View     `xml:"View"`
	Version Version  `xml:"Version"`
}

func main() {
	data, err := os.ReadFile("/Users/victorbravo/GolandProjects/stories/frontend/view.jxm")
	if err != nil {
		panic(err)
	}
	var jxm CInter
	err = xml.Unmarshal(data, &jxm)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%v", jxm.View.Text.Value)
}
