package repository

import (
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
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
	GetOrderRange(form form.GetOrder) ([]model.OrderSummary, int, error)
	GetOrderById(id string) (*model.OrderDetail, int, error)
	RemoveOrderById(id string) (*model.OrderDetail, int, error)
	GetOrderItemByOrderId(orderId string) ([]model.OrderItemDetail, int, error)
	GetTotalOrderId(orderId string) (float64, error)
	GetPaymentByOrderId(orderId string) (*model.Payment, int, error)
	RemovePaymentByOrderId(orderId string) (*model.Payment, int, error)
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
		Message:     form.Message,
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
		Total:       form.Total,
		Change:      form.Change,
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

func (entity orderEntity) GetOrderRange(form form.GetOrder) ([]model.OrderSummary, int, error) {
	logrus.Info("GetOrderRange")
	ctx, cancel := core.InitContext()
	defer cancel()
	var items []model.OrderSummary

	cursor, err := entity.orderRepo.Find(ctx, bson.M{"createdDate": bson.M{
		"$gt": form.StartDate,
		"$lt": form.EndDate,
	},
	})
	if err != nil {
		logrus.Error(err)
		return nil, http.StatusBadRequest, err
	}
	for cursor.Next(ctx) {
		var data model.OrderSummary
		err = cursor.Decode(&data)
		if err != nil {
			logrus.Error(err)
		}
		data.Total, _ = entity.GetTotalOrderId(data.Id.Hex())
		items = append(items, data)
	}
	if items == nil {
		items = []model.OrderSummary{}
	}
	return items, http.StatusOK, nil
}

func (entity orderEntity) GetOrderById(id string) (*model.OrderDetail, int, error) {
	logrus.Info("GetOrderById")
	ctx, cancel := core.InitContext()
	defer cancel()
	var data model.OrderDetail
	objId, _ := primitive.ObjectIDFromHex(id)
	err := entity.orderRepo.FindOne(ctx, bson.M{"_id": objId}).Decode(&data)
	if err != nil {
		logrus.Error(err)
		return nil, http.StatusBadRequest, err
	}
	data.Total, _ = entity.GetTotalOrderId(id)
	payment, _, err := entity.GetPaymentByOrderId(id)
	if err != nil {
		logrus.Error(err)
		return nil, http.StatusBadRequest, err
	}
	data.Payment = *payment
	items, _, err := entity.GetOrderItemByOrderId(id)
	if err != nil {
		logrus.Error(err)
		return nil, http.StatusBadRequest, err
	}
	data.Items = items
	return &data, http.StatusOK, nil
}

func (entity orderEntity) GetOrderItemByOrderId(orderId string) ([]model.OrderItemDetail, int, error) {
	logrus.Info("GetOrderItemByOrderId")
	ctx, cancel := core.InitContext()
	defer cancel()
	objId, _ := primitive.ObjectIDFromHex(orderId)
	cursor, err := entity.orderItemRepo.Aggregate(ctx, []bson.M{
		{
			"$match": bson.M{
				"orderId": objId,
			},
		},
		{
			"$lookup": bson.M{
				"from":         "products",
				"localField":   "productId",
				"foreignField": "_id",
				"as":           "product",
			},
		},
		{"$unwind": "$product"},
	})
	if err != nil {
		logrus.Error(err)
		return nil, http.StatusBadRequest, err
	}
	var items []model.OrderItemDetail
	for cursor.Next(ctx) {
		var data model.OrderItemDetail
		err = cursor.Decode(&data)
		if err != nil {
			logrus.Error(err)
		}
		items = append(items, data)
	}
	if items == nil {
		items = []model.OrderItemDetail{}
	}
	return items, http.StatusOK, nil
}

func (entity orderEntity) RemoveOrderItemByOrderId(orderId string) ([]model.OrderItemDetail, int, error) {
	logrus.Info("RemoveOrderItemByOrderId")
	ctx, cancel := core.InitContext()
	defer cancel()
	objId, _ := primitive.ObjectIDFromHex(orderId)
	items, _, err := entity.GetOrderItemByOrderId(orderId)
	if err != nil {
		logrus.Error(err)
		return nil, http.StatusBadRequest, err
	}
	_, err = entity.orderItemRepo.DeleteMany(ctx, bson.M{"orderId": objId})
	if err != nil {
		logrus.Error(err)
		return nil, http.StatusBadRequest, err
	}
	return items, http.StatusOK, nil
}

func (entity orderEntity) GetTotalOrderId(orderId string) (float64, error) {
	ctx, cancel := core.InitContext()
	objId, _ := primitive.ObjectIDFromHex(orderId)
	defer cancel()
	pipeline := []bson.M{
		{
			"$match": bson.M{
				"orderId": objId,
			},
		},
		{
			"$group": bson.M{
				"_id":   "",
				"total": bson.M{"$sum": "$price"},
			},
		},
	}
	var result []bson.M
	cursor, err := entity.orderItemRepo.Aggregate(ctx, pipeline)
	if err != nil {
		logrus.Error(err)
		return 0, err
	}
	err = cursor.All(ctx, &result)
	if result == nil {
		return 0, err
	}
	return result[0]["total"].(float64), nil
}

func (entity orderEntity) GetPaymentByOrderId(orderId string) (*model.Payment, int, error) {
	logrus.Info("GetPaymentByOrderId")
	ctx, cancel := core.InitContext()
	defer cancel()
	var data model.Payment
	objId, _ := primitive.ObjectIDFromHex(orderId)
	err := entity.paymentRepo.FindOne(ctx, bson.M{"orderId": objId}).Decode(&data)
	if err != nil {
		logrus.Error(err)
		return nil, http.StatusBadRequest, err
	}
	return &data, http.StatusOK, nil
}

func (entity orderEntity) RemovePaymentByOrderId(orderId string) (*model.Payment, int, error) {
	logrus.Info("RemovePaymentByOrderId")
	ctx, cancel := core.InitContext()
	defer cancel()
	var data model.Payment
	objId, _ := primitive.ObjectIDFromHex(orderId)
	err := entity.paymentRepo.FindOne(ctx, bson.M{"orderId": objId}).Decode(&data)
	if err != nil {
		logrus.Error(err)
		return nil, http.StatusBadRequest, err
	}
	_, err = entity.paymentRepo.DeleteMany(ctx, bson.M{"orderId": objId})
	if err != nil {
		logrus.Error(err)
		return nil, http.StatusBadRequest, err
	}
	return &data, http.StatusOK, nil
}

func (entity orderEntity) RemoveOrderById(id string) (*model.OrderDetail, int, error) {
	logrus.Info("RemoveOrderById")
	ctx, cancel := core.InitContext()
	defer cancel()
	var data model.OrderDetail
	objId, _ := primitive.ObjectIDFromHex(id)
	err := entity.orderRepo.FindOne(ctx, bson.M{"_id": objId}).Decode(&data)
	if err != nil {
		logrus.Error(err)
		return nil, http.StatusBadRequest, err
	}
	_, err = entity.orderRepo.DeleteOne(ctx, bson.M{"_id": data.Id})
	if err != nil {
		logrus.Error(err)
		return nil, http.StatusBadRequest, err
	}
	data.Total, _ = entity.GetTotalOrderId(id)
	payment, _, _ := entity.RemovePaymentByOrderId(id)
	data.Payment = *payment
	items, _, _ := entity.RemoveOrderItemByOrderId(id)
	data.Items = items
	return &data, http.StatusOK, nil
}
