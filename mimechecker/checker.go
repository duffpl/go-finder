package mimechecker

type Checker interface {
	ByPath(path string) (string, error)
}