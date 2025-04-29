package jsondiff

import (
	"slices"
)

var TokenEnders = [4]byte{':', ',', '}', ']'}
var TokenMods = [2]byte{'{', '['}

const (
	VALUE  = iota
	OBJECT = iota
	ARRAY  = iota
)

type jsonEntry struct {
	jsonType     int
	length       int
	ObjectKeys   *[]jsonEntry
	ObjectValues *[]jsonEntry
	Array        *[]jsonEntry
	Value        *[]byte
}

func (first *jsonEntry) Equal(other *jsonEntry) (bool, error) {
	if first.jsonType != other.jsonType {
		return false, nil
	}
	switch first.jsonType {
	case ARRAY:
		if len(*first.Array) != len(*other.Array) {
			return false, nil
		}
		for idx, val := range *first.Array {
			isEqual, err := val.Equal(&(*other.Array)[idx])
			if err != nil {
				return false, err
			}
			if !isEqual {
				return false, nil
			}
		}
		return true, nil
	case OBJECT:
		if (len(*first.ObjectValues)) >= len(*other.ObjectKeys) {
			return false, nil
		}
		for idx, key := range *first.ObjectKeys {
			isEqual, err := key.Equal(&(*other.ObjectKeys)[idx])
			if err != nil {
				return false, err
			}
			if !isEqual {
				return false, nil
			}
			isEqual, err = (*first.ObjectValues)[idx].Equal(&(*other.ObjectValues)[idx])
			if err != nil {
				return false, err
			}
			if !isEqual {
				return false, nil
			}
		}
		return true, nil
	default:
		return slices.Equal(*first.Value, *other.Value), nil
	}
}

func (je jsonEntry) String() string {
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
			if idx != 0 {
				result += ","
			}
			result += val.String()
		}
		result += "]"
	}
	return result
}

func halfDiffArray(source *jsonEntry, target *jsonEntry) []byte {
	result := make([]byte, 0)
	forAppend := make([]*jsonEntry, 0)
	for _, val := range *target.Array {
		isMatched := false
		for _, valTgt := range *source.Array {

			isMatched, _ = val.Equal(&valTgt)
		}
		if !isMatched {
			forAppend = append(forAppend, &val)
		}
	}
	if (len(forAppend)) <= 0 {
		return result
	}
	result = append(result, '[')
	for idx, val := range forAppend {
		if idx != 0 {
			result = append(result, ',')
		}
		result = append(result, []byte(val.String())...)
	}
	result = append(result, ']')
	return result
}

func halfDiffObject(source *jsonEntry, target *jsonEntry) []byte {
	result := make([]byte, 0)

	forAppend := make([][2]*jsonEntry, 0)
	for idx, key := range *target.ObjectKeys {
		isMatched := false
		for idxSrc, keySrc := range *source.ObjectKeys {
			isKeyMatch, _ := key.Equal(&keySrc)
			if !isKeyMatch {
				continue
			}
			isMatched, _ = (*target.ObjectValues)[idx].Equal(&(*source.ObjectValues)[idxSrc])
		}
		if !isMatched {
			forAppend = append(forAppend, [2]*jsonEntry{&key, &(*source.ObjectValues)[idx]})
		}

	}
	if (len(forAppend)) <= 0 {
		return result
	}
	result = append(result, '{')
	for idx := range forAppend {
		if idx != 0 {
			result = append(result, ',')
		}
		result = append(result, []byte(forAppend[idx][0].String())...)
		result = append(result, ':')
		result = append(result, []byte(forAppend[idx][1].String())...)
	}
	result = append(result, '}')

	return result
}

func Diff(source []byte, target []byte) ([]byte, []byte) {
	added := make([]byte, 0)
	deleted := make([]byte, 0)
	srcTokens := getTokens(source)
	tgtTokens := getTokens(target)

	if srcTokens.jsonType != tgtTokens.jsonType {
		return source, target
	}
	switch srcTokens.jsonType {
	case OBJECT:
		objAdded := halfDiffObject(&srcTokens, &tgtTokens)
		objDeleted := halfDiffObject(&tgtTokens, &srcTokens)
		added = append(added, objAdded...)
		deleted = append(deleted, objDeleted...)
	case ARRAY:
		arrAdded := halfDiffArray(&srcTokens, &tgtTokens)
		arrDeleted := halfDiffArray(&tgtTokens, &srcTokens)

		added = append(added, arrAdded...)
		deleted = append(deleted, arrDeleted...)
	default:
		isEqual, err := srcTokens.Equal(&tgtTokens)
		if err != nil {
			panic(err)
		}
		if !isEqual {
			added = append(added, source...)
			deleted = append(deleted, target...)
		}
	}

	return added, deleted
}

func getTokens(source []byte) jsonEntry {
	first := source[0]
	var jsonType int
	switch first {
	case '{':
		return getObject(source)
	case '[':
		return getArray(source)
	default:
		value := getToken(source)
		return jsonEntry{jsonType, len(value), nil, nil, nil, &value}
	}
}

func getArray(source []byte) jsonEntry {
	tokens := make([]jsonEntry, 0)
	last := 0
	for idx, b := range source {
		if b == ']' {
			last = idx
			break
		}
	}
	for idx := 1; idx < last; idx++ {
		char := source[idx]
		if char == ',' {
			continue
		}
		tokenData := getTokens(source[idx:])
		tokens = append(tokens, tokenData)
		idx += tokenData.length
	}
	return jsonEntry{ARRAY, last, nil, nil, &tokens, nil}
}

func getObject(source []byte) jsonEntry {
	tokens := make([]jsonEntry, 0)
	keys := make([]jsonEntry, 0)
	values := make([]jsonEntry, 0)

	last := 0
	for idx, b := range source {
		if b == '}' {
			last = idx
			break
		}
	}
	for idx := 1; idx < last; idx++ {
		char := source[idx]
		if char == ',' {
			continue
		}
		tokenData := getTokens(source[idx:])
		tokens = append(tokens, tokenData)
		idx += tokenData.length
	}
	for idx, tkn := range tokens {
		if idx%2 == 0 {
			keys = append(keys, tkn)
		} else {
			values = append(values, tkn)
		}
	}
	return jsonEntry{OBJECT, last, &keys, &values, nil, nil}
}

func getToken(source []byte) []byte {
	for idx, b := range source {
		if slices.Contains(TokenEnders[:], b) {
			return source[:idx]
		}
	}
	return source
}
