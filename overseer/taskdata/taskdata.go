package taskdata

//GroupData  - helper data for managing task
type GroupData struct {
	Group string `validate:"required,max=20,resname"`
}

//GroupNameData  - helper data for managing task
type GroupNameData struct {
	GroupData
	Name string `validate:"required,max=32,resname"`
}

//TaskNameModel  - contains base task properties
type TaskNameModel struct {
	Group       string
	Name        string
	Description string
}
