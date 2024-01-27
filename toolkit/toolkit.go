package toolkit

import "strings"

func SliceToString(data []string) string {
	return "{" + strings.Join(data, ",") + "}"
}
