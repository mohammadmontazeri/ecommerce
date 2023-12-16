package category

import "ecommerce/models"

type categoryService struct {
	categoryService models.CategoryRepository
}

func NewCategoryService(s models.CategoryRepository) *categoryService {
	return &categoryService{
		categoryService: s,
	}
}

func (cs *categoryService) Create(input models.Category) error {

	cat := models.Category{}

	cat.Name = input.Name
	cat.ParentID = input.ParentID

	err := cs.categoryService.InsertCategory(cat)

	if err != nil {
		return err
	} else {
		return nil
	}
}

func (cs *categoryService) Read(categoryID int) (models.Category, error) {

	category := models.Category{}

	category, err := cs.categoryService.GetCategory(category, categoryID)

	if err != nil {
		return category, err
	} else {
		return category, nil
	}

}

func (cs *categoryService) Update(categoryID int, input models.Category) error {

	cat := models.Category{}
	cat.Name = input.Name
	cat.ParentID = input.ParentID

	err := cs.categoryService.UpdateRow(cat, categoryID)

	if err != nil {
		return err
	} else {
		return nil
	}

}

func (cs *categoryService) Delete(categoryID int) error {

	category := models.Category{}
	err := cs.categoryService.DeleteRow(category, categoryID)

	if err != nil {
		return err
	} else {
		return nil
	}
}
