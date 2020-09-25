package bolt

import (
	"context"
	g_errors "errors"
	"os"
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

	prop1 := &model.Property{
		Name:        "test.name.1",
		Description: "test.description.1",
		Value:       "test.value.1",
	}

	err := repo.Create(context.Background(), prop1)
	assert.Equal(t, g_errors.New("database not open"), err)

}

func TestReadAll(t *testing.T) {
	repo := setup()
	defer tearDown(repo)

	prop1 := &model.Property{Name: "test.name.1", Description: "test.description.1", Value: "test.value.1"}
	prop2 := &model.Property{Name: "test.name.2", Description: "test.description.2", Value: "test.value.2"}

	repo.Create(context.Background(), prop1)
	repo.Create(context.Background(), prop2)

	readProps, _ := repo.ReadAll(context.Background())

	if readProps[0].ID == prop2.ID {
		tmp := prop1
		prop1 = prop2
		prop2 = tmp
	}

	assert.Equal(t, true, (readProps[0].ID != ""))
	assert.Equal(t, prop1.Name, readProps[0].Name)
	assert.Equal(t, prop1.Description, readProps[0].Description)
	assert.Equal(t, prop1.Value, readProps[0].Value)
	assert.Equal(t, true, (readProps[1].ID != ""))
	assert.Equal(t, prop2.Name, readProps[1].Name)
	assert.Equal(t, prop2.Description, readProps[1].Description)
	assert.Equal(t, prop2.Value, readProps[1].Value)
}

func TestReadAllFiltered(t *testing.T) {
	repo := setup()
	defer tearDown(repo)

	prop1 := &model.Property{Name: "test.name.1", Description: "test.description.1", Value: "test.value.1"}
	prop2 := &model.Property{Name: "test.name.2", Description: "test.description.2", Value: "test.value.2"}
	prop3 := &model.Property{Name: "test.name.3", Description: "test.description.3", Value: "test.value.3"}

	repo.Create(context.Background(), prop1)
	repo.Create(context.Background(), prop2)
	repo.Create(context.Background(), prop3)

	readProps, _ := repo.ReadAllFiltered(context.Background(), []string{prop1.Name, prop2.Name})

	if readProps[0].ID == prop2.ID {
		tmp := prop1
		prop1 = prop2
		prop2 = tmp
	}

	assert.Equal(t, true, (readProps[0].ID != ""))
	assert.Equal(t, prop1.Name, readProps[0].Name)
	assert.Equal(t, prop1.Description, readProps[0].Description)
	assert.Equal(t, prop1.Value, readProps[0].Value)
	assert.Equal(t, true, (readProps[1].ID != ""))
	assert.Equal(t, prop2.Name, readProps[1].Name)
	assert.Equal(t, prop2.Description, readProps[1].Description)
	assert.Equal(t, prop2.Value, readProps[1].Value)
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

	prop1 := &model.Property{
		Name:        "test.name.1",
		Description: "test.description.1",
		Value:       "test.value.1",
	}
	prop2 := &model.Property{
		Name:        "test.name.2",
		Description: "test.description.2",
		Value:       "test.value.2",
	}

	repo.Create(context.Background(), prop1)
	repo.Create(context.Background(), prop2)

	var readProps []*model.Property
	readProps = make([]*model.Property, 0, 2)

	prop, _ := repo.FindByID(context.Background(), prop1.ID)
	readProps = append(readProps, prop)

	prop, _ = repo.FindByID(context.Background(), prop2.ID)
	readProps = append(readProps, prop)

	assert.Equal(t, true, (readProps[0].ID != ""))
	assert.Equal(t, prop1.Name, readProps[0].Name)
	assert.Equal(t, prop1.Description, readProps[0].Description)
	assert.Equal(t, prop1.Value, readProps[0].Value)
	assert.Equal(t, true, (readProps[1].ID != ""))
	assert.Equal(t, prop2.Name, readProps[1].Name)
	assert.Equal(t, prop2.Description, readProps[1].Description)
	assert.Equal(t, prop2.Value, readProps[1].Value)
}

func TestFindByIdNotFound(t *testing.T) {
	repo := setup()
	defer tearDown(repo)

	id := "test.notfound.id"
	_, err := repo.FindByID(context.Background(), id)

	assert.Equal(t, errors.NewEntityNotFound(model.Property{}, id), err)
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

	prop1 := &model.Property{
		Name:        "test.name.1",
		Description: "test.description.1",
		Value:       "test.value.1",
	}
	prop2 := &model.Property{
		Name:        "test.name.2",
		Description: "test.description.2",
		Value:       "test.value.2",
	}

	repo.Create(context.Background(), prop1)
	repo.Create(context.Background(), prop2)

	var readProps []*model.Property
	readProps = make([]*model.Property, 0, 2)

	prop, _ := repo.FindByName(context.Background(), prop1.Name)
	readProps = append(readProps, prop)

	prop, _ = repo.FindByName(context.Background(), prop2.Name)
	readProps = append(readProps, prop)

	assert.Equal(t, true, (readProps[0].ID != ""))
	assert.Equal(t, prop1.Name, readProps[0].Name)
	assert.Equal(t, prop1.Description, readProps[0].Description)
	assert.Equal(t, prop1.Value, readProps[0].Value)
	assert.Equal(t, true, (readProps[1].ID != ""))
	assert.Equal(t, prop2.Name, readProps[1].Name)
	assert.Equal(t, prop2.Description, readProps[1].Description)
	assert.Equal(t, prop2.Value, readProps[1].Value)
}

func TestFindByNameNotFound(t *testing.T) {
	repo := setup()
	defer tearDown(repo)

	name := "test.notfound.name"
	_, err := repo.FindByName(context.Background(), name)

	assert.Equal(t, errors.NewEntityNotFound(model.Property{}, name), err)
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

	prop1 := &model.Property{
		Name:        "test.name.1",
		Description: "test.description.1",
		Value:       "test.value.1",
	}

	repo.Create(context.Background(), prop1)

	readProps, _ := repo.ReadAll(context.Background())

	assert.Equal(t, true, (readProps[0].ID != ""))
	assert.Equal(t, prop1.Name, readProps[0].Name)
	assert.Equal(t, prop1.Description, readProps[0].Description)
	assert.Equal(t, prop1.Value, readProps[0].Value)

	repo.Delete(context.Background(), readProps[0].ID)
	id := readProps[0].ID
	_, err := repo.FindByID(context.Background(), id)

	assert.Equal(t, errors.NewEntityNotFound(model.Property{}, id), err)
}

func TestDeleteNotFound(t *testing.T) {
	repo := setup()
	defer tearDown(repo)

	id := "test.notfound.id"
	err := repo.Delete(context.Background(), id)
	assert.Equal(t, errors.NewEntityNotFound(model.Property{}, id), err)
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

	prop1 := &model.Property{
		Name:        "test.name.1",
		Description: "test.description.1",
		Value:       "test.value.1",
	}

	repo.Create(context.Background(), prop1)

	readProps, _ := repo.ReadAll(context.Background())
	assert.Equal(t, true, (readProps[0].ID != ""))
	assert.Equal(t, prop1.Name, readProps[0].Name)
	assert.Equal(t, prop1.Description, readProps[0].Description)
	assert.Equal(t, prop1.Value, readProps[0].Value)

	prop1.Description = "test.description.1.2"

	repo.Update(context.Background(), prop1)

	found, _ := repo.FindByID(context.Background(), prop1.ID)
	assert.Equal(t, found.ID, prop1.ID)
	assert.Equal(t, found.Name, prop1.Name)
	assert.Equal(t, found.Description, prop1.Description)
	assert.Equal(t, found.Value, prop1.Value)
}

func TestUpdateNotFound(t *testing.T) {
	repo := setup()
	defer tearDown(repo)

	prop1 := &model.Property{
		Name:        "test.name.1",
		Description: "test.description.1",
		Value:       "test.value.1",
	}

	repo.Create(context.Background(), prop1)

	readProps, _ := repo.ReadAll(context.Background())
	assert.Equal(t, true, (readProps[0].ID != ""))
	assert.Equal(t, prop1.Name, readProps[0].Name)
	assert.Equal(t, prop1.Description, readProps[0].Description)
	assert.Equal(t, prop1.Value, readProps[0].Value)

	prop1.Description = "test.description.1.2"
	repo.Delete(context.Background(), prop1.ID)
	err := repo.Update(context.Background(), prop1)
	assert.Equal(t, errors.NewEntityNotFound(model.Property{}, prop1.ID), err)
}

func TestUpdateUnExpected(t *testing.T) {
	repo := setup()

	defer tearDown(repo)

	prop1 := &model.Property{
		Name:        "test.name.1",
		Description: "test.description.1",
		Value:       "test.value.1",
	}

	repo.Create(context.Background(), prop1)

	readProps, _ := repo.ReadAll(context.Background())
	assert.Equal(t, true, (readProps[0].ID != ""))
	assert.Equal(t, prop1.Name, readProps[0].Name)
	assert.Equal(t, prop1.Description, readProps[0].Description)
	assert.Equal(t, prop1.Value, readProps[0].Value)

	prop1.Description = "test.description.1.2"
	repo.db.Close()

	err := repo.Update(context.Background(), prop1)
	assert.Equal(t, g_errors.New("database not open"), err)
}

func BenchmarkReadAll(b *testing.B) {
	repo := setup()
	defer tearDown(repo)

	prop1 := &model.Property{
		Name:        "test.name.1",
		Description: "test.description.1",
		Value:       "test.value.1",
	}
	prop2 := &model.Property{
		Name:        "test.name.2",
		Description: "test.description.2",
		Value:       "test.value.2",
	}

	repo.Create(context.Background(), prop1)
	repo.Create(context.Background(), prop2)

	for n := 0; n < b.N; n++ {
		repo.ReadAll(context.Background())
	}

}

func setup() *PropertyRepository {
	util.CreateParentFolder(defaultDB)

	db, _ := storm.Open(defaultDB)

	return &PropertyRepository{
		db: db,
	}
}

func tearDown(repo *PropertyRepository) {
	repo.db.Close()

	os.Remove(defaultDB)
	os.Remove(defaultDir)
}
