package bolt

import (
	"context"
	g_errors "errors"
	"os"
	"reflect"
	"testing"

	"github.com/asdine/storm/v3"
	"github.com/rghiorghisor/basic-go-rest-api/errors"
	"github.com/rghiorghisor/basic-go-rest-api/model"
	"github.com/rghiorghisor/basic-go-rest-api/util"
	"gopkg.in/go-playground/assert.v1"
)

var defaultDir = "../../../../tests/local-repo"
var defaultDB = "../../../../tests/local-repo/testsdb"

func TestCreateUnexpected(t *testing.T) {
	repo := setup()
	repo.db.Close()
	defer tearDown(repo)

	prop1 := &model.PropertySet{Name: "test.name.1", Values: []string{"test.value.1.1", "test.value.1.2"}}

	err := repo.Create(context.Background(), prop1)
	assert.Equal(t, g_errors.New("database not open"), err)

}

func TestReadAll(t *testing.T) {
	repo := setup()
	defer tearDown(repo)

	prop1 := &model.PropertySet{Name: "test.name.1", Values: []string{"test.value.1.1", "test.value.1.2"}}
	prop2 := &model.PropertySet{Name: "test.name.2", Values: []string{"test.value.2.1", "test.value.2.2"}}

	repo.Create(context.Background(), prop1)
	repo.Create(context.Background(), prop2)

	readProps, _ := repo.ReadAll(context.Background())

	if readProps[0].Name == prop2.Name {
		tmp := prop1
		prop1 = prop2
		prop2 = tmp
	}

	assert.Equal(t, true, (readProps[0].Name != ""))
	assert.Equal(t, prop1.Name, readProps[0].Name)
	assert.Equal(t, prop1.Values, readProps[0].Values)
	assert.Equal(t, true, (readProps[1].Name != ""))
	assert.Equal(t, prop2.Name, readProps[1].Name)
	assert.Equal(t, prop2.Values, readProps[1].Values)
}

func TestReadAllUnexpected(t *testing.T) {
	repo := setup()
	repo.db.Close()

	defer tearDown(repo)

	_, err := repo.ReadAll(context.Background())
	assert.Equal(t, g_errors.New("database not open"), err)
}

func TestFindById(t *testing.T) {
	repo := setup()
	defer tearDown(repo)

	prop1 := &model.PropertySet{Name: "test.name.1", Values: []string{"test.value.1.1", "test.value.1.2"}}
	prop2 := &model.PropertySet{Name: "test.name.2", Values: []string{"test.value.2.1", "test.value.2.2"}}

	repo.Create(context.Background(), prop1)
	repo.Create(context.Background(), prop2)

	var readProps []*model.PropertySet
	readProps = make([]*model.PropertySet, 0, 2)

	prop, _ := repo.FindByID(context.Background(), prop1.Name)
	readProps = append(readProps, prop)

	prop, _ = repo.FindByID(context.Background(), prop2.Name)
	readProps = append(readProps, prop)

	assert.Equal(t, true, (readProps[0].Name != ""))
	assert.Equal(t, prop1.Name, readProps[0].Name)
	assert.Equal(t, prop1.Values, readProps[0].Values)
	assert.Equal(t, true, (readProps[1].Name != ""))
	assert.Equal(t, prop2.Name, readProps[1].Name)
	assert.Equal(t, prop2.Values, readProps[1].Values)
}

func TestFindByIdNotFound(t *testing.T) {
	repo := setup()
	defer tearDown(repo)

	id := "test.notfound.id"
	_, err := repo.FindByID(context.Background(), id)

	assert.Equal(t, errors.NewEntityNotFound(reflect.TypeOf((*model.PropertySet)(nil)).Elem(), id), err)
}

func TestFindByIdUnexpected(t *testing.T) {
	repo := setup()
	repo.db.Close()
	defer tearDown(repo)

	id := "test.notfound.id"
	_, err := repo.FindByID(context.Background(), id)

	assert.Equal(t, g_errors.New("database not open"), err)
}

func TestFindByName(t *testing.T) {
	repo := setup()
	defer tearDown(repo)

	prop1 := &model.PropertySet{Name: "test.name.1", Values: []string{"test.value.1.1", "test.value.1.2"}}
	prop2 := &model.PropertySet{Name: "test.name.2", Values: []string{"test.value.2.1", "test.value.2.2"}}

	repo.Create(context.Background(), prop1)
	repo.Create(context.Background(), prop2)

	var readProps []*model.PropertySet
	readProps = make([]*model.PropertySet, 0, 2)

	prop, _ := repo.FindByName(context.Background(), prop1.Name)
	readProps = append(readProps, prop)

	prop, _ = repo.FindByName(context.Background(), prop2.Name)
	readProps = append(readProps, prop)

	assert.Equal(t, true, (readProps[0].Name != ""))
	assert.Equal(t, prop1.Name, readProps[0].Name)
	assert.Equal(t, prop1.Values, readProps[0].Values)
	assert.Equal(t, true, (readProps[1].Name != ""))
	assert.Equal(t, prop2.Name, readProps[1].Name)
	assert.Equal(t, prop2.Values, readProps[1].Values)
}

func TestFindByNameNotFound(t *testing.T) {
	repo := setup()
	defer tearDown(repo)

	name := "test.notfound.name"
	_, err := repo.FindByName(context.Background(), name)

	assert.Equal(t, errors.NewEntityNotFound(reflect.TypeOf((*model.PropertySet)(nil)).Elem(), name), err)
}

func TestFindByNameUnexpected(t *testing.T) {
	repo := setup()
	repo.db.Close()
	defer tearDown(repo)

	name := "test.notfound.name"
	_, err := repo.FindByName(context.Background(), name)

	assert.Equal(t, g_errors.New("database not open"), err)
}

func TestDelete(t *testing.T) {
	repo := setup()
	defer tearDown(repo)

	prop1 := &model.PropertySet{Name: "test.name.1", Values: []string{"test.value.1.1", "test.value.1.2"}}

	repo.Create(context.Background(), prop1)

	readProps, _ := repo.ReadAll(context.Background())

	assert.Equal(t, true, (readProps[0].Name != ""))
	assert.Equal(t, prop1.Name, readProps[0].Name)
	assert.Equal(t, prop1.Values, readProps[0].Values)

	repo.Delete(context.Background(), readProps[0].Name)
	id := readProps[0].Name
	_, err := repo.FindByID(context.Background(), id)

	assert.Equal(t, errors.NewEntityNotFound(reflect.TypeOf((*model.PropertySet)(nil)).Elem(), id), err)
}

func TestDeleteNotFound(t *testing.T) {
	repo := setup()
	defer tearDown(repo)

	id := "test.notfound.id"
	err := repo.Delete(context.Background(), id)
	assert.Equal(t, errors.NewEntityNotFound(reflect.TypeOf((*model.PropertySet)(nil)).Elem(), id), err)
}

func TestDeleteUnexpectedFind(t *testing.T) {
	repo := setup()
	repo.db.Close()
	defer tearDown(repo)

	id := "test.notfound.id"
	err := repo.Delete(context.Background(), id)

	assert.Equal(t, g_errors.New("database not open"), err)
}

func TestUpdate(t *testing.T) {
	repo := setup()
	defer tearDown(repo)

	prop1 := &model.PropertySet{Name: "test.name.1", Values: []string{"test.value.1.1", "test.value.1.2"}}

	repo.Create(context.Background(), prop1)

	readProps, _ := repo.ReadAll(context.Background())
	assert.Equal(t, true, (readProps[0].Name != ""))
	assert.Equal(t, prop1.Name, readProps[0].Name)
	assert.Equal(t, prop1.Values, readProps[0].Values)

	expectedValues := append(prop1.Values, "test.value.1.3")
	prop1.Values = expectedValues

	repo.Update(context.Background(), prop1)

	found, _ := repo.FindByID(context.Background(), prop1.Name)
	assert.Equal(t, found.Name, prop1.Name)
	assert.Equal(t, found.Values, prop1.Values)
}

func TestUpdateNotFound(t *testing.T) {
	repo := setup()
	defer tearDown(repo)

	prop1 := &model.PropertySet{Name: "test.name.1", Values: []string{"test.value.1.1", "test.value.1.2"}}

	repo.Create(context.Background(), prop1)

	readProps, _ := repo.ReadAll(context.Background())
	assert.Equal(t, true, (readProps[0].Name != ""))
	assert.Equal(t, prop1.Name, readProps[0].Name)
	assert.Equal(t, prop1.Values, readProps[0].Values)

	expectedValues := append(prop1.Values, "test.value.1.3")
	prop1.Values = expectedValues

	repo.Delete(context.Background(), prop1.Name)
	err := repo.Update(context.Background(), prop1)
	assert.Equal(t, errors.NewEntityNotFound(reflect.TypeOf((*model.PropertySet)(nil)).Elem(), prop1.Name), err)
}

func TestUpdateUnExpected(t *testing.T) {
	repo := setup()

	defer tearDown(repo)

	prop1 := &model.PropertySet{Name: "test.name.1", Values: []string{"test.value.1.1", "test.value.1.2"}}

	repo.Create(context.Background(), prop1)

	readProps, _ := repo.ReadAll(context.Background())
	assert.Equal(t, true, (readProps[0].Name != ""))
	assert.Equal(t, prop1.Name, readProps[0].Name)
	assert.Equal(t, prop1.Values, readProps[0].Values)

	expectedValues := append(prop1.Values, "test.value.1.3")
	prop1.Values = expectedValues
	repo.db.Close()

	err := repo.Update(context.Background(), prop1)
	assert.Equal(t, g_errors.New("database not open"), err)
}

func BenchmarkReadAll(b *testing.B) {
	repo := setup()
	defer tearDown(repo)

	prop1 := &model.PropertySet{Name: "test.name.1", Values: []string{"test.value.1.1", "test.value.1.2"}}
	prop2 := &model.PropertySet{Name: "test.name.2", Values: []string{"test.value.2.1", "test.value.2.2"}}

	repo.Create(context.Background(), prop1)
	repo.Create(context.Background(), prop2)

	for n := 0; n < b.N; n++ {
		repo.ReadAll(context.Background())
	}

}

func setup() *PropertySetRepository {
	util.CreateParentFolder(defaultDB)

	db, _ := storm.Open(defaultDB)

	return &PropertySetRepository{
		db: db,
	}
}

func tearDown(repo *PropertySetRepository) {
	repo.db.Close()

	os.Remove(defaultDB)
	os.Remove(defaultDir)
}
