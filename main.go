package main

import (
	"encoding/json"
	"fmt"
)

func main() {
	// init things
	parse_flags()
	logger_options()
	// cmus checks
	cmus_output := cmusQuery()
	cmus_status := cmusParse(cmus_output)
	// json output, ignoring error
	b, _ := json.Marshal(cmus_status)
	fmt.Println(string(b))
}
