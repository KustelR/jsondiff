package main

import (
	"encoding/json"
	"fmt"
)

/*
	[WIP]!

Finds difference between two JSONs and returns two new, first containing entries missing from target, and second - lines missing from source
*/
func Diff(source []byte, target []byte) ([]byte, []byte) {
	added := make([]byte, 0)
	deleted := make([]byte, 0)
	srcTokens := getTokens(source)
	fmt.Println(srcTokens)
	return added, deleted
}

func getTokens(source []byte) JsonEntry {
	first := source[0]
	var jsonType int
	switch first {
	case '{':
		return getObject(source)
	case '[':
		return getArray(source)
	default:
		value := getToken(source)
		return JsonEntry{jsonType, len(value), nil, nil, nil, &value}
	}
}

	srcLen := len(source)
	last := 0
	for idx := 1; idx < srcLen; idx++ {
		if idx >= last {

		}
		char := source[idx]
		if char == ']' {
			return JsonEntry{ARRAY, last, nil, nil, &tokens, nil}
		}

		tokenData := getTokens(source[idx:])
		tokens = append(tokens, tokenData)
		idx += tokenData.length
		last = idx
	}
	return JsonEntry{ARRAY, last, nil, nil, &tokens, nil}
}

func getObject(source []byte) JsonEntry {
	last := 0
	srcLen := len(source)
	tokens := make([]JsonEntry, 0)
	keys := make([]JsonEntry, 0)
	values := make([]JsonEntry, 0)
	for idx := 1; idx < srcLen; idx++ {
		tokenData := getTokens(source[idx:])
		tokens = append(tokens, tokenData)
		idx += tokenData.length
		last = idx
	}
	for idx, tkn := range tokens {
		if idx%2 == 0 {
			keys = append(keys, tkn)
		} else {
			values = append(values, tkn)
		}
	}
	return JsonEntry{OBJECT, last, &keys, &values, nil, nil}
}

		}
	}
	return new, old
}

/*
Pushes identifier and value to dest
*/
func pushEntry(dest []byte, identifier []byte, value []byte) []byte {

	idCopy := make([]byte, len(identifier))
	valueCopy := make([]byte, len(value))
	copy(idCopy, identifier)
	copy(valueCopy, value)
	dest = append(dest, idCopy...)
	dest = append(dest, valueCopy...)

	return dest
}

type smth struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Eatable     bool     `json:"eatable"`
	Beliefs     []string `json:"beliefs"`
}

func main() {
	item1 := smth{
		"john",
		"humans",
		true,
		make([]string, 0),
	}

	item2 := smth{
		"jown",
		"humane",
		true,
		make([]string, 0),
	}
	item1.Beliefs = append(item1.Beliefs, "noodle monster")
	item2.Beliefs = append(item2.Beliefs, "noodle monster")
	json1, err := json.Marshal(item1)
	if err != nil {
		panic(err)
	}
	json2, err := json.Marshal(item2)
	if err != nil {
		panic(err)
	}
	old, newer := Diff(json1, json2)
	fmt.Println(string(newer))
	fmt.Println(string(old))
	fmt.Println(string(json1), string(json2))
}
