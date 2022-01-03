package auth

import (
	"context"

	"github.com/przebro/overseer/datastore"
	"github.com/przebro/overseer/overseer/config"

	"github.com/przebro/databazaar/collection"
)

//RoleAssociationManager - provides basic operations on role association model
type RoleAssociationManager struct {
	col collection.DataCollection
}

//NewRoleAssociationManager - creates a new RoleAssociationManager
func NewRoleAssociationManager(conf config.SecurityConfiguration, provider *datastore.Provider) (*RoleAssociationManager, error) {

	var col collection.DataCollection
	var err error

	if col, err = provider.GetCollection(conf.Collection); err != nil {
		return nil, err
	}

	return &RoleAssociationManager{col: col}, nil

}

//Get - gets a role association with given user
func (m *RoleAssociationManager) Get(username string) (RoleAssociationModel, bool) {

	dsassoc := dsRoleAssociationModel{}

	if err := m.col.Get(context.Background(), idFormatter(assocNamespace, username), &dsassoc); err != nil {
		return RoleAssociationModel{}, false
	}

	return dsassoc.RoleAssociationModel, true

}

//Create - creates a new role association
func (m *RoleAssociationManager) Create(model RoleAssociationModel) error {

	dsassoc := dsRoleAssociationModel{RoleAssociationModel: model, ID: idFormatter(assocNamespace, model.UserID)}

	_, err := m.col.Create(context.Background(), &dsassoc)
	return err

}

//Modify - modifies a role association
func (m *RoleAssociationManager) Modify(model RoleAssociationModel) error {

	dsassoc := dsRoleAssociationModel{}

	if err := m.col.Get(context.Background(), idFormatter(assocNamespace, model.UserID), &dsassoc); err != nil {
		return err
	}

	dsassoc.RoleAssociationModel = model

	return m.col.Update(context.Background(), &dsassoc)

}

//Delete - deletes a role association
func (m *RoleAssociationManager) Delete(username string) error {

	return m.col.Delete(context.Background(), idFormatter(assocNamespace, username))
}
