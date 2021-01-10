package resources

import (
	"context"
	"goscheduler/datastore"

	"github.com/przebro/databazaar/collection"
)

type readWriter interface {
	Load() (map[string]interface{}, error)
	Write(items map[string]interface{}) error
}

type flagReadWriter struct {
	colname  string
	objectID string
	rev      string
	col      collection.DataCollection
}

func NewFlagReadWriter(colname, objectID string, provider *datastore.Provider) (readWriter, error) {

	col, err := provider.GetCollection(colname)

	if err != nil {
		return nil, err
	}

	return &flagReadWriter{colname: colname, col: col, objectID: objectID}, nil

}

func (cl *flagReadWriter) Load() (map[string]interface{}, error) {

	model := TicketsResourceModel{Tickets: []TicketResource{}}

	err := cl.col.Get(context.Background(), cl.objectID, &model)
	if err != nil {
		return nil, err
	}
	cl.rev = model.REV

	data := map[string]interface{}{}

	for _, t := range model.Tickets {

		key := t.Name + string(t.Odate)
		data[key] = TicketResource{Name: t.Name, Odate: t.Odate}
	}

	return data, nil
}
func (cl *flagReadWriter) Write(items map[string]interface{}) error {

	model := []FlagResource{}

	for _, v := range items {
		model = append(model, v.(FlagResource))
	}

	frm := FlagsResourceModel{ID: cl.objectID, REV: cl.rev}

	err := cl.col.Update(context.Background(), frm)

	return err
}
