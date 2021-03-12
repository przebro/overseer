package resources

import (
	"context"
	"overseer/datastore"

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

//newFlagReadWriter - creates a new readWriter
func newFlagReadWriter(colname, objectID string, provider *datastore.Provider) (readWriter, error) {

	col, err := provider.GetCollection(colname)

	if err != nil {
		return nil, err
	}

	return &flagReadWriter{colname: colname, col: col, objectID: objectID}, nil
}

//Load - load items from a persistent store
func (cl *flagReadWriter) Load() (map[string]interface{}, error) {

	model := FlagsResourceModel{Flags: []FlagResource{}}

	err := cl.col.Get(context.Background(), cl.objectID, &model)
	if err != nil {
		if err == collection.ErrNoDocuments {
			model.ID = cl.objectID
			cl.col.Create(context.Background(), &model)

		} else {
			return nil, err
		}
	}
	cl.rev = model.REV

	data := map[string]interface{}{}

	for _, t := range model.Flags {

		key := t.Name
		data[key] = FlagResource{Name: t.Name, Count: t.Count, Policy: t.Policy}
	}

	return data, nil
}

//Write - writes items to the persistent store
func (cl *flagReadWriter) Write(items map[string]interface{}) error {

	model := []FlagResource{}

	for _, v := range items {
		model = append(model, v.(FlagResource))
	}

	frm := FlagsResourceModel{ID: cl.objectID, REV: cl.rev, Flags: model}

	err := cl.col.Update(context.Background(), frm)

	return err
}
