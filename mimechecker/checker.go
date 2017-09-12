package mimechecker

type Checker interface {
	TypeByFile(path string) (string, error)
}