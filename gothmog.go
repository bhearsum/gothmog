package main

import (
    "fmt"
    "net/http"
    "strings"
)

type GothmogHandler struct {
}

type morgothFields struct {
    product string
    version string
    buildid string
    buildTarget string
    locale string
    channel string
    osVersion string
    instructionSet string // split out of the systemCapabilities field
    memory string // split out of the systemCapabilities field
    distribution string
    distVersion string
}

func splitFields(fields string) (morgothFields, bool) {
    sections := strings.Split(fields, "/")
    if len(sections) != 14 {
        return morgothFields{}, false
    }

    var instructionSet, memory string
    systemCapabilities := strings.Split(sections[10], ",")
    switch len(systemCapabilities) {
    case 0:
        ;
    case 1:
        instructionSet = systemCapabilities[0]
    default:
        instructionSet = systemCapabilities[0]
        memory = systemCapabilities[1]
    }
    return morgothFields{
        product: sections[3],
        version: sections[4],
        buildid: sections[5],
        buildTarget: sections[6],
        locale: sections[7],
        channel: sections[8],
        osVersion: sections[9],
        instructionSet: instructionSet,
        memory: memory,
        distribution: sections[11],
        distVersion: sections[12],
    }, true
}

func (b *GothmogHandler) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
    fields, ok := splitFields(strings.Split(req.URL.RequestURI(), "?")[0])
    if ok != true {
        rw.Header().Set("Content-Type", "text/plain")
        rw.Write([]byte("wrong number of fields"))
    } else {
        rw.Header().Set("Content-Type", "text/plain")
        rw.Write([]byte(fmt.Sprintf("%v", fields)))
    }
}