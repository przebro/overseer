package launcher

import (
	"errors"
	common "overseer/common/types"
	"overseer/ovsworker/fragments"
	"overseer/ovsworker/msgheader"
)

type factMethod func(header msgheader.TaskHeader, data []byte) (fragments.WorkFragment, error)

var (
	factories map[common.TaskType]factMethod
)

//FragmentCreator - fragment factory
type FragmentCreator struct {
	launcher *FragmentLauncher
}

//FragmentFactory - Creates a new fragment factory
func FragmentFactory(launcher *FragmentLauncher) *FragmentCreator {

	factories = map[common.TaskType]factMethod{
		common.TypeDummy: fragments.FactoryDummy,
		common.TypeOs:    fragments.FactoryOS,
	}
	return &FragmentCreator{launcher: launcher}
}

//CreateFragment - Creates a new fragment thaht will be executed.
func (creator *FragmentCreator) CreateFragment(header msgheader.TaskHeader, data []byte) error {

	var w fragments.WorkFragment = nil
	var err error

	f, exists := factories[header.Type]
	if !exists {
		return errors.New("unable to construct fragment")
	}

	w, err = f(header, data)

	if err != nil {
		return err
	}

	err = creator.launcher.addFragment(w)

	return err
}
