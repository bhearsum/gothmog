package main

import "encoding/json"

// Rule contains all information needed to determine which Release an update request should receive. Specifically:
// * All of the gothmogFields, which are used to determine which Rules are relevant to an update request
// * A priority, which is used to tie-break when multiple matching Rules exist
// * A release_mapping, which is the Release that is returned when the most matching Rule is found
type Rule struct {
	properties      gothmogFields
	release_mapping string
	priority        int
}
type Rules []Rule

// balrogRules is an intermediate structure that contains all fields that Balrog's
// Rules have.
type balrogRules []struct {
	// TODO: Is it possible to only define the fields that we're actually going to use?
	alias              string
	backgroundRate     int
	buildID            string
	buildTarget        string
	channel            string
	comment            string
	data_version       int
	distVersion        string
	distribution       string
	headerArchitecture string
	instructionSet     string
	jaws               bool
	locale             string
	mapping            string
	memory             string
	mig64              bool
	osVersion          string
	priority           int
	product            string
	rule_id            int
	update_type        string
	version            string
}

// parseRules transforms Balrog Rules into Gothmog Rules
// by plucking out the parts of the Balrog Rules we care about.
func parseRules(data []byte) (Rules, error) {
	var importedRules balrogRules
	var parsedRules Rules
	err := json.Unmarshal(data, &importedRules)
	if err != nil {
		return parsedRules, err
	}

	for _, rule := range importedRules {
		parsedRules = append(parsedRules, Rule{
			properties: gothmogFields{
				product:        rule.product,
				version:        rule.version,
				buildid:        rule.buildID,
				buildTarget:    rule.buildTarget,
				locale:         rule.locale,
				channel:        rule.channel,
				osVersion:      rule.osVersion,
				instructionSet: rule.instructionSet,
				memory:         rule.memory,
				distribution:   rule.distribution,
				distVersion:    rule.distVersion,
			},
			release_mapping: rule.mapping,
			priority:        rule.priority,
		})
	}

	return parsedRules, nil
}
