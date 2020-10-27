package main

import (
	"encoding/json"
	"log"
	"strings"
)

// Rule contains all information needed to determine which Release an update request should receive. Specifically:
// * All of the gothmogFields, which are used to determine which Rules are relevant to an update request
// * A priority, which is used to tie-break when multiple matching Rules exist
// * A release_mapping, which is the Release that is returned when the most matching Rule is found
type Rule struct {
	properties      gothmogFields
	release_mapping string
	priority        int
}

// TODO: is this even useful to define?
type Rules []Rule

// balrogRules is an intermediate structure that contains all fields that Balrog's
// Rules have.
type balrogRule struct {
	BuildID        string `json:"buildID"`
	BuildTarget    string `json:"buildTarget"`
	Channel        string `json:"channel"`
	DistVersion    string `json:"distVersion"`
	Distribution   string `json:"distribution"`
	InstructionSet string `json:"instructionSet"`
	Locale         string `json:"locale"`
	Mapping        string `json:"mapping"`
	Memory         string `json:"memory"`
	OsVersion      string `json:"osVersion"`
	Priority       int    `json:"priority"`
	Product        string `json:"product"`
	Version        string `json:"version"`
}

// parseRules transforms Balrog Rules into Gothmog Rules
// by plucking out the parts of the Balrog Rules we care about.
func parseRules(data []byte) (Rules, error) {
	var importedRules []balrogRule
	var parsedRules Rules
	err := json.Unmarshal(data, &importedRules)
	if err != nil {
		return parsedRules, err
	}

	for _, rule := range importedRules {
		parsedRules = append(parsedRules, Rule{
			properties: gothmogFields{
				product:        rule.Product,
				version:        rule.Version,
				buildid:        rule.BuildID,
				buildTarget:    rule.BuildTarget,
				locale:         rule.Locale,
				channel:        rule.Channel,
				osVersion:      rule.OsVersion,
				instructionSet: rule.InstructionSet,
				memory:         rule.Memory,
				distribution:   rule.Distribution,
				distVersion:    rule.DistVersion,
			},
			release_mapping: rule.Mapping,
			priority:        rule.Priority,
		})
	}

	return parsedRules, nil
}

// matchCsv determines whether or not any of the comma separated
// values of `field` match `value`. `substring` controls whether
// a full or partial string match is performed.
func matchCsv(field string, value string, substring bool) bool {
	for _, f := range strings.Split(field, ",") {
		if substring {
			// TODO: do a substring match!
			return false
		} else {
			if f == value {
				return true
			}
		}
	}

	return false
}

// matchComparison tests whether `value` matches the
// test in `field`. `field` may begin be a plain string
// or begin with <, <=, >, or >=. If the latter, the operator
// given is used to compare `value` against the non-operator
// portion of `field`
func matchComparison(field string, value string) bool {
	// TODO: letting this anonymous function update `prefix`
	// is a little nasty. It would be better if we could split
	// by the list of all prefixes instead.
	var prefix string
	f := func(c rune) bool {
		is_prefix := c == '<' || c == '>' || c == '='
		if is_prefix {
			prefix = prefix + string(c)
		}
		return is_prefix
	}
	field_value := strings.FieldsFunc(field, f)[0]

	if prefix == "" {
		if prefix == value {
			return true
		}
		return false
	}

	// TODO: there must be a better way to do this.
	// something like python's `operator` library maybe?
	if prefix == "<" && value < field_value {
		return true
	}
	if prefix == "<=" && value <= field_value {
		return true
	}
	if prefix == ">" && value > field_value {
		return true
	}
	if prefix == ">=" && value >= field_value {
		return true
	}

	return false
}

// findMatchingRule compares an incoming request against a set of
// Rules and returns the best matching Rule.
// TODO: this needs tests!
func findMatchingRule(rules *Rules, req gothmogFields) Rule {
	// TODO: this should be define outside of the function as a general
	// sentinel value
	var matchingRule Rule
	matchingRule.priority = -1

	for _, rule := range *rules {
		if rule.properties.product != "" && rule.properties.product != req.product {
			continue
		}
		// TODO: support version comparison and csv
		if rule.properties.version != "" && !matchComparison(rule.properties.version, req.version) {
			continue
		}
		// TODO: support version comparison
		if rule.properties.buildid != "" && rule.properties.buildid != req.buildid {
			continue
		}
		// TODO: support comma separated values
		if rule.properties.buildTarget != "" && !matchCsv(rule.properties.buildid, req.version, false) {
			continue
		}
		// TODO: support comma separated values
		if rule.properties.locale != "" && rule.properties.locale != req.locale {
			continue
		}
		// TODO: support * globbing
		if rule.properties.channel != "" && rule.properties.channel != req.channel {
			continue
		}
		// support comma separated values and partial matches
		if rule.properties.osVersion != "" && rule.properties.osVersion != req.osVersion {
			continue
		}
		if rule.properties.instructionSet != "" && rule.properties.instructionSet != req.instructionSet {
			continue
		}
		// TODO: support comparison
		if rule.properties.memory != "" && rule.properties.memory != req.memory {
			continue
		}
		if rule.properties.distribution != "" && rule.properties.distribution != req.distribution {
			continue
		}
		if rule.properties.distVersion != "" && rule.properties.distVersion != req.distVersion {
			continue
		}

		if rule.priority > matchingRule.priority {
			log.Printf("Replacing matchingRule %v with %v", matchingRule, rule)
			matchingRule = rule
		}
	}

	return matchingRule
}
