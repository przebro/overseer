package resources

import (
	"context"
	"overseer/datastore"

	"github.com/przebro/databazaar/collection"
)

type ticketReadWriter struct {
	colname  string
	objectID string
	rev      string
	col      collection.DataCollection
}

//newTicketReadWriter - Creates a new readWriter
func newTicketReadWriter(colname, objectID string, provider *datastore.Provider) (readWriter, error) {

	col, err := provider.GetCollection(colname)
	if err != nil {
		return nil, err
	}

	return &ticketReadWriter{colname: colname, col: col, objectID: objectID}, nil
}

//Load - load items from a persistent store
func (cl *ticketReadWriter) Load() (map[string]interface{}, error) {

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

//Write - writes items to a persistent store
func (cl *ticketReadWriter) Write(items map[string]interface{}) error {

	model := []TicketResource{}

	for _, v := range items {
		model = append(model, v.(TicketResource))
	}

	trm := TicketsResourceModel{ID: cl.objectID, REV: cl.rev, Tickets: model}

	err := cl.col.Update(context.Background(), trm)

	return err
}
