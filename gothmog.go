package main

import (
	"fmt"
	"net/http"
	"strings"
)

type GothmogHandler struct {
	rules *Rules
}

// gothmogFields hold all of the relevant information contained in an update request URI
type gothmogFields struct {
	product        string
	version        string
	buildid        string // TODO: this could potentially be an int, not sure if it should be
	buildTarget    string
	locale         string
	channel        string
	osVersion      string
	instructionSet string // split out of the systemCapabilities field
	memory         string // split out of the systemCapabilities field
	distribution   string
	distVersion    string
}

// splitFields takes an entire update URI and parses the useful parts into a gothmogFields instance
// Example update URI: /update/6/Firefox/55.0/20170731163142/Linux_x86_64-gcc3/en-GB/beta/Linux 4.11.3-202.fc25.x86_64 (GTK 3.22.15,libpulse 10.0.0)/NA,8196/default/default/update.xml
// Notably, the first two parts and the final `update.xml` part are unused.
func splitFields(fields string) (gothmogFields, bool) {
	// TODO: make sure bad data is handled correctly - add tests

	sections := strings.Split(fields, "/")
	if len(sections) != 14 {
		return gothmogFields{}, false
	}

	// These two fields are extracted first to so the final return statement
	// can be a simple gothmogFields literal.
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

func (g *GothmogHandler) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	// Notably, we're throwing away query args here. In reality there are
	// a few that we should be paying attention to, but for this simple
	// implementation we're just ignoring them.
	fields, ok := splitFields(strings.Split(req.URL.RequestURI(), "?")[0])
	if ok != true {
		rw.Header().Set("Content-Type", "text/plain")
		rw.Write([]byte("Couldn't parse update URI"))
	}

	rule := findMatchingRule(g.rules, fields)
	rw.Header().Set("Content-Type", "text/plain")
	rw.Write([]byte(fmt.Sprintf("Matching rule is: %v", rule)))
}
