package taskdata

//GroupNameData  - helper data for managing task
type GroupNameData struct {
	Group string `validate:"required,max=20"`
	Name  string `validate:"required,max=32"`
}
