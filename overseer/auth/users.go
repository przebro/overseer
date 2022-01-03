package auth

import (
	"context"
	"strings"

	"github.com/przebro/overseer/datastore"
	"github.com/przebro/overseer/overseer/config"

	"golang.org/x/crypto/bcrypt"

	"github.com/przebro/databazaar/collection"
)

//UserManager - provides basic operations on user model
type UserManager struct {
	col collection.DataCollection
}

//NewUserManager - Creates an new instance of UserManager
func NewUserManager(conf config.SecurityConfiguration, provider *datastore.Provider) (*UserManager, error) {

	var col collection.DataCollection
	var err error

	if col, err = provider.GetCollection(conf.Collection); err != nil {
		return nil, err
	}

	return &UserManager{col: col}, nil
}

//Get - gets a user, returns empty model and false if user not found
func (m *UserManager) Get(username string) (UserModel, bool) {

	model := dsUserModel{}

	if err := m.col.Get(context.Background(), idFormatter(userNamespace, username), &model); err != nil {

		return UserModel{}, false
	}

	return model.UserModel, true
}

//Create - creates a new user
func (m *UserManager) Create(model UserModel) error {

	dsmodel := dsUserModel{UserModel: model, ID: idFormatter(userNamespace, model.Username)}

	_, err := m.col.Create(context.Background(), &dsmodel)
	return err
}

//Modify - modifies a user
func (m *UserManager) Modify(model UserModel) error {

	dsmodel := dsUserModel{}

	if err := m.col.Get(context.Background(), idFormatter(userNamespace, model.Username), &dsmodel); err != nil {
		return err
	}

	dsmodel.UserModel = model

	return m.col.Update(context.Background(), &dsmodel)
}

//Delete - deletes a user
func (m *UserManager) Delete(username string) error {

	return m.col.Delete(context.Background(), idFormatter(userNamespace, username))
}

//All - returns a list of users
func (m *UserManager) All(filter string) ([]UserModel, error) {

	crsr, err := m.col.All(context.Background())
	if err != nil {
		return nil, err
	}

	umodel := dsUserModel{}
	result := []UserModel{}

	for crsr.Next(context.Background()) {
		if err := crsr.Decode(&umodel); err != nil {
			return nil, err
		}

		if strings.HasPrefix(umodel.ID, userNamespace) {
			result = append(result, umodel.UserModel)
		}
	}

	return result, nil
}

//CheckChangePassword - checks if an old password match and if succeed, create and returns a new one
func (m *UserManager) CheckChangePassword(crypt, old, new []byte) (string, error) {

	var err error

	if err = bcrypt.CompareHashAndPassword(crypt, old); err != nil {
		return "", err
	}

	return HashPassword(new)
}
