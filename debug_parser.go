package main

import (
	"fmt"

	"github.com/goccy/go-yaml/lexer"
	"github.com/goccy/go-yaml/token"
)

func main() {
	yaml := "- &anchorName value\n- *anchorName"
	tokens := lexer.Tokenize(yaml)

	for i, t := range tokens {
		fmt.Printf("[%d] Type: %s (%d) | Value: %q | Origin: %q\n",
			i, t.Type, t.Type, t.Value, t.Origin)

		if t.Type == token.AnchorType || t.Type == token.AliasType {
			fmt.Println("   -> Is Anchor/Alias")
		}
	}
}
