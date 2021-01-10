package auth

import (
	"context"
	"overseer/datastore"

	"github.com/przebro/databazaar/collection"
)

//UserManager - provides basic operations on user model
type UserManager struct {
	col collection.DataCollection
}

func NewUserManager(collectionName string, provider *datastore.Provider) (*UserManager, error) {

	var col collection.DataCollection
	var err error

	if col, err = provider.GetCollection(collectionName); err != nil {
		return nil, err
	}

	return &UserManager{col: col}, nil
}

func (m *UserManager) Get(username string) (UserModel, bool) {

	model := dsUserModel{}
	if err := m.col.Get(context.Background(), username, &model); err != nil {

		return UserModel{}, false
	}

	return model.UserModel, true
}
func (m *UserManager) Create(model UserModel) error {

	dsmodel := dsUserModel{UserModel: model, ID: idFormatter(userNamespace, model.Username)}

	_, err := m.col.Create(context.Background(), &dsmodel)
	return err
}
func (m *UserManager) Modify(model UserModel) error {

	dsmodel := dsUserModel{}

	if err := m.col.Get(context.Background(), idFormatter(userNamespace, model.Username), &dsmodel); err != nil {
		return err
	}

	dsmodel.UserModel = model

	return m.col.Update(context.Background(), &dsmodel)
}
func (m *UserManager) Delete(username string) error {

	return m.col.Delete(context.Background(), idFormatter(userNamespace, username))
}
