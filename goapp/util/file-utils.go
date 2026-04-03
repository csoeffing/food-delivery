package util

import "strings"

func GetFilenameAndExtension(fullFilename string) (string, string) {
	fname := ""
	ext := ""

	p := strings.LastIndex(fullFilename, ".")

	if p == -1 {
		fname = fullFilename
	} else {
		fname = fullFilename[:p]
		ext = fullFilename[p+1:]
	}

	return fname, ext
}

func GetLastPathComponent(path string) string {
	p := strings.LastIndex(path, "/")

	if p == -1 {
		return path
	} else {
		return path[p+1:]
	}
}
