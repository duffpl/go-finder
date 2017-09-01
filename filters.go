package finder

import (
	"regexp"
	"log"
	"github.com/pkg/errors"
)

func (f *Finder) Exclude(excludes []string) *Finder {
	f.addFilter(newExcludeFilter(excludes))
	return f
}

func newExcludeFilter(excludes []string) FilterFunction{
	patterns := make([]*regexp.Regexp, len(excludes))
	for i, excludePattern := range excludes {
		compiled, err := regexp.Compile(excludePattern)
		if err != nil {
			log.Fatal(errors.Wrap(err, "Cannot compile pattern"))
		}
		patterns[i] = compiled
	}
	return func(input *FileEntry) bool {
		for _, pattern := range patterns {
			res := pattern.FindStringIndex(input.RelativePath)
			if res != nil {
				return false
			}
		}
		return true
	}
}
