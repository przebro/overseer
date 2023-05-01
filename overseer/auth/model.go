package auth

// UserModel -
type UserModel struct {
	Username string `json:"username" validate:"required,max=50,username"`
	FullName string `json:"fullname" validate:"max=100,auth"`
	Password string `json:"password" validate:"required"`
	Mail     string `json:"mail" validate:"email,max=128"`
	Enabled  bool   `json:"enabled"`
}

type UserDetailModel struct {
	UserModel
	Roles []string `json:"roles"`
}

// datastoreUserModel - internal database representation of a UserModel
type dsUserModel struct {
	ID  string `json:"_id" bson:"_id" validate:"required,max=55"` // prefix='user' + @ + username
	Rev string `json:"_rev,omitempty" bson:"_rev,omitempty"`
	UserModel
}

// RoleAssociationModel -
type RoleAssociationModel struct {
	UserID string   `json:"userid" validate:"required,max=50"`
	Roles  []string `json:"uroles"`
}

type dsRoleAssociationModel struct {
	RoleAssociationModel
	ID  string `json:"_id" bson:"_id" validate:"required,max=58"` // prefix='assoc' + @ + username
	Rev string `json:"_rev,omitempty" bson:"_rev,omitempty"`
}

// RoleModel -
type RoleModel struct {
	Name        string `json:"name" validate:"required,max=32"`
	Description string `json:"description" validate:"max=100,auth"`
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
