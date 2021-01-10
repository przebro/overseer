package auth

import (
	"context"
	"goscheduler/datastore"

	"github.com/przebro/databazaar/collection"
)

//RoleAssociationManager - provides basic operations on role association model
type RoleAssociationManager struct {
	col collection.DataCollection
}

func NewRoleAssociationManager(collectionName string, provider *datastore.Provider) (*RoleAssociationManager, error) {

	var col collection.DataCollection
	var err error

	if col, err = provider.GetCollection(collectionName); err != nil {
		return nil, err
	}

	return &RoleAssociationManager{col: col}, nil

}

func (m *RoleAssociationManager) Get(username string) (RoleAssociationModel, bool) {

	dsassoc := dsRoleAssociationModel{}

	if err := m.col.Get(context.Background(), idFormatter(assocNamespace, username), &dsassoc); err != nil {
		return RoleAssociationModel{}, false
	}

	return dsassoc.RoleAssociationModel, false

}
func (m *RoleAssociationManager) Create(model RoleAssociationModel) error {

	dsassoc := dsRoleAssociationModel{RoleAssociationModel: model, ID: idFormatter(assocNamespace, model.Username)}

	_, err := m.col.Create(context.Background(), &dsassoc)
	return err

}
func (m *RoleAssociationManager) Modify(model RoleAssociationModel) error {

	dsassoc := dsRoleAssociationModel{}

	if err := m.col.Get(context.Background(), idFormatter(assocNamespace, model.Username), &dsassoc); err != nil {
		return err
	}

	dsassoc.RoleAssociationModel = model

	return m.col.Update(context.Background(), &dsassoc)

}
func (m *RoleAssociationManager) Delete(username string) error {

	return m.col.Delete(context.Background(), idFormatter(assocNamespace, username))
}
