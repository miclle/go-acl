package acl

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type Staff struct {
	ID   bson.ObjectId
	Name string
}

// ACLType interface
func (staff *Staff) ACLType() string {
	return "staff"
}

// ACLIdentity interface
func (staff *Staff) ACLIdentity() interface{} {
	return staff.ID
}

type Article struct {
	ID int
}

// ACLType interface
func (article *Article) ACLType() string {
	return "article"
}

// ACLIdentity interface
func (article *Article) ACLIdentity() interface{} {
	return fmt.Sprintf("%d", article.ID)
}

func TestAllow(t *testing.T) {
	assert := assert.New(t)

	session, err := mgo.Dial("127.0.0.1:27017")
	assert.Nil(err)

	coll := session.DB("acl_test").C("acls_test")

	_, err = coll.RemoveAll(nil)
	assert.Nil(err)

	store, err := NewMongoStore(coll)
	assert.Nil(err)

	acl := ACL{store}

	staff := &Staff{bson.NewObjectId(), "miclle"}
	article := &Article{123}

	// Allow
	err = acl.Allow(staff, "view", article)
	assert.Nil(err)

	err = acl.Allow(staff, "create", article)
	assert.Nil(err)

	err = acl.Allow(staff, "update", article)
	assert.Nil(err)

	err = acl.Allow(staff, "delete", article)
	assert.Nil(err)

	err = acl.Allow("guest", "view", "doc")
	assert.Nil(err)

	// Can
	can, err := acl.Can(staff, "view", article)
	assert.Nil(err)
	assert.True(can, fmt.Sprintf("expected %v to be able to view %v", staff, article))

	can, err = acl.Can(staff, "delete", article)
	assert.Nil(err)
	assert.True(can, fmt.Sprintf("expected %v to be able to delete %v", staff, article))

	can, err = acl.Can(staff, "rm", article)
	assert.Nil(err)
	assert.False(can, fmt.Sprintf("expected %v to be unable to rm %v", staff, article))

	// RemovePermission
	err = acl.RemovePermission(staff, "delete", article)
	assert.Nil(err)

	can, err = acl.Can(staff, "delete", article)
	assert.Nil(err)
	assert.False(can, fmt.Sprintf("expected %v to be unable to delete %v", staff, article))

	can, err = acl.Can("guest", "view", "doc")
	assert.Nil(err)
	assert.True(can, "expected user to be able to view doc")

	// RemovePermission
	err = acl.RemovePermission("guest", "view", "doc")
	assert.Nil(err)

	can, err = acl.Can("guest", "view", "doc")
	assert.Nil(err)
	assert.False(can, "expected user to be unable to view doc")
}
