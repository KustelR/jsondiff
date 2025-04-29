package jsondiff_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/KustelR/jsondiff"
)

type CardJson struct {
	Id          string   `json:"id"`
	ColumnId    string   `json:"columnId"`
	Name        string   `json:"name"`
	Order       int      `json:"order"`
	Description string   `json:"description"`
	TagIds      []string `json:"tagIds"`
	CreatedAt   int      `json:"createdAt"`
	UpdatedAt   int      `json:"updatedAt"`
	CreatedBy   string   `json:"createdBy"`
	UpdatedBy   string   `json:"updatedBy"`
}

func TestHelloWorld(t *testing.T) {

	MockCard1 := CardJson{
		Id:          "same",
		ColumnId:    "nt same",
		Name:        "asdf sadfas",
		Order:       123,
		CreatedAt:   12321,
		TagIds:      make([]string, 0),
		CreatedBy:   "",
		Description: "dsfasd asdfs adfasdf asdf",
	}
	MockCard2 := CardJson{
		Id:          "same",
		ColumnId:    "not same",
		Name:        "asdf sadfas",
		Order:       12234,
		TagIds:      make([]string, 0),
		CreatedAt:   12321,
		CreatedBy:   "",
		Description: "dsfasd asdfs adfasdf asdf",
	}
	json1, _ := json.Marshal(MockCard1)
	json2, _ := json.Marshal(MockCard2)
	diff1, diff2 := jsondiff.Diff(json1, json2)
	var recreatedCard1 CardJson
	var recreatedCard2 CardJson
	err := json.Unmarshal(diff1, &recreatedCard1)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(diff2, &recreatedCard2)
	if err != nil {
		panic(err)
	}
	if recreatedCard1.Name != "" {
		t.Errorf("created property which shouldn't")
	}
	fmt.Println(string(diff1), string(diff2))
}
