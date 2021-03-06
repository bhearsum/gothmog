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
			t.Errorf("%v failed when reading %v: %v", name, testcase.file, err)
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
				buildid:        "2020101010101010",
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
		"version <": {
			req: gothmogFields{
				product: "Firefox",
				channel: "esr",
				version: "78.0",
			},
			want: rules[5],
		},
		// TODO: test for version <= & >
		// But we don't have rules that use these right now
		"version > matches higher": {
			req: gothmogFields{
				product:   "Firefox",
				channel:   "beta*",
				version:   "45.0",
				osVersion: "GTK 2,GTK 3.0.,GTK 3.1.,GTK 3.2.,GTK 3.3.",
			},
			want: rules[120],
		},
		"version > matches exact": {
			req: gothmogFields{
				product:   "Firefox",
				channel:   "beta*",
				version:   "44.0",
				osVersion: "GTK 2,GTK 3.0.,GTK 3.1.,GTK 3.2.,GTK 3.3.",
			},
			want: rules[120],
		},
		"version matches csv": {
			req: gothmogFields{
				product:   "Firefox",
				channel:   "release*",
				version:   "56.0",
				osVersion: "Windows_NT",
			},
			want: rules[52],
		},
		"version matches csv2": {
			req: gothmogFields{
				product:   "Firefox",
				channel:   "release*",
				version:   "56.0.1",
				osVersion: "Windows_NT",
			},
			want: rules[52],
		},
		"version no match csv": {
			req: gothmogFields{
				product:   "Firefox",
				channel:   "release*",
				version:   "99.0",
				osVersion: "Windows_NT",
			},
			want: Rule{priority: -1},
		},
		"buildid <": {
			req: gothmogFields{
				product:   "Firefox",
				channel:   "nightly*",
				buildid:   "20151209095499",
				osVersion: "Windows_NT",
			},
			want: rules[108],
		},
		// TODO: tests for other operators for buildID, when
		// we have rules with them in it
		"buildTarget csv": {
			req: gothmogFields{
				product:     "Firefox",
				channel:     "bhrelease-localtest",
				version:     "56.0",
				buildTarget: "WINNT_x86-msvc-x64",
				memory:      ">2048",
			},
			want: rules[47],
		},
		"buildTarget csv2": {
			req: gothmogFields{
				product:     "Firefox",
				channel:     "bhrelease-localtest",
				version:     "56.0",
				buildTarget: "WINNT_x86_64-msvc",
				memory:      ">2048",
			},
			want: rules[47],
		},
		"locale csv": {
			req: gothmogFields{
				product: "Firefox",
				channel: "release*",
				version: "56.0",
				locale:  "as",
			},
			want: rules[33],
		},
		"locale csv2": {
			req: gothmogFields{
				product: "Firefox",
				channel: "release*",
				version: "56.0",
				locale:  "mai",
			},
			want: rules[33],
		},
		"locale csv3": {
			req: gothmogFields{
				product: "Firefox",
				channel: "release*",
				version: "56.0",
				locale:  "bn-IN",
			},
			want: rules[33],
		},
		"memory >": {
			req: gothmogFields{
				product:     "Firefox",
				channel:     "bhrelease-localtest",
				buildTarget: "WINNT_x86_64-msvc",
				version:     "56.0",
				memory:      "2049",
			},
			want: rules[47],
		},
		// TODO: test other operators when we have rules with them
		"osVersion csv": {
			req: gothmogFields{
				product:   "Firefox",
				channel:   "aurora",
				buildid:   "20170123004000",
				osVersion: "Windows_NT 5.1",
			},
			want: rules[58],
		},
		"osVersion csv2": {
			req: gothmogFields{
				product:   "Firefox",
				channel:   "aurora",
				buildid:   "20170123004000",
				osVersion: "Windows_NT 5.2",
			},
			want: rules[58],
		},
		"osVersion csv3": {
			req: gothmogFields{
				product:   "Firefox",
				channel:   "aurora",
				buildid:   "20170123004000",
				osVersion: "Windows_NT 6.0",
			},
			want: rules[58],
		},
		"osVersion substring match": {
			req: gothmogFields{
				product:   "Firefox",
				channel:   "esr*",
				osVersion: "Windows_NT 6.0",
				version:   "38.0",
			},
			want: rules[114],
		},
		"channel exact match": {
			req: gothmogFields{
				product: "Firefox",
				channel: "date-localtest",
			},
			want: rules[1],
		},
		"channel match glob": {
			req: gothmogFields{
				product:   "Firefox",
				channel:   "esr",
				version:   "38.0.0",
				osVersion: "Windows_NT",
			},
			want: rules[114],
		},
		"channel match glob2": {
			req: gothmogFields{
				product:   "Firefox",
				channel:   "esr-localtest",
				version:   "38.0.0",
				osVersion: "Windows_NT",
			},
			want: rules[114],
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
