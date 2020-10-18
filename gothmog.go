package main

import (
	"fmt"
	"net/http"
	"strings"
)

type GothmogHandler struct {
}

type gothmogFields struct {
	product        string
	version        string
	buildid        string
	buildTarget    string
	locale         string
	channel        string
	osVersion      string
	instructionSet string // split out of the systemCapabilities field
	memory         string // split out of the systemCapabilities field
	distribution   string
	distVersion    string
}

// TODO: make sure bad data is handled
func splitFields(fields string) (gothmogFields, bool) {
	sections := strings.Split(fields, "/")
	if len(sections) != 14 {
		return gothmogFields{}, false
	}

	var instructionSet, memory string
	systemCapabilities := strings.Split(sections[10], ",")
	switch len(systemCapabilities) {
	case 0:

	case 1:
		instructionSet = systemCapabilities[0]
	default:
		instructionSet = systemCapabilities[0]
		memory = systemCapabilities[1]
	}
	return gothmogFields{
		product:        sections[3],
		version:        sections[4],
		buildid:        sections[5],
		buildTarget:    sections[6],
		locale:         sections[7],
		channel:        sections[8],
		osVersion:      sections[9],
		instructionSet: instructionSet,
		memory:         memory,
		distribution:   sections[11],
		distVersion:    sections[12],
	}, true
}

func (b *GothmogHandler) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	fields, ok := splitFields(strings.Split(req.URL.RequestURI(), "?")[0])
	if ok != true {
		rw.Header().Set("Content-Type", "text/plain")
		rw.Write([]byte("wrong number of fields"))
    }

    fakeRule := make(map[string]interface{})
    fakeRule["product"] = "p"
    fakeRule["mapping"] = "m"
    fakeRule["priority"] = 99
    rule, ok := parseRule(fakeRule)
	if ok != true {
		rw.Header().Set("Content-Type", "text/plain")
		rw.Write([]byte("couldn't parse rule"))
	} else {
		rw.Header().Set("Content-Type", "text/plain")
		rw.Write([]byte(fmt.Sprintf("%v", fields)))
        rw.Write([]byte(fmt.Sprintf("%v", rule)))
	}
}
