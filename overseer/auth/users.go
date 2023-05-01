package auth

import (
	"context"
	"fmt"
	"sync"

	"github.com/przebro/overseer/datastore"
	"github.com/przebro/overseer/overseer/config"

	"golang.org/x/crypto/bcrypt"

	"github.com/przebro/databazaar/collection"
	"github.com/przebro/databazaar/selector"
)

// UserManager - provides basic operations on user model
type UserManager struct {
	lock sync.RWMutex
	col  collection.DataCollection
}

// NewUserManager - Creates an new instance of UserManager
func NewUserManager(conf config.SecurityConfiguration, provider *datastore.Provider) (*UserManager, error) {

	var col collection.DataCollection
	var err error

	if col, err = provider.GetCollection(context.Background(), collectionName); err != nil {
		return nil, err
	}

	return &UserManager{col: col, lock: sync.RWMutex{}}, nil
}

// Get - gets a user, returns empty model and false if user not found
func (m *UserManager) Get(ctx context.Context, username string) (UserModel, bool) {

	m.lock.RLock()
	defer m.lock.RUnlock()

	model := dsUserModel{}

	if err := m.col.Get(ctx, idFormatter(userNamespace, username), &model); err != nil {

		return UserModel{}, false
	}

	return model.UserModel, true
}

// Create - creates a new user
func (m *UserManager) Create(ctx context.Context, model UserModel) error {

	m.lock.Lock()
	defer m.lock.Unlock()

	dsmodel := dsUserModel{UserModel: model, ID: idFormatter(userNamespace, model.Username)}

	if err := m.col.Get(ctx, idFormatter(userNamespace, model.Username), &dsUserModel{}); err == nil {
		return fmt.Errorf("user %s already exists", model.Username)
	}

	_, err := m.col.Create(ctx, &dsmodel)
	return err
}

// Modify - modifies a user
func (m *UserManager) Modify(ctx context.Context, model UserModel) error {

	m.lock.Lock()
	defer m.lock.Unlock()

	dsmodel := dsUserModel{}

	if err := m.col.Get(ctx, idFormatter(userNamespace, model.Username), &dsmodel); err != nil {
		return err
	}

	dsmodel.UserModel = model

	return m.col.Update(context.Background(), &dsmodel)
}

// Delete - deletes a user
func (m *UserManager) Delete(ctx context.Context, username string) error {

	m.lock.Lock()
	defer m.lock.Unlock()

	return m.col.Delete(ctx, idFormatter(userNamespace, username))
}

// All - returns a list of users
func (m *UserManager) All(filter string) ([]UserModel, error) {

	m.lock.RLock()
	defer m.lock.RUnlock()

	q, err := m.col.AsQuerable()
	if err != nil {
		return nil, err
	}
	var sel selector.Expr

	if q.Type() == "badger" {
		sel = selector.Prefix("_id", selector.String(userNamespace))
	} else {
		sel = selector.Eq("_id", selector.String(userNamespace))
	}

	crsr, err := q.Select(context.Background(), sel, nil)
	if err != nil {
		return nil, err
	}

	result := []UserModel{}

	for crsr.Next(context.Background()) {
		umodel := dsUserModel{}
		if err := crsr.Decode(&umodel); err != nil {
			return nil, err
		}

		result = append(result, umodel.UserModel)
	}

	return result, nil
}

// CheckChangePassword - checks if an old password match and if succeed, create and returns a new one
func (m *UserManager) CheckChangePassword(crypt, old, new []byte) (string, error) {

	var err error

	if err = bcrypt.CompareHashAndPassword(crypt, old); err != nil {
		return "", err
	}

	return HashPassword(new)
}
