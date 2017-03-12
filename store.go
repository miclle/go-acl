package acl

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// MongoStore struct
type MongoStore struct {
	coll *mgo.Collection
}

// NewMongoStore return store
func NewMongoStore(coll *mgo.Collection) (*MongoStore, error) {

	if err := coll.EnsureIndex(mgo.Index{
		Key:        []string{"g", "gid"},
		Unique:     false,
		DropDups:   true,
		Background: true,
	}); err != nil {
		return nil, err
	}

	if err := coll.EnsureIndex(mgo.Index{
		Key:        []string{"r", "rid"},
		Unique:     false,
		DropDups:   true,
		Background: true,
	}); err != nil {
		return nil, err
	}
	return &MongoStore{coll}, nil
}

// Get permission
func (store *MongoStore) Get(args *Permission) (permission *Permission, err error) {
	if args.ID.Valid() {
		err = store.coll.FindId(args.ID).One(&permission)
		return
	}
	err = store.coll.Find(bson.M{
		"g":   args.GType,
		"gid": args.GID,
		"a":   args.Action,
		"r":   args.RType,
		"rid": args.RID,
	}).One(&permission)
	return
}

// Upsert permission
func (store *MongoStore) Upsert(p *Permission) (err error) {
	_, err = store.coll.Upsert(bson.M{
		"g":   p.GType,
		"gid": p.GID,
		"a":   p.Action,
		"r":   p.RType,
		"rid": p.RID,
	}, p)
	return
}

// Remove permission
func (store *MongoStore) Remove(p *Permission) error {
	if p.ID.Valid() {
		return store.coll.RemoveId(p.ID)
	}
	return store.coll.Remove(bson.M{
		"g":   p.GType,
		"gid": p.GID,
		"a":   p.Action,
		"r":   p.RType,
		"rid": p.RID,
	})
}
