package main

import (
	"io/ioutil"
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
		req   gothmogFields
		want Rule
	}{
		"simple string matches": {
			req: gothmogFields{
				product: "Firefox",
				channel: "aurora",
				instructionSet: "SSE",
				osVersion: "Linux",
			},
			want: rules[46],
		},
		/*"product no match": {
			req: gothmogFields{
				product: "NotFirefox",
			},
			want: Rule{
				priority: -1,
			},
		},
		"version exact match": {
			req: gothmogFields{
				version: "58.0",
			},
			want: rules[100],
		},
		"version no match": {
			req: gothmogFields{
				version: "NotFirefox",
			},
			want: Rule{
				priority: -1,
			},
		},*/
	}

	for name, testcase := range tests {
		got := findMatchingRule(&rules, testcase.req)
		if got != testcase.want {
			t.Errorf("%v failed. wanted: %v, got: %v", name, testcase.want, got)
			continue
		}
	}
}
