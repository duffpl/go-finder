package finder

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"sort"
	"github.com/duffpl/go-finder/os"
	"time"
	os2 "os"
)

func TestFinder_Size_integration(t *testing.T) {
	testExpectations := []struct {
		op  CmpOperator
		val int
		res []string
	}{
		{Equal, 100, []string{"size-100.dat"}},
		{MoreThan, 50, []string{"size-100.dat", "size-150.dat"}},
		{MoreOrEqual, 100, []string{"size-100.dat", "size-150.dat"}},
		{LessThan, 100, []string{"size-50.dat"}},
		{LessOrEqual, 100, []string{"size-100.dat", "size-50.dat"}},
	}
	for _, tex := range testExpectations {
		t.Run(string(tex.op), func(t *testing.T) {
			result, _ := New().
				Size(tex.op, int64(tex.val)).
				Glob("./**/size-*.dat")
			assert.Equal(t, tex.res, getFileNamesFromResult(result))
		})
	}
	t.Run("ErrorOnInvalidOperator", func(t *testing.T) {
		_, err := New().Size(":(", 0).Glob(".")
		assert.Error(t, err)
		assert.Equal(t, Errors.InvalidSizeOperator, err)
	})
}

func TestFinder_Size_unit(t *testing.T) {
	testExpectations := []struct {
		operator CmpOperator
		value    int
		result   []string
	}{
		{Equal, 100, []string{"100"}},
		{MoreThan, 50, []string{"100", "150"}},
		{MoreOrEqual, 100, []string{"100", "150"}},
		{LessThan, 100, []string{"50"}},
		{LessOrEqual, 100, []string{"100", "50"}},
	}
	mockData := []os.FileInfoEx{
		&mockFileInfoEx{name: "100", size: 100,},
		&mockFileInfoEx{name: "150", size: 150,},
		&mockFileInfoEx{name: "50", size: 50,},
	}

	for _, expectation := range testExpectations {
		t.Run(string(expectation.operator), func(t *testing.T) {
			sut := New().SetGlobFunc(func(pattern string) ([]os.FileInfoEx, error) {
				return mockData, nil
			}).Size(expectation.operator, int64(expectation.value))

			result, _ := sut.Glob("?")
			assert.Equal(t, expectation.result, getFileNamesFromResult(result))
		})
	}
}

func TestFinder_Mime(t *testing.T) {
	testExpectations := []struct {
		mime   string
		result []string
	}{
		{"image/jpeg", []string{"image.jpg"}},
		{"image/png", []string{"image.png"}},
		{"application/pdf", []string{"pdf.pdf"}},
		{"audio/mpeg", []string{"audio-no-extension", "audio.mp3"}},
		{"image/gif", []string{"image-gif-fake-audio.mp3"}},
	}
	for _, expectation := range testExpectations {
		result, err := New().
			Mime(expectation.mime).
			Glob("./test_files/mime/*")
		if err != nil {
			t.Error(err)
		}
		assert.Equal(t, expectation.result, getFileNamesFromResult(result))
	}
}

func TestFinder_Mime_unit(t *testing.T) {
	testExpectations := []struct {
		mime string
		res  []string
	}{
		{"application/mock-data-a", []string{"mock-data-a-1"}},
		{"application/mock-data-b", []string{"mock-data-b-1", "mock-data-b-2"}},
	}
	for _, tex := range testExpectations {
		testGlobFunc := func(pattern string) (result []os.FileInfoEx, err error) {
			return []os.FileInfoEx{
				&mockFileInfoEx{name: "mock-data-a-1", mime: "application/mock-data-a"},
				&mockFileInfoEx{name: "mock-data-b-1", mime: "application/mock-data-b"},
				&mockFileInfoEx{name: "mock-data-b-2", mime: "application/mock-data-b"},
			}, nil
		}
		result, err := New().
			SetGlobFunc(testGlobFunc).
			Mime(tex.mime).
			Glob("test-glob")
		if err != nil {
			t.Error(err)
		}
		assert.Equal(t, tex.res, getFileNamesFromResult(result))
	}
}

func TestFinder_MimeRegexp(t *testing.T) {
	testExpectations := []struct {
		pattern string
		result     []string
	}{
		{"image.*?p", []string{"image.jpg", "image.png"}},
		{"^audio", []string{"audio-no-extension", "audio.mp3"}},
		{"pdf$", []string{"pdf.pdf"}},
	}
	for _, expectation := range testExpectations {
		result, err := New().
			MimeRegexp(expectation.pattern).
			Glob("./test_files/mime/*")
		if err != nil {
			t.Error(err)
		}
		assert.Equal(t, expectation.result, getFileNamesFromResult(result))
	}
}

func getFileNamesFromResult(result []os.FileInfoEx) (names []string) {
	for _, item := range result {
		names = append(names, item.Name())
	}
	sort.Slice(names, func(i, j int) bool {
		return names[i] <= names[j]
	})
	return names
}

