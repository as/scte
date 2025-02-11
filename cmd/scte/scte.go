package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/as/scte"
)

func main() {
	data, err := io.ReadAll(os.Stdin)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error reading from stdin: %v", err)
		os.Exit(1)
	}

	p, err := scte.Parse(data)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error parsing packet: %v", err)
		os.Exit(1)
	}

	printjson(p)

}

func printjson(v any) {
	p, _ := json.Marshal(v)
	fmt.Println(string(p))
}
