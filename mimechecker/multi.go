package mimechecker

import "errors"

type Multi struct {
	checkers []Checker
}

func (c Multi) TypeByFile(path string) (m string, err error) {
	defer func() {
		if err != nil {
			err = errors.New("multi mimechecker: " + err.Error())
		}
	}()
	for _, checker := range c.checkers {
		m, err = checker.TypeByFile(path)
		if m != "" || err != nil {
			return
		}
	}
	return
}

func NewMulti(checkers ...Checker) *Multi {
	return &Multi{checkers}
}

