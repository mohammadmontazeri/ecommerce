package category

import (
	categorymocks "ecommerce/internal/mocks/categorymocks"
	"ecommerce/models"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateCategoryService(t *testing.T) {

	repo := &categorymocks.CategoryRepository{}

	repo.On("InsertCategory", mock.AnythingOfType("models.Category")).
		Return(nil).
		Once()

	categoryService := NewCategoryService(repo)

	parentID := 1
	categoryInput := models.Category{
		Name:     "category_test_name",
		ParentID: &parentID,
	}

	err := categoryService.Create(categoryInput)

	assert.Nil(t, err)

}

func TestReadCategoryService(t *testing.T) {
	repo := &categorymocks.CategoryRepository{}

	repo.On("GetCategory", mock.AnythingOfType("models.Category"), mock.AnythingOfType("int")).
		Return(func(cat models.Category, id int) models.Category {
			return cat
		}, nil).
		Once()

	categoryService := NewCategoryService(repo)

	categoryID := 10
	_, err := categoryService.Read(categoryID)

	assert.Nil(t, err)

}

func TestUpdateCategoryService(t *testing.T) {
	repo := &categorymocks.CategoryRepository{}

	repo.On("UpdateRow", mock.AnythingOfType("models.Category"), mock.AnythingOfType("int")).
		Return(nil).
		Once()

	categoryService := NewCategoryService(repo)

	parentID := 1
	categoryInput := models.Category{
		Name:     "category_test_name",
		ParentID: &parentID,
	}
	categoryID := 10
	err := categoryService.Update(categoryID, categoryInput)

	assert.Nil(t, err)

}

func TestDeleteCategoryService(t *testing.T) {
	repo := &categorymocks.CategoryRepository{}

	repo.On("DeleteRow", mock.AnythingOfType("models.Category"), mock.AnythingOfType("int")).
		Return(nil).
		Once()

	categoryService := NewCategoryService(repo)

	categoryID := 10
	err := categoryService.Delete(categoryID)

	assert.Nil(t, err)

}
