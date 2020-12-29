package gr8http

import (
	"io/ioutil"
	"net/http"
	"strings"
)

// NewDir returns an initialized *Dir
func NewDir(dir http.Dir) *Dir {
	return &Dir{
		dir:   dir,
		files: map[string]func() (http.File, error){},
	}
}

// Dir is an http.Filesystem that extends Dir and provides
// the ability to server arbitrary files
type Dir struct {
	dir   http.Dir
	files map[string]func() (http.File, error)
}

// AddFile adds a file at name provided by fn()
func (d *Dir) AddFile(name string, fn func() (http.File, error)) {
	if !strings.HasPrefix(name, "/") {
		name = "/" + name
	}

	d.files[name] = fn
}

// Open satisfies http.Filesystem
func (d *Dir) Open(name string) (http.File, error) {
	if fn, ok := d.files[name]; ok {
		return fn()
	}

	return d.dir.Open(name)
}

// ReadFile reads file contents from the file name in fs
func ReadFile(fs http.FileSystem, name string) ([]byte, error) {
	f, err := fs.Open(name)
	if err != nil {
		return nil, err
	}

	defer f.Close()
	return ioutil.ReadAll(f)
}
