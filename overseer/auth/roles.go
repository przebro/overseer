package auth

import (
	"context"
	"overseer/datastore"
	"overseer/overseer/config"
	"strings"

	"github.com/przebro/databazaar/collection"
)

//RoleManager - provides basic operations on role model
type RoleManager struct {
	col collection.DataCollection
}

func NewRoleManager(conf config.SecurityConfiguration, provider *datastore.Provider) (*RoleManager, error) {

	var col collection.DataCollection
	var err error

	if col, err = provider.GetCollection(conf.Collection); err != nil {
		return nil, err
	}

	return &RoleManager{col: col}, nil
}

func (m *RoleManager) Get(name string) (RoleModel, bool) {

	dsrole := dsRoleModel{}

	if err := m.col.Get(context.Background(), idFormatter(rolesNamespace, name), &dsrole); err != nil {
		return RoleModel{}, false
	}

	return dsrole.RoleModel, true
}
func (m *RoleManager) Create(model RoleModel) error {

	dsrole := dsRoleModel{RoleModel: model, ID: idFormatter(rolesNamespace, model.Name)}

	_, err := m.col.Create(context.Background(), &dsrole)

	return err
}

func (m *RoleManager) Modify(model RoleModel) error {

	dsrole := dsRoleModel{}

	if err := m.col.Get(context.Background(), idFormatter(rolesNamespace, model.Name), &dsrole); err != nil {
		return err
	}

	dsrole.RoleModel = model

	return m.col.Update(context.Background(), &dsrole)

}

func (m *RoleManager) Delete(name string) error {

	return m.col.Delete(context.Background(), idFormatter(rolesNamespace, name))
}

func (m *RoleManager) All(filter string) ([]RoleModel, error) {

	crsr, err := m.col.All(context.Background())
	if err != nil {
		return nil, err
	}

	rmodel := dsRoleModel{}
	result := []RoleModel{}

	for crsr.Next(context.Background()) {
		if err := crsr.Decode(&rmodel); err != nil {
			return nil, err
		}

		if strings.HasPrefix(rmodel.ID, rolesNamespace) {
			result = append(result, rmodel.RoleModel)
		}
	}

	return result, nil
}
