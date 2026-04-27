package modifytags

type Transform int

const (
	CamelCase = iota
	SnakeCase
	PascalCase
)

type Modification struct {
	Add       []string
	Remove    []string
	Transform Transform
	Override  bool
	Sort      bool
	Clear     bool
}

func (mod *Modification) Apply() error {
	return nil
}

func (mod *Modification) rewrite() error {
	return nil
}
