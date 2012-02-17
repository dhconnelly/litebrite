package main

import (
	"go/scanner"
	"go/token"
	"fmt"
	"os"
	"io/ioutil"
	"html/template"
)

// get token types for each position in src
func tokenize(filename string, src []byte) []token.Token {
	tokens := make([]token.Token, 0)
	var s scanner.Scanner
	fset := token.NewFileSet()
	file := fset.AddFile(filename, fset.Base(), len(src))
	s.Init(file, src, nil, 0)
	for {
		_, tok, _ := s.Scan()
		if tok == token.EOF {
			break
		}
		tokens = append(tokens, tok)
	}
	return tokens
}

// breaks src into substrings, returning a map from position to string
func split(filename string, src []byte) []string {
	splits := make([]string, 0)
	positions := make([]int, 0)

	// find the starting positions of all tokens
	var s scanner.Scanner
	fset := token.NewFileSet()
	file := fset.AddFile(filename, fset.Base(), len(src))
	s.Init(file, src, nil, 0)
	for {
		pos, tok, _ := s.Scan()
		if tok == token.EOF {
			break
		}
		positions = append(positions, int(pos) - 1) // TODO wtf -1
	}
	positions = append(positions, positions[len(positions)-1] + 1)

	// split the source at each position to get a slice of substrings
	for i := 1; i < len(positions); i++ {
		start, end := positions[i-1], positions[i]
		splits = append(splits, string(src[start:end]))
	}
	return splits
}

// wrap each substring with corresponding token tag and join
func tag(src []string, tokens []token.Token, t *template.Template) string {
	fmt.Print("<pre>")
	for i, line := range src {
		tok := tokens[i]
		s := struct{Tag, Code string}{tok.String(), line}
		t.Execute(os.Stdout, s)
	}
	fmt.Print("</pre>")
	return ""
}

const TAG = `<div style="display: inline" class="{{.Tag}}">{{.Code}}</div>`
var t = template.Must(template.New("golang-code").Parse(TAG))

func main() {
	files := os.Args[1:]
	for _, filename := range files {
		src, err := ioutil.ReadFile(filename)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err.Error())
			return
		}
		tokens := tokenize(filename, src)
		splits := split(filename, src)
		fmt.Println(tag(splits, tokens, t))
	}
}
