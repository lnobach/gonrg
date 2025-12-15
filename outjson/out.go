package outjson

import (
	"encoding/json"
	"fmt"

	"github.com/lnobach/gonrg/obis"
)

func PrintJSON(res *obis.OBISListResult) {
	bytes, err := json.MarshalIndent(res, "", "  ")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(bytes))
}
