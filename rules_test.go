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
			file: "fixtures/good_rules.json",
			err:  false,
		},
		"bad rules": {
			file: "fixtures/bad_rules.json",
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
	rules := Rules{
		Rule{
			properties: gothmogFields{
				product: "Firefox",
			},
			release_mapping: "Firefox-76.0-build1",
			priority: 100,
		},
		Rule{
			properties: gothmogFields{
				product: "Thunderbird",
			},
			release_mapping: "Thunderbird-56.0-build1",
			priority: 100,
		},
	}
	tests := map[string]struct {
		req   gothmogFields
		want Rule
	}{
		"product match": {
			req: gothmogFields{
				product: "Firefox",
			},
			want: rules[0],
		},
		"product no match": {
			req: gothmogFields{
				product: "NotFirefox",
			},
			want: Rule{
				priority: -1,
			},
		},
	}

	for name, testcase := range tests {
		got := findMatchingRule(&rules, testcase.req)
		if got != testcase.want {
			t.Errorf("%v failed. wanted: %v, got: %v", name, testcase.want, got)
			continue
		}
	}
}
