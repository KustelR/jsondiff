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
	srcLen := len(source)
	tgtLen := len(target)
	if srcLen != tgtLen {
		panic("different length not supported")
	}
	srcTokens := getTokens(source)
	tgtTokens := getTokens(target)
	new := make([]byte, tgtLen)
	old := make([]byte, srcLen)
	var identifier []byte
	for idx, srcToken := range srcTokens {
		tgtToken := tgtTokens[idx]
		isSame := compareTokens(srcToken, tgtToken)
		fmt.Println(isSame, "id:", string(identifier), "values:", string(srcToken), string(tgtToken))
		if identifier != nil {
			if isSame {
				identifier = nil
				continue
			}
			new = pushEntry(new, identifier, srcToken)
			old = pushEntry(old, identifier, srcToken)
			identifier = nil

		} else {
			if isSame {
				identifier = srcToken
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
