package main

import (
	"encoding/json"
	"fmt"
	"slices"
)

var TokenEnders = [4]string{":", ",", "}", "]"}
var TokenMods = [2]string{"{", "["}

const (
	VALUE  = iota
	OBJECT = iota
	ARRAY  = iota
)

type JsonEntry struct {
	jsonType     int
	length       int
	ObjectKeys   *[]JsonEntry
	ObjectValues *[]JsonEntry
	Array        *[]JsonEntry
	Value        *[]byte
}

func (je JsonEntry) String() string {
	result := ""
	switch je.jsonType {
	case VALUE:
		result += string(*je.Value)
	case OBJECT:
		result += "{"
		for idx, key := range *je.ObjectKeys {
			if idx != 0 {
				result += ","
			}
			val := (*je.ObjectValues)[idx]
			result += key.String()
			result += ":"
			result += val.String()
		}
		result += "}"
	case ARRAY:
		result += "["
		for idx, val := range *je.Array {
			//fmt.Println(idx, val)
			if idx != 0 {
				result += ","
			}
			result += val.String()
		}
		result += "]"
	}
	return result
}
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

func getArray(source []byte) JsonEntry {
	tokens := make([]JsonEntry, 0)
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

func getToken(source []byte) []byte {
	for idx, b := range source {
		if slices.Contains(TokenEnders[:], string(b)) {
			fmt.Println(string(b))
			if b == ']' {
			}
			fmt.Println(source[:idx])
			return source[:idx]
		}
	}
	return source
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
	item1.Beliefs = append(item1.Beliefs, "jesus christ")
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
	//fmt.Println(string(json1), string(json2))
}
