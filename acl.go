package acl

import (
	"fmt"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// Type interface
type Type interface {
	ACLType() string
}

// Identity interface
type Identity interface {
	ACLIdentity() interface{}
}

// Permission struct
type Permission struct {
	ID     bson.ObjectId `bson:"_id,omitempty" json:"id"`
	GType  string        `bson:"g"             json:"g"`
	GID    interface{}   `bson:"gid"           json:"gid"`
	Action string        `bson:"a"             json:"a"`
	RType  string        `bson:"r"             json:"r"`
	RID    interface{}   `bson:"rid"           json:"rid"`
}

// Store interface
type Store interface {
	Get(args *Permission) (*Permission, error)
	Upsert(perm *Permission) error
	Remove(perm *Permission) error
}

// ACL access control list
type ACL struct {
	store Store
}

// New return ACL
func New(store Store) (*ACL, error) {
	return &ACL{store}, nil
}

// Allow permission
func (acl *ACL) Allow(rolable interface{}, action string, resource interface{}) (err error) {
	perm, err := permissionFromInterface(rolable, action, resource)
	if err != nil {
		return err
	}
	return acl.store.Upsert(perm)
}

// RemovePermission remove permission
func (acl *ACL) RemovePermission(rolable interface{}, action string, resource interface{}) error {
	perm, err := permissionFromInterface(rolable, action, resource)
	if err != nil {
		return err
	}
	return acl.store.Remove(perm)
}

// Can able to action resource
func (acl *ACL) Can(rolable interface{}, action string, resource interface{}) (bool, error) {
	permission, err := permissionFromInterface(rolable, action, resource)
	if err != nil {
		return false, err
	}
	perm, err := acl.store.Get(permission)
	if err == mgo.ErrNotFound {
		return false, nil
	}
	if perm != nil {
		return true, nil
	}
	return false, err
}

func rolableFromInterface(rolable interface{}) (aclType string, identity interface{}) {
	identity = "*"
	if aclType, ok := rolable.(string); ok {
		return aclType, identity
	}
	if i, ok := rolable.(Type); ok {
		aclType = i.ACLType()
	}
	if i, ok := rolable.(Identity); ok {
		identity = i.ACLIdentity()
	}
	return
}

func permissionFromInterface(rolable interface{}, action string, resource interface{}) (*Permission, error) {
	ptype, pid := rolableFromInterface(rolable)
	if ptype == "" {
		return nil, fmt.Errorf("rolable should be a string, ACLType or ACLIdentity. Was: %v", rolable)
	}

	rtype, rid := rolableFromInterface(resource)
	if ptype == "" {
		return nil, fmt.Errorf("resource should be a string, ACLType or ACLIdentity. Was: %v", resource)
	}

	return &Permission{
		GType:  ptype,
		GID:    pid,
		Action: action,
		RType:  rtype,
		RID:    rid,
	}, nil
}
