package repository

import (
	"context"
	"mirza/gogen/refactor/t201/p1/model/entity"
	"mirza/gogen/refactor/t201/p1/model/service"
)

type SaveOrderRepo interface {
	SaveOrder(ctx context.Context, order *entity.Order, svc service.PublishMessage) error
}
