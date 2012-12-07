---
title: litebrite
layout: default
---

litebrite

## ABOUT

litebrite is a library for generating syntax-highlighted HTML from your Go
source code.  For an example of its use, check out its annotated source code:
dhconnelly.github.com/litebrite.

## INSTALL

go get github.com/dhconnelly/litebrite
go install github.com/dhconnelly/litebrite

## USAGE

Make sure you import "github.com/dhconnelly/litebrite" in your source file.

Then you can

    h := new(litebrite.Highlighter)
    h.CommentClass = "commentz"
    h.OperatorClass = "opz"
    // add some more classes names, see below
    html := h.Highlight(myCodez)

This will return a string of HTML where every comment in the string myCodez
is wrapped with a `<div class="commentz">` tag and every operator is wrapped
with a `<div class="opz">` tag.

The following string fields are available on a Highlighter struct:

CommentClass
OperatorClass
IdentClass
LiteralClass
KeywordClass

Setting a field to a non-nil value causes tokens of that type to be wrapped
with a `<div>` that has your specified CSS class name.

## AUTHOR

docgo was written by Daniel Connelly.  You can find my stuff at dhconnelly.com.

## LICENSE

litebrite is released under a BSD-style license available in LICENSE.md.
