package product

import (
	"context"
	"ecommerce/models"

	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()

type Repository interface {
	SetProduct(key string, value models.Product) error
	GetProduct(key string) (models.Product, error)
}

type repository struct {
	Client redis.Cmdable
}

func NewRedisRepository(c redis.Cmdable) Repository {
	return &repository{Client: c}
}

func (r *repository) SetProduct(key string, value models.Product) error {

	return r.Client.HSet(ctx, key, value).Err()

}

func (r *repository) GetProduct(key string) (models.Product, error) {
	product := models.Product{}
	err := r.Client.HGetAll(ctx, key).Scan(&product)

	if err != nil {
		return product, err
	} else {
		return product, nil
	}

}
