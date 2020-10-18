package main

import "encoding/json"

type Rule struct {
    properties gothmogFields
    release_mapping string
    priority int
}

type Rules []struct {
    properties gothmogFields
    release_mapping string
    priority int
}

// We support importing Rules directly from Balrog's API. We don't use all of them,
// so this is just used as an intermediary until we get them into our proper Rules
// structure.
type balrogRules []struct {
    alias string
    backgroundRate int
    buildID string
    buildTarget string
    channel string
    comment string
    data_version int
    distVersion string
    distribution string
    headerArchitecture string
    instructionSet string
    jaws bool
    locale string
    mapping string
    memory string
    mig64 bool
    osVersion string
    priority int
    product string
    rule_id int
    update_type string
    version string
}

func parseRules(data []byte) (Rules, error) {
    var importedRules balrogRules
    var parsedRules Rules
    err := json.Unmarshal(data, &importedRules)
    if (err != nil) {
        return parsedRules, err
    }

    for _, rule := range importedRules {
        parsedRules = append(parsedRules, Rule{
            properties: gothmogFields{
                product: rule.product,
                version: rule.version,
                buildid: rule.buildID,
                buildTarget: rule.buildTarget,
                locale: rule.locale,
                channel: rule.channel,
                osVersion: rule.osVersion,
                instructionSet: rule.instructionSet,
                memory: rule.memory,
                distribution: rule.distribution,
                distVersion: rule.distVersion,
            },
            release_mapping: rule.mapping,
            priority: rule.priority,
        })
    }

    return parsedRules, nil
}

func getString(field interface{}) (string, bool) {
    val, ok := field.(string)
    if ok == true {
        return val, true
    } else {
        return "", false
    }
}

func getInt(field interface{}) (int, bool) {
    val, ok := field.(int)
    if ok == true {
        return val, true
    } else {
        return -1, false
    }
}

func parseRule(rawRule map[string]interface{}) (Rule, bool) {
    product, ok := getString(rawRule["product"])
    if ok == false {
        return Rule{}, ok
    }
    mapping, ok := getString(rawRule["mapping"])
    if ok == false {
        return Rule{}, ok
    }
    priority, ok := getInt(rawRule["priority"])
    if ok == false {
        return Rule{}, ok
    }

    return Rule{
        properties: gothmogFields{
            product: product,
        },
        release_mapping: mapping,
        priority: priority,
    }, true
}
