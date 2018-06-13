package pathutil

import (
	"path/filepath"
	"runtime"
	"strings"

	"github.com/pkg/errors"
)

const (
	projectPrefix = "kavirajk/bookshop" // should be in sync with actual project prefix
)

// ProjectRoot returns the project root's directory
// from wherever its called. Can be useful to read project-leve
// files like config and migrations in tests.
func ProjectRoot() (string, error) {
	_, path, _, ok := runtime.Caller(0)
	if !ok {
		return "", errors.New("project-root.unknown-error")
	}
	last := strings.LastIndex(path, projectPrefix)
	if last < 0 {
		return "", errors.New("wrong project path")
	}
	return filepath.Clean(filepath.Join(path[0:last], projectPrefix)), nil
}
