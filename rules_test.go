package main

import (
    "io/ioutil"
    "testing"
)

func TestParseRules(t *testing.T) {
    tests := map[string]struct {
        file string
        err error
    } {
        "good rules": {
            file: "fixtures/good_rules.json",
            err: nil,
        },
    }

    for name, testcase := range tests {
        data, err := ioutil.ReadFile(testcase.file)
        if err != testcase.err {
            t.Errorf("%v failed when reading %v: %v", name, testcase.file, testcase.err)
        }

        _, err = parseRules(data)
        if err != nil {
            t.Errorf("%v failed when parsing rules: %v", name, testcase.err)
        }
    }
}
