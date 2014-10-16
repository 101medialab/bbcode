// Copyright 2014 Frustra. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package bbcode

import (
	"html"
	"net/url"
	"sort"
	"strings"
)

// HTMLTag represents a DOM node.
type HTMLTag struct {
	Name     string
	Value    string
	Attrs    map[string]string
	Children []*HTMLTag
}

// NewHTMLTag creates a new HTMLTag with string contents specified by value.
func NewHTMLTag(value string) *HTMLTag {
	return &HTMLTag{
		Value:    value,
		Attrs:    make(map[string]string),
		Children: make([]*HTMLTag, 0),
	}
}

func (t *HTMLTag) String() string {
	var value string
	if t.Value != "" {
		value = html.EscapeString(t.Value)
	}
	var attrString string

	keys := make([]string, len(t.Attrs))
	i := 0
	for key := range t.Attrs {
		keys[i] = key
		i++
	}
	sort.Strings(keys)
	for _, key := range keys {
		attrString += " " + key + `="` + strings.Replace(html.EscapeString(t.Attrs[key]), "\n", "", -1) + `"`
	}
	if len(t.Children) > 0 {
		var childrenString string
		for _, child := range t.Children {
			childrenString += child.String()
		}
		if t.Name != "" {
			return value + "<" + t.Name + attrString + ">" + childrenString + "</" + t.Name + ">"
		} else {
			return value + childrenString
		}
	} else if t.Name != "" {
		return value + "<" + t.Name + attrString + ">"
	} else {
		return value
	}
}

func (t *HTMLTag) AppendChild(child *HTMLTag) *HTMLTag {
	if child == nil {
		t.Children = append(t.Children, NewHTMLTag(""))
	} else {
		t.Children = append(t.Children, child)
	}
	return t
}

func InsertNewlines(out *HTMLTag) {
	if strings.ContainsRune(out.Value, '\n') {
		parts := strings.Split(out.Value, "\n")
		for i, part := range parts {
			if i == 0 {
				out.Value = parts[i]
			} else {
				out.AppendChild(NewlineTag())
				if len(part) > 0 {
					out.AppendChild(NewHTMLTag(part))
				}
			}
		}
	}
}

// Returns a new HTMLTag representing a line break
func NewlineTag() *HTMLTag {
	var out = NewHTMLTag("")
	out.Name = "br"
	return out
}

func ValidURL(raw string) string {
	u, err := url.Parse(raw)
	if err != nil {
		return ""
	}
	return u.String()
}
