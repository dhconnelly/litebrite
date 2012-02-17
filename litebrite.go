package litebrite

import (
	"bytes"
	"go/scanner"
	"go/token"
	"html/template"
)

type codeSegment struct {
	Code  string // a segment of source code with the same style
	Pos   int    // the position of the segment in the source
	Tok   token.Token // the token type of the segment
	Class string
}

// getSegments splits the source into same-token-type chunks.
func getSegments(src []byte) []*codeSegment {
	segments := make([]*codeSegment, 0)

	// find the starting positions of all segments
	var s scanner.Scanner
	fset := token.NewFileSet()
	file := fset.AddFile("", fset.Base(), len(src))
	s.Init(file, src, nil, 1)
	for {
		pos, tok, _ := s.Scan()
		if tok == token.EOF {
			break
		}
		segment := &codeSegment{Pos: int(pos)-1, Tok: tok} // WTF -1
		segments = append(segments, segment)
	}

	// split the source at each position to get segments
	for i, segment := range segments {
		if i+1 == len(segments) {
			segment.Code = string(src[segment.Pos:])
		} else {
			next := segments[i+1]
			segment.Code = string(src[segment.Pos:next.Pos])
		}
	}

	return segments
}

const ELEM = `{{with .Class}}<div class="{{.}}">{{end}}{{.Code}}{{with .Class}}</div>{{end}}`

var elemT = template.Must(template.New("golang-elem").Parse(ELEM))

// buildHTML constructs an HTML string of elements from the segments in src.
func buildHTML(src []*codeSegment) string {
	var b bytes.Buffer
	for _, segment := range src {
		elemT.Execute(&b, segment)
	}

	return string(b.Bytes())
}

// Highlighter contains the CSS class names that are applied to the
// corresponding source code token types.
type Highlighter struct {
	OperatorClass string
	IdentClass string
	LiteralClass string
	KeywordClass string
	CommentClass string
}

// highlightSegment adds the CSS class name specified by h to the segment.
func (h *Highlighter) highlightSegment(s *codeSegment) {
	switch {
	case s.Tok.IsKeyword():
		s.Class = h.KeywordClass
	case s.Tok.IsLiteral():
		if s.Tok == token.IDENT {
			s.Class = h.IdentClass
		} else {
			s.Class = h.LiteralClass
		}
	case s.Tok.IsOperator():
		if s.Tok == token.SEMICOLON && s.Code != ";" {
			return
		}
		s.Class = h.OperatorClass
	case s.Tok == token.COMMENT:
		s.Class = h.CommentClass
	default:
		panic("unknown token type!")
	}
}

// Highlight returns an HTML fragment containing elements for all Go tokens
// in src.  The elements will be of the form <div class="TYPE_CLASS">CODE</div>
// where TYPE_CLASS is the CSS class name provided in h corresponding to the
// token type of CODE.  For instance, if CODE is a keyword, then TYPE_CLASS
// will be h.keywordClass.
func (h *Highlighter) Highlight(src string) string {
	data := []byte(src)
	segments := getSegments(data)
	for _, segment := range segments {
		h.highlightSegment(segment)
	}
	return buildHTML(segments)
}
