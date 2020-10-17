package main

type Rule struct {
    properties gothmogFields
    release_mapping string
    priority int
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
