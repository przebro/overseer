package resources

import (
	"context"
	"fmt"

	"github.com/przebro/databazaar/collection"
	"github.com/przebro/databazaar/selector"
	"github.com/przebro/overseer/datastore"
)

type resourceReadWriter struct {
	col collection.DataCollection
}

func newResourceReadWriter(collectionName string, provider datastore.CollectionProvider) (*resourceReadWriter, error) {

	fmt.Println("is provider nil:", provider)
	col, err := provider.GetCollection(context.Background(), collectionName)
	if err != nil {
		return nil, err
	}

	return &resourceReadWriter{col: col}, nil
}
func (r *resourceReadWriter) Get(ctx context.Context, key string, data *ResourceModel) error {

	return r.col.Get(ctx, key, data)
}
func (r *resourceReadWriter) Insert(ctx context.Context, data *ResourceModel) error {

	_, err := r.col.Create(ctx, data)
	return err
}
func (r *resourceReadWriter) Update(ctx context.Context, data *ResourceModel) error {

	return r.col.Update(ctx, data)
}
func (r *resourceReadWriter) Delete(ctx context.Context, key string) error {

	return r.col.Delete(ctx, key)
}

func (r *resourceReadWriter) AllFlags(ctx context.Context) (collection.BazaarCursor, error) {
	q, _ := r.col.AsQuerable()

	var sel selector.Expr

	if r.col.Type() == "badger" {
		sel = selector.Prefix("id", selector.String("f:"))
	} else {
		sel = selector.Eq("type", selector.String("flag"))
	}

	return q.Select(ctx, sel, selector.Fields{"_id", "type", "value"})
}

func (r *resourceReadWriter) AllTickets(ctx context.Context) (collection.BazaarCursor, error) {
	q, _ := r.col.AsQuerable()

	var sel selector.Expr

	if r.col.Type() == "badger" {
		sel = selector.Prefix("id", selector.String("t:"))
	} else {
		sel = selector.Eq("type", selector.String("ticket"))
	}
	return q.Select(ctx, sel, selector.Fields{"_id", "type", "value"})
}
