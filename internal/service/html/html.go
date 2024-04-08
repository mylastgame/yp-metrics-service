package html

import "fmt"

func SliceToOlLi(title string, s []fmt.Stringer) string {
	var html string
	if title != "" {
		html = fmt.Sprintf("%s: <ol>", title)
	} else {
		html = "<ol>"
	}

	for _, str := range s {
		html += fmt.Sprintf("<li>%s</li>", str)
	}
	html += "</ol>"

	return html
}

func Tag(tag string, val fmt.Stringer) string {
	return fmt.Sprintf("<%s>%s</%s>", tag, val, tag)
}
