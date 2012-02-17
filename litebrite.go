package litebrite

import (
	"bytes"
	"fmt"
	"go/scanner"
	"go/token"
	"html/template"
	"sort"
)

type codeSegment struct {
	Code  string // a segment of source code with the same style
	Pos   int    // the position of the segment in the source
	Class string // the CSS class for the segment
}

// token types
const (
	KEYWORD  = "keyword"
	IDENT    = "ident"
	LITERAL  = "literal"
	OPERATOR = "operator"
)

// getClass returns the CSS class name for the specified token type.
func getClass(tok token.Token) (class string) {
	switch {
	case tok.IsKeyword():
		return KEYWORD
	case tok.IsLiteral():
		if tok == token.IDENT {
			return IDENT
		} else {
			return LITERAL
		}
	case tok.IsOperator():
		return OPERATOR
	default:
		panic("unknown token type!")
	}
	return
}

// getClasses returns a map from source token positions to CSS class names.
func getClasses(src []byte) map[int]string {
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
		tokens[int(pos)-1] = getClass(tok) // WTF -1
	}

	return tokens
}

// getSegments splits the source into same-token-type chunks.
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
		positions = append(positions, int(pos)-1) // WTF -1
	}

	// to help with slicing, add one more position at the end
	positions = append(positions, positions[len(positions)-1]+1)

	// split the source at each position to get segments
	for i := 1; i < len(positions); i++ {
		start, end := positions[i-1], positions[i]
		segment := string(src[start:end])
		segments[start] = &codeSegment{Code: segment, Pos: start}
	}

	return segments
}

// styleSegments adds the CSS class names in classes to the segments in src.
// The class name classes[i] is applied to the segment src[i].
func styleSegments(src map[int]*codeSegment, classes map[int]string) {
	for pos, class := range classes {
		src[pos].Class = class
	}
}

const CODE = "<pre><code class=\"golang\">%s</code></pre>"
const ELEM = `<div class="{{.Class}}">{{.Code}}</div>`

var elemT = template.Must(template.New("golang-elem").Parse(ELEM))

// buildHTML constructs an HTML string of elements from the segments in src.
func buildHTML(src map[int]*codeSegment) string {
	indices := make([]int, 0)
	for pos := range src {
		indices = append(indices, pos)
	}
	sort.Ints(indices)

	var b bytes.Buffer
	for _, pos := range indices {
		elemT.Execute(&b, src[pos])
	}

	return fmt.Sprintf(CODE, string(b.Bytes()))
}

// Highlight returns an HTML fragment containing elements for all Go tokens
// in src.  The elements will be of the form <div class=TYPE>CODE</div>, where
// TYPE is the token type ("keyword", "operator", "literal", or "ident")
// and CODE is the source fragment.  The entire fragment is wrapped with <pre>
// and <code class="golang"> tags.
func Highlight(src string) string {
	data := []byte(src)
	segs := getSegments(data)
	classes := getClasses(data)
	styleSegments(segs, classes)
	return buildHTML(segs)
}
