// Wraps some essential Operating System/Input-Output (OS/IO) path features with helper functions
// to make them more digestible and clarify their API.
//
// None of those functionalities rely on existing files or folders, they simply operate on a string.
package path

import (
	"path/filepath"
	"strings"
)

func New(elements ...string) string {
	return normalize(filepath.Join(elements...))
}

func IsDirectory(path string) bool {
	path = normalize(path)
	if strings.HasSuffix(path, "/") {
		return true
	}
	return Extension(path) == ""
}
func IsFile(path string) bool {
	return Extension(path) != ""
}

func LastPart(path string) string {
	return filepath.Base(path)
}
func Folder(path string) string {
	return normalize(filepath.Dir(path))
}
func Extension(path string) string {
	return normalize(filepath.Ext(path))
}
func RemoveExtension(path string) string {
	var ext = Extension(path)
	if ext == "" {
		return path
	}
	return normalize(path[:len(path)-len(ext)])
}

// private ========================================================

func normalize(path string) string {
	return strings.ReplaceAll(path, "\\", "/")
}
