package finder

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"sort"
	"github.com/duffpl/go-finder/file"
)

func TestFinder_Size(t *testing.T) {
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
	mockGlob := newMockGlobFunc([]file.FileInfoEx{
		&mockFileInfoEx{name: "100", size: 100,},
		&mockFileInfoEx{name: "150", size: 150,},
		&mockFileInfoEx{name: "50", size: 50,},
	})

	for _, expectation := range testExpectations {
		t.Run(string(expectation.operator), func(t *testing.T) {
			sut := New().
				SetGlobFunc(mockGlob).
				Size(expectation.operator, int64(expectation.value))
			result, _ := sut.Glob("*")
			assert.Equal(t, expectation.result, getFileNamesFromResult(result))
		})
	}
}

func TestFinder_Mime_integration(t *testing.T) {
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
	mockGlob := newMockGlobFunc([]file.FileInfoEx{
		&mockFileInfoEx{name: "mock-data-a-1", mime: "application/mock-data-a"},
		&mockFileInfoEx{name: "mock-data-b-1", mime: "application/mock-data-b"},
		&mockFileInfoEx{name: "mock-data-b-2", mime: "application/mock-data-b"},
	})
	for _, tex := range testExpectations {
		result, err := New().
			SetGlobFunc(mockGlob).
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
		result  []string
	}{
		{"^r.*?p", []string{"ramp", "rump"}},
		{"audio$", []string{"some-audio"}},
	}
	mockGlob := newMockGlobFunc([]file.FileInfoEx{
		&mockFileInfoEx{name: "ramp", mime: "raaaamp/x"},
		&mockFileInfoEx{name: "gamp", mime: "gaaaamp/x"},
		&mockFileInfoEx{name: "rump", mime: "raaummp/x"},
		&mockFileInfoEx{name: "some-audio", mime: "mpeg/audio"},
		&mockFileInfoEx{name: "audio-on-start", mime: "audio/mpeg"},
	})
	for _, expectation := range testExpectations {
		result, err := New().
			SetGlobFunc(mockGlob).
			MimeRegexp(expectation.pattern).
			Glob("*")
		if err != nil {
			t.Error(err)
		}
		assert.Equal(t, expectation.result, getFileNamesFromResult(result))
	}
}

func TestFinder_Checksum_unit(t *testing.T) {
	testExpectations := []struct {
		checksum string
		result   []string
	}{
		{"01", []string{"1"}},
		{"03", []string{"3"}},
	}
	mockGlob := newMockGlobFunc([]file.FileInfoEx{
		&mockFileInfoEx{name: "1", checksum: []byte{1}},
		&mockFileInfoEx{name: "2", checksum: []byte{2}},
		&mockFileInfoEx{name: "3", checksum: []byte{3}},
		&mockFileInfoEx{name: "4", checksum: []byte{4}},
	})
	for _, expectation := range testExpectations {
		result, err := New().
			SetGlobFunc(mockGlob).
			Checksum(expectation.checksum).
			Glob("*")
		if err != nil {
			t.Error(err)
		}
		assert.Equal(t, expectation.result, getFileNamesFromResult(result))
	}
}

func TestFinder_Checksum_integration(t *testing.T) {
	testExpectations := []struct {
		checksum string
		result   []string
	}{
		{"60b725f10c9c85c70d97880dfe8191b3", []string{"60b725f10c9c85c70d97880dfe8191b3"}},
		{"3b5d5c3712955042212316173ccf37be", []string{"3b5d5c3712955042212316173ccf37be"}},
	}
	for _, expectation := range testExpectations {
		result, err := New().
			Checksum(expectation.checksum).
			Glob("./test_files/checksum/*")
		if err != nil {
			t.Error(err)
		}
		assert.Equal(t, expectation.result, getFileNamesFromResult(result))
	}
}

func TestFinder_RegexpName(t *testing.T) {
	testExpectations := []struct {
		pattern string
		res  []string
	}{
		{"match", []string{"also-match", "match"}},
		{"^match", []string{"match"}},
	}
	mockGlob := newMockGlobFunc([]file.FileInfoEx{
		&mockFileInfoEx{name: "match", abs:"abs/match"},
		&mockFileInfoEx{name: "not", abs:"abs/not"},
		&mockFileInfoEx{name: "also-match", abs: "abs/also-match"},
	})
	for _, expectation := range testExpectations {
		result, err := New().
			SetGlobFunc(mockGlob).
			RegexpName(expectation.pattern).
			Glob("test-glob")
		if err != nil {
			t.Error(err)
		}
		assert.Equal(t, expectation.res, getFileNamesFromResult(result))
	}
}

func TestFinder_RegexpPath(t *testing.T) {
	testExpectations := []struct {
		pattern string
		res  []string
	}{
		{"^abs/.*h", []string{"match"}},
		{"abs/oth", []string{"other-match"}},
	}
	mockGlob := newMockGlobFunc([]file.FileInfoEx{
		&mockFileInfoEx{name: "match", abs:"abs/match"},
		&mockFileInfoEx{name: "not", abs:"abs/not"},
		&mockFileInfoEx{name: "other-match", abs: "other-abs/other-match"},
	})
	for _, expectation := range testExpectations {
		result, err := New().
			SetGlobFunc(mockGlob).
			RegexpPath(expectation.pattern).
			Glob("test-glob")
		if err != nil {
			t.Error(err)
		}
		assert.Equal(t, expectation.res, getFileNamesFromResult(result), "a")
	}
}


func getFileNamesFromResult(result []file.FileInfoEx) (names []string) {
	for _, item := range result {
		names = append(names, item.Name())
	}
	sort.Slice(names, func(i, j int) bool {
		return names[i] <= names[j]
	})
	return names
}
