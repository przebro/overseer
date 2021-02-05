package auth

//UserModel -
type UserModel struct {
	Username string `json:"username" validate:"required,max=32"`
	FullName string `json:"fullname" validate:"max=64,auth"`
	Password string `json:"password" validate:"required"`
	Mail     string `json:"mail" validate:"email,max=64"`
	Enabled  bool   `json:"enabled"`
}

//datastoreUserModel - internal database representation of a UserModel
type dsUserModel struct {
	UserModel
	ID  string `json:"_id" bson:"_id" validate:"required,max=37"` // prefix='user' + @ + username
	Rev string `json:"_rev,omitempty" bson:"_rev,omitempty"`
}

//RoleAssociationModel -
type RoleAssociationModel struct {
	UserID string   `json:"userid" validate:"required,max=32"`
	Roles  []string `json:"uroles"`
}

type dsRoleAssociationModel struct {
	RoleAssociationModel
	ID  string `json:"_id" bson:"_id" validate:"required,max=38"` // prefix='assoc' + @ + username
	Rev string `json:"_rev,omitempty" bson:"_rev,omitempty"`
}

//RoleModel -
type RoleModel struct {
	Name        string `json:"name" validate:"required,max=32"`
	Description string `json:"description" validate:"max=64,auth"`
	//ServerManagement
	//UsersManagement
	Administration bool `json:"administration"`
	//Active Pool Management
	Restart      bool `json:"restart"`
	SetToOK      bool `json:"setok"`
	AddTicket    bool `json:"addticket"`
	RemoveTicket bool `json:"removeticket"`
	SetFlag      bool `json:"setflag"`
	Confirm      bool `json:"confirm"`
	Bypass       bool `json:"bypass"`
	Hold         bool `json:"hold"`
	Free         bool `json:"free"`
	//Definition Management
	Order      bool `json:"order"`
	Force      bool `json:"force"`
	Definition bool `json:"definition"`
}

type dsRoleModel struct {
	RoleModel
	ID  string `json:"_id" bson:"_id" validate:"required,max=37"` // prefix='role' + @ + name
	Rev string `json:"_rev,omitempty" bson:"_rev,omitempty"`
}
