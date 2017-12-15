package main

import (
	"io/ioutil"
	"log"
	"os/exec"
)

// Chrome service
type Chrome struct {
	chromePath string
	args       []string
	source     string
	pdf        pdf
}

//NewChrome Instanciate a new Chrome instance
func NewChrome(src string, name string) Chrome {
	pdf := pdf{
		"storage/pdf/",
		name + ".pdf",
	}
	return Chrome{
		chromePath: "/usr/bin/google-chrome",
		args: []string{
			"--headless",
			"--disable-gpu",
			"--print-to-pdf=" + pdf.GetPath(),
			src,
		},
		source: src,
		pdf:    pdf,
	}
}

func (c *Chrome) run() pdf {
	_, err := exec.Command(c.chromePath, c.args...).Output()
	if err != nil {
		log.Fatal(err)
	}

	return c.pdf
}

//Pdf model
type pdf struct {
	root     string
	filename string
}

//GetPath return the absolute pasth to the generated pdf
func (p *pdf) GetPath() string {
	return p.root + p.filename
}

func (p *pdf) GetContent() []byte {
	content, err := ioutil.ReadFile(p.GetPath())
	if err != nil {
		panic(err)
	}
	return content
}
