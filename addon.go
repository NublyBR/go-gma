package gma

import "io"

type metadata struct {
	Description string   `json:"description"`
	Type        string   `json:"type"`
	Tags        []string `json:"tags"`
}

type gma struct {
	stream io.ReadSeekCloser

	name   string
	author string
	meta   *metadata

	start int64

	files   []*entry
	pathMap map[string]*entry
}

func (a *gma) Name() string {
	return a.name
}

func (a *gma) Author() string {
	return a.author
}

func (a *gma) Description() string {
	return a.meta.Description
}

func (a *gma) Type() string {
	return a.meta.Type
}

func (a *gma) Tags() []string {
	return a.meta.Tags
}

func (a *gma) Entries() []Entry {
	ret := make([]Entry, len(a.files))
	for i, f := range a.files {
		ret[i] = f
	}
	return ret
}

func (a *gma) Open(path string) (io.ReadSeekCloser, error) {
	if e, ok := a.pathMap[path]; ok {
		return e.Open(), nil
	}

	return nil, ErrInvalidPath
}

func (a *gma) Close() error {
	return a.stream.Close()
}
