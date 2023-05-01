package auth

/*
var aprovider *datastore.Provider

func aprepare(t *testing.T) {

	if provider != nil {
		return
	}

	var err error
	f, _ := os.Create("../../data/tests/authtest.json")
	f.Write([]byte("{}"))
	f.Close()

	provider, err = datastore.NewDataProvider(astorecfg)
	if err != nil {
		t.Fatal("unable to init store")
	}

}

func TestNewAssocManager(t *testing.T) {

	aprepare(t)

	_, err := NewRoleAssociationManager(aconf, provider)
	if err != nil {
		t.Error("unexpected resutlt:", err)
	}
}
func TestCreateAssociation(t *testing.T) {

	m, err := NewRoleAssociationManager(aconf, provider)
	if err != nil {
		t.Error("unexpected resutlt:", err)
	}

	model := RoleAssociationModel{
		Roles:  []string{"test", "admin"},
		UserID: "testuser",
	}

	err = m.Create(model)
	if err != nil {
		t.Error("unexpected result:", err)
	}

	err = m.Create(model)
	if err == nil {
		t.Error("unexpected result:", err)
	}

}
func TestGetAssociation(t *testing.T) {

	m, err := NewRoleAssociationManager(aconf, provider)
	if err != nil {
		t.Error("unexpected resutlt:", err)
	}

	_, ok := m.Get("testuser")

	if ok != true {
		t.Error("unexpected result:", ok)
	}

	_, ok = m.Get("testuser2")

	if ok == true {
		t.Error("unexpected result:", ok)
	}
}

func TestModifyAssociation(t *testing.T) {

	m, err := NewRoleAssociationManager(aconf, provider)
	if err != nil {
		t.Error("unexpected resutlt:", err)
	}

	model := RoleAssociationModel{
		Roles:  []string{"test", "admin"},
		UserID: "testuser2",
	}

	err = m.Modify(model)
	if err == nil {
		t.Error("unexpected result:", nil)
	}
	err = m.Create(model)

	if err != nil {
		t.Error("unexpected result:", err)
	}

	model.Roles = []string{"admin", "other"}

	err = m.Modify(model)
	if err != nil {
		t.Error("unexpected result:", err)
	}
}

func TestDeleteAssociation(t *testing.T) {

	m, err := NewRoleAssociationManager(aconf, provider)
	if err != nil {
		t.Error("unexpected resutlt:", err)
	}

	err = m.Delete("testuser3")
	if err != nil {
		t.Error("unexpected result:", err)
	}

	err = m.Delete("testuser2")
	if err != nil {
		t.Error("unexpected result:", err)
	}

}
*/
