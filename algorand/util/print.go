package util

import (
	"encoding/json"
	"fmt"

	prettyjson "github.com/hokaccha/go-prettyjson"
)

func PrintJSON(o interface{}) {
	s, _ := prettyjson.Marshal(o)
	fmt.Println(string(s))
}

func PrettyPrintJSON(i interface{}) {
	s, _ := json.MarshalIndent(i, "", "\t")
	fmt.Printf("%+v", string(s))

}
