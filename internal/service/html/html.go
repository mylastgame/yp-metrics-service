package html

import "fmt"

func Tag(tag string, val string) string {
	return fmt.Sprintf("<%s>%s</%s>", tag, val, tag)
}
