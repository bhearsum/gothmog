package main

import (
	"io/ioutil"
	"log"
	"testing"
)

func TestParseRules(t *testing.T) {
	tests := map[string]struct {
		file string
		err  bool
	}{
		"good rules": {
			file: "testdata/good_rules.json",
			err:  false,
		},
		"bad rules": {
			file: "testdata/bad_rules.json",
			err:  true,
		},
	}

	for name, testcase := range tests {
		data, err := ioutil.ReadFile(testcase.file)
		if err != nil {
			t.Errorf("%v failed when reading %v: %v", name, testcase.file, testcase.err)
			continue
		}

		rules, err := parseRules(data)
		if testcase.err {
			if err == nil {
				t.Errorf("%v should've failed but didn't when parsing rules %v", name, testcase.file)
				continue
			}
		} else if err != nil {
			t.Errorf("%v failed when parsing rules: %v", name, err)
			continue
		}

		for _, rule := range rules {
			if rule.release_mapping == "" {
				t.Errorf("%v failed: no mapping found for rule: %v", name, rule)
				break
			}
		}
	}
}

// TestFindMatchingRule is not necessarily the ideal way to test
// findMatchingRule, as we don't fully isolate each individual
// property under test - we test the entire function at once.
// In reality, there are never rules with just one property set,
// and it's probably better to test using realistic rules than
// fully isolate each property under test.
// The main downside here is that a bug in evaluating one
// property may cause seemingly unrelated tests to fail.
func TestFindMatchingRule(t *testing.T) {
	data, err := ioutil.ReadFile("testdata/good_rules.json")
	if err != nil {
		t.Errorf("Couldn't read rules to run rule matching tests")
		return
	}

	rules, err := parseRules(data)
	if err != nil {
		t.Errorf("Couldn't parse rules to run rule matching tests")
		return
	}

	tests := map[string]struct {
		req  gothmogFields
		want Rule
	}{
		"simple string matches": {
			req: gothmogFields{
				product:        "Firefox",
				channel:        "aurora",
				instructionSet: "SSE",
				osVersion:      "Linux",
			},
			want: rules[46],
		},
		"simple string no matches": {
			req: gothmogFields{
				product:      "NotFirefox",
				channel:      "fake",
				osVersion:    "fake",
				distribution: "fake",
			},
			want: Rule{
				priority: -1,
			},
		},
		"version less than": {
			req: gothmogFields{
				product: "Firefox",
				channel: "esr",
				version: "78.0",
			},
			want: rules[5],
		},
	}

	for name, testcase := range tests {
		log.Printf("Running test: %v", name)
		got := findMatchingRule(&rules, testcase.req)
		if got != testcase.want {
			t.Errorf("%v failed. wanted: %v, got: %v", name, testcase.want, got)
			continue
		}
	}
}
