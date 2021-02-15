package fragments

import (
	"errors"
	common "overseer/common/types"
	"overseer/ovsworker/msgheader"
)

//FragmentFactorytMethod - Creates a framgent
type FragmentFactorytMethod func(header msgheader.TaskHeader, data []byte) (WorkFragment, error)

var factories = map[common.TaskType]FragmentFactorytMethod{
	common.TypeDummy: FactoryDummy,
	common.TypeOs:    FactoryOS,
}

//CreateWorkFragment - Creates a task on worker that will be executed
func CreateWorkFragment(header msgheader.TaskHeader, data []byte) (WorkFragment, error) {

	var fragment WorkFragment
	var err error

	method, exists := factories[header.Type]
	if !exists {
		return nil, errors.New("unable to construct fragment")
	}

	if fragment, err = method(header, data); err != nil {
		return nil, err
	}

	return fragment, nil

}
