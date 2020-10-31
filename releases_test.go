package main

import (
	"io/ioutil"
	"testing"
)

func TestLoadRelease(t *testing.T) {
	tests := map[string]struct {
		file string
		err  bool
	}{
		"schema 1": {
			file: "testdata/releases/No-Update.json",
			err:  false,
		},
		"schema 2": {
			file: "testdata/releases/Firefox-12.0-build1.json",
			err:  false,
		},
		"schema 4": {
			file: "testdata/releases/Firefox-mozilla-central-nightly-latest.json",
			err:  false,
		},
		"schema 9": {
			file: "testdata/releases/Firefox-82.0-build2.json",
			err:  false,
		},
		"schema 50": {
			file: "testdata/releases/SSE-Desupport.json",
			err:  false,
		},
		// TODO: add invalid blobs
	}

	for name, testcase := range tests {
		data, err := ioutil.ReadFile(testcase.file)
		if err != nil {
			t.Errorf("%v failed when reading %v: %v", name, testcase.file, err)
			continue
		}

		release, err := parseRelease(data)
		if testcase.err {
			if err == nil {
				t.Errorf("%v should've failed but didn't when parsing release %v", name, testcase.file)
				continue
			}
		} else if err != nil {
			t.Errorf("%v failed when parsing release: %v", name, err)
			continue
		}

		if int(release["schema_version"].(float64)) < 1 {
			t.Errorf("%v failed: parsed release has no schema version set", name)
		}
	}
}
