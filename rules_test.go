package main

import (
    "reflect"
    "testing"
)

func TestParseRule(t *testing.T) {
    tests := map[string]struct {
        rawRule map[string]interface{}
        want Rule
    } {
        "simple": {
            rawRule: map[string]interface{}{
                "product": "Firefox",
                "mapping": "Firefox-58.0-build1",
                "channel": "release",
                "priority": 100,
            },
            want: Rule{
                properties: gothmogFields{
                    product: "Firefox",
                    channel: "release",
                },
                release_mapping: "Firefox-58.0-build1",
                priority: 100,
            },
        },
        "fallback mapping ignored": {
            rawRule: map[string]interface{}{
                "product": "Firefox",
                "mapping": "Firefox-58.0-build1",
                "channel": "release",
                "priority": 100,
                "fallbackMapping": "ignored",
            },
            want: Rule{
                properties: gothmogFields{
                    product: "Firefox",
                    channel: "release",
                },
                release_mapping: "Firefox-58.0-build1",
                priority: 100,
            },
        },
    }

    for name, testcase := range tests {
        got, _ := parseRule(testcase.rawRule)

        if !reflect.DeepEqual(testcase.want, got) {
            t.Errorf("%v: Parsed rule does not match expected value: %v, %v", name, got, testcase.want)
        }
    }
}
