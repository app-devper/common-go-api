package repository

import (
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"mgo-gin/app/core"
	"mgo-gin/app/featues/order/form"
	"mgo-gin/app/featues/order/model"
	"mgo-gin/db"
	"mgo-gin/utils/constant"
	"net/http"
	"time"
)

var OrderEntity IOrder

type orderEntity struct {
	resource      *db.Resource
	orderRepo     *mongo.Collection
	orderItemRepo *mongo.Collection
	paymentRepo   *mongo.Collection
}

type IOrder interface {
	CreateOrder(form form.Order) (*model.Order, int, error)
}

func NewOrderEntity(resource *db.Resource) IOrder {
	orderRepo := resource.DB.Collection("orders")
	orderItemRepo := resource.DB.Collection("order_items")
	paymentRepo := resource.DB.Collection("payments")
	OrderEntity = &orderEntity{resource: resource, orderRepo: orderRepo, orderItemRepo: orderItemRepo, paymentRepo: paymentRepo}
	return OrderEntity
}

func (entity orderEntity) CreateOrder(form form.Order) (*model.Order, int, error) {
	logrus.Info("CreateOrder")
	ctx, cancel := core.InitContext()
	defer cancel()
	var orderId = primitive.NewObjectID()
	data := model.Order{
		Id:          orderId,
		Status:      constant.ACTIVE,
		CreatedDate: time.Now(),
		UpdatedDate: time.Now(),
	}
	_, err := entity.orderRepo.InsertOne(ctx, data)
	if err != nil {
		logrus.Error(err)
		return nil, http.StatusBadRequest, err
	}
	count := len(form.Items)
	orderItem := make([]interface{}, count)
	for i := 0; i < count; i++ {
		formItem := form.Items[i]
		productId, _ := primitive.ObjectIDFromHex(formItem.ProductId)
		item := model.OrderItem{
			Id:          primitive.NewObjectID(),
			OrderId:     orderId,
			ProductId:   productId,
			Quantity:    formItem.Quantity,
			Price:       formItem.Price,
			Discount:    formItem.Discount,
			CreatedDate: time.Now(),
			UpdatedDate: time.Now(),
		}
		orderItem[i] = item
	}
	_, err = entity.orderItemRepo.InsertMany(ctx, orderItem)
	if err != nil {
		logrus.Error(err)
		return nil, http.StatusBadRequest, err
	}
	payment := model.Payment{
		Id:          primitive.NewObjectID(),
		OrderId:     orderId,
		Status:      constant.ACTIVE,
		Amount:      form.Amount,
		Type:        form.Type,
		CreatedDate: time.Now(),
		UpdatedDate: time.Now(),
	}
	_, err = entity.paymentRepo.InsertOne(ctx, payment)
	if err != nil {
		logrus.Error(err)
		return nil, http.StatusBadRequest, err
	}
	return &data, http.StatusOK, nil
}
