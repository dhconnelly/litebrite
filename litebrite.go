package main

import (
	"go/scanner"
	"go/token"
	"fmt"
	"os"
	"io/ioutil"
	"html/template"
	"bytes"
	"sort"
)

type codeSegment struct {
	Code string // a segment of source code with the same style
	Pos int // the position of the segment in the source
	Tag string // the CSS class for the segment
}

const (
	KEYWORD = "keyword"
	IDENT = "ident"
	LITERAL = "literal"
	OPERATOR = "operator"
)

func getTag(tok token.Token) (tag string) {
	switch {
	case tok.IsKeyword():
		tag = KEYWORD
	case tok.IsLiteral():
		if tok == token.IDENT {
			tag = IDENT
		} else {
			tag = LITERAL
		}
	case tok.IsOperator():
		tag = OPERATOR
	default:
		panic("unknown token type!")
	}
	return
}

// get token tags for each position in src
func getTags(src []byte) map[int]string {
	tokens := make(map[int]string)
	var s scanner.Scanner
	fset := token.NewFileSet()
	file := fset.AddFile("", fset.Base(), len(src))
	s.Init(file, src, nil, 0)
	for {
		pos, tok, _ := s.Scan()
		if tok == token.EOF {
			break
		}
		tokens[int(pos) - 1] = getTag(tok) // WTF -1
	}
	return tokens
}

// breaks src into segments; returns a map from segment position in src to segment
func getSegments(src []byte) map[int]*codeSegment {
	segments := make(map[int]*codeSegment)
	positions := make([]int, 0)

	// find the starting positions of all tokens
	var s scanner.Scanner
	fset := token.NewFileSet()
	file := fset.AddFile("", fset.Base(), len(src))
	s.Init(file, src, nil, 0)
	for {
		pos, tok, _ := s.Scan()
		if tok == token.EOF {
			break
		}
		positions = append(positions, int(pos) - 1) // WTF -1
	}
	positions = append(positions, positions[len(positions)-1] + 1)

	// split the source at each position to get a slice of substrings
	for i := 1; i < len(positions); i++ {
		start, end := positions[i-1], positions[i]
		segments[start] = &codeSegment{Code: string(src[start:end]), Pos: start}
	}
	
	return segments
}

func tagSegments(src map[int]*codeSegment, tags map[int]string) {
	for pos, tag := range tags {
		src[pos].Tag = tag
	}
}

const TAG = `<div style="display: inline" class="{{.Tag}}">{{.Code}}</div>`
var t = template.Must(template.New("golang-code").Parse(TAG))

func buildHTML(src map[int]*codeSegment) string {
	indices := make([]int, 0)
	for pos := range src {
		indices = append(indices, pos)
	}
	sort.Ints(indices)
	
	var b bytes.Buffer
	for _, pos := range indices {
		t.Execute(&b, src[pos])
	}
	return "<pre><code>" + string(b.Bytes())  + "</code></pre>"
}

func Highlight(src []byte) string {
	segs := getSegments(src)
	tagSegments(segs, getTags(src))
	return buildHTML(segs)
}

func main() {
	files := os.Args[1:]
	for _, filename := range files {
		src, err := ioutil.ReadFile(filename)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err.Error())
			return
		}
		fmt.Println(Highlight(src))
	}
}
