package mimechecker

import (
	"testing"
	"github.com/duffpl/go-finder/mimechecker/mock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"errors"
)

func TestMulti_ByPath(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mock1 := mock_mimechecker.NewMockChecker(ctrl)
	mock2 := mock_mimechecker.NewMockChecker(ctrl)
	multi := NewMulti(mock1, mock2)
	testPath := "test.path"
	t.Run("ReturnFirstResult", func(t *testing.T) {
		mock1.EXPECT().TypeByFile(testPath).Return("mock-1/mime", nil)
		mock2.EXPECT().TypeByFile(gomock.Any()).Times(0)
		result, _ := multi.TypeByFile(testPath)
		assert.Equal(t, "mock-1/mime", result)
	})
	t.Run("CheckUntilMimeHasBeenFound", func(t *testing.T) {
		mock1.EXPECT().TypeByFile(testPath).Return("", nil)
		mock2.EXPECT().TypeByFile(testPath).Return("mock-2/mime", nil)
		result, _ := multi.TypeByFile(testPath)
		assert.Equal(t, "mock-2/mime", result)
	})
	t.Run("ReturnWrappedErrAsSoonAsErrorIsEncountered", func(t *testing.T) {
		mock1.EXPECT().TypeByFile(testPath).Return("", errors.New("mock-1-error"))
		mock2.EXPECT().TypeByFile(gomock.Any()).Times(0)
		result, err := multi.TypeByFile(testPath)
		assert.Equal(t, "", result)
		assert.Error(t, err)
		assert.Equal(t, "multi mimechecker: mock-1-error", err.Error())
	})
}
