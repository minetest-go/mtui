package cdb

import (
	"fmt"
	"strings"
)

func GetAuthorName(pkgname string) (string, string) {
	parts := strings.Split(pkgname, "/")
	return parts[0], parts[1]
}

func GetPackagename(author, name string) string {
	return fmt.Sprintf("%s/%s", author, name)
}
