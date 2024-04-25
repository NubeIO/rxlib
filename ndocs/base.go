package ndocs

import (
	"encoding/json"
	"github.com/NubeIO/rxlib/libs/docgen"
	"github.com/texttheater/golang-levenshtein/levenshtein"
	"strings"
)

type Docs struct {
	docs []*docgen.Helper
}

func New(name string) *Docs {
	docs := &Docs{}
	docs.LoadFromString(name)
	return docs
}

func (d *Docs) LoadFromString(jsonString string) error {
	return json.Unmarshal([]byte(jsonString), &d.docs)
}

func (d *Docs) All() []*docgen.Helper {
	return d.docs
}

func (d *Docs) Find(name string) *docgen.Helper {
	for _, method := range d.docs {
		if method.Name == name {
			return method
		}
	}
	return nil
}

func cleanSearchTerm(in string) string {
	keyWordsToRemove := []string{"runtime"}

	// If "in" has ".", let's do a string split and then remove keywords
	if strings.Contains(in, ".") {
		parts := strings.Split(in, ".")
		for i := range parts {
			for _, kw := range keyWordsToRemove {
				if parts[i] == kw {
					parts[i] = ""
				}
			}
		}
		in = strings.Join(parts, ".")
	} else {
		for _, kw := range keyWordsToRemove {
			in = strings.ReplaceAll(in, kw, "")
		}
	}
	return removeFirstDot(in)
}

func removeFirstDot(s string) string {
	if len(s) > 0 && s[0] == '.' {
		return s[1:]
	}
	return s
}

func (d *Docs) Fuzzy(name string) []string {
	var matches []string
	for _, method := range d.docs {
		if levenshtein.DistanceForStrings([]rune(strings.ToLower(method.Name)), []rune(strings.ToLower(name)), levenshtein.DefaultOptions) <= 5 {
			matches = append(matches, method.Name)
		}
	}
	return matches
}
