package remove

import "strings"

func FileExtension(filename string, extension string) string {
	parts := strings.Split(filename, ".")
	if len(parts) == 0 {
		return filename
	}
	if !strings.EqualFold(parts[len(parts)-1], extension) {
		return filename
	}
	return strings.Join(parts[0:len(parts)-1], ".")
}
