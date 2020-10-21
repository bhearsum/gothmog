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
type balrogRule struct {
	BuildID            string `json:"buildID"`
	BuildTarget        string `json:"buildTarget"`
	Channel            string `json:"channel"`
	DistVersion        string `json:"distVersion"`
	Distribution       string `json:"distribution"`
	InstructionSet     string `json:"instructionSet"`
	Locale             string `json:"locale"`
	Mapping            string `json:"mapping"`
	Memory             string `json:"memory"`
	OsVersion          string `json:"osVersion"`
	Priority           int    `json:"priority"`
	Product            string `json:"product"`
	Version            string `json:"version"`
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
