package category

import "ecommerce/category/model"

type categoryService struct {
	categoryService model.CategoryRepository
}

func NewCategoryService(s model.CategoryRepository) *categoryService {
	return &categoryService{
		categoryService: s,
	}
}

func (cs *categoryService) Create(input model.Category) error {

	cat := model.Category{}

	cat.Name = input.Name
	cat.ParentID = input.ParentID

	err := cs.categoryService.InsertCategory(cat)

	if err != nil {
		return err
	} else {
		return nil
	}
}

func (cs *categoryService) Read(categoryID int) (model.Category, error) {

	category := model.Category{}

	category, err := cs.categoryService.GetCategory(category, categoryID)

	if err != nil {
		return category, err
	} else {
		return category, nil
	}

}

func (cs *categoryService) Update(categoryID int, input model.Category) error {

	cat := model.Category{}
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

	category := model.Category{}
	err := cs.categoryService.DeleteRow(category, categoryID)

	if err != nil {
		return err
	} else {
		return nil
	}
}
