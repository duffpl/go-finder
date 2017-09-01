package finder

import (
	"log"
	"io/ioutil"
	"path/filepath"
	"bytes"
)

type FilterFunction func(entry *FileEntry) bool

type Finder struct {
	filters []FilterFunction
}

func New() *Finder {
	return &Finder{}
}

func (f *Finder) Find(rootPath string) (result []*FileEntry, err error){
	absoluteRootPath, err := filepath.Abs(rootPath)
	if err != nil {
		return
	}
	result = f.getListing(absoluteRootPath, absoluteRootPath)

	return
}


func (f *Finder) checkFilters(input *FileEntry) bool {
	for _, filterFunction := range f.filters {
		if !filterFunction(input) {
			return false
		}
	}
	return true
}

func (f *Finder) addFilter(ff FilterFunction) {
	f.filters = append(f.filters, ff)
}

func (f *Finder) getListing(dirPath string, sourcePath string) (result []*FileEntry) {
	files, err := ioutil.ReadDir(dirPath)
	if err != nil {
		log.Fatal(err)
	}
	for _,file := range files {
		absolutePath := dirPath + "/" +file.Name()
		relativePath := absolutePath[len(sourcePath)+1:]
		entry := newFileEntryFromFileInfo(file)
		entry.AbsolutePath = absolutePath
		entry.RelativePath = relativePath
		if !f.checkFilters(entry) {
			continue
		}
		result = append(result, entry)
		if entry.IsDir {
			result = append(result, f.getListing(absolutePath, sourcePath)...)
		}
	}
	return
}

type FileCollection []*FileEntry

type CompareResult struct {
	MissingInSource DifferenceResult
	MissingInDestination DifferenceResult
}

type DifferenceResult struct {
	Missing FileCollection
	Different FileCollection
}
var TotalCompares int = 0
func CompareCollections(src FileCollection, dst FileCollection) CompareResult {
	return CompareResult{
		MissingInSource: getDifference(dst, src),
		MissingInDestination: getDifference(src, dst),
	}
}


func getDifference(src FileCollection, dst FileCollection) (result DifferenceResult) {
	for _, srcFile := range src {
		switch searchResult := isEntryInCollection(srcFile, dst); searchResult {
		case ItemMissing:
			result.Missing = append(result.Missing, srcFile)
		case ItemDiffers:
			result.Different = append(result.Different, srcFile)
		}
	}
	return
}

func isEntryInCollection(entry *FileEntry, collection FileCollection) CollectionSearchResult {
	for _, entryInCollection := range collection {
		TotalCompares++
		if entryInCollection.RelativePath == entry.RelativePath {
			if !bytes.Equal(entryInCollection.Checksum, entry.Checksum) {
				return ItemDiffers
			}
			return ItemFound
		}
	}
	return ItemMissing
}

type CollectionSearchResult int8

const (
	ItemMissing CollectionSearchResult = iota
	ItemDiffers
	ItemFound
)