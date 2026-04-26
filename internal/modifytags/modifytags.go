package modifytags

import "errors"

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
	if err := mod.validate(); err != nil {
		return err
	}

	return nil
}

func (mod *Modification) validate() error {
	if len(mod.Add) == 0 && !mod.Clear && len(mod.Remove) == 0 {
		return errors.New("one of [-add-tags, -add-options, -remove-tags, -remove-options, -clear-tags," +
			" -clear-options] should be defined")
	}

	return nil
}
