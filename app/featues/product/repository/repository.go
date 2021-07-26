package repository

import (
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"mgo-gin/app/core"
	"mgo-gin/app/featues/product/form"
	"mgo-gin/app/featues/product/model"
	"mgo-gin/db"
	"net/http"
	"time"
)

var ProductEntity IProduct

type productEntity struct {
	resource    *db.Resource
	productRepo *mongo.Collection
	lotRepo     *mongo.Collection
}

type IProduct interface {
	CreateIndex() (string, error)
	GetAll() ([]model.Product, int, error)
	GetOneBySerialNumber(serialNumber string) (*model.Product, int, error)
	GetOneById(id string) (*model.Product, int, error)
	CreateOne(form form.Product) (*model.Product, int, error)
	RemoveOneById(id string) (*model.Product, int, error)
	UpdateOneById(id string, form form.UpdateProduct) (*model.Product, int, error)
	RemoveQuantityById(id string, quantity int) (*model.Product, int, error)

	CreateLotOne(productId string, form form.Product) (*model.ProductLot, int, error)
	GetLotAllByProductId(productId string) ([]model.ProductLot, int, error)
	GetLotOneById(id string) (*model.ProductLot, int, error)
	UpdateLotOneById(id string, form form.ProductLot) (*model.ProductLot, int, error)
}

func NewProductEntity(resource *db.Resource) IProduct {
	productRepo := resource.DB.Collection("products")
	lotRepo := resource.DB.Collection("product_lots")
	ProductEntity = &productEntity{resource: resource, productRepo: productRepo, lotRepo: lotRepo}
	return ProductEntity
}

func (entity *productEntity) CreateIndex() (string, error) {
	ctx, cancel := core.InitContext()
	defer cancel()
	mod := mongo.IndexModel{
		Keys: bson.M{
			"serialNumber": 1,
		},
		Options: options.Index().SetUnique(true),
	}
	ind, err := entity.productRepo.Indexes().CreateOne(ctx, mod)
	return ind, err
}

func (entity *productEntity) GetAll() ([]model.Product, int, error) {
	logrus.Info("GetAll")
	ctx, cancel := core.InitContext()
	defer cancel()
	var products []model.Product
	cursor, err := entity.productRepo.Find(ctx, bson.M{})
	if err != nil {
		logrus.Error(err)
		return nil, http.StatusBadRequest, err
	}
	for cursor.Next(ctx) {
		var user model.Product
		err = cursor.Decode(&user)
		if err != nil {
			logrus.Error(err)
		}
		products = append(products, user)
	}
	if products == nil {
		products = []model.Product{}
	}
	return products, http.StatusOK, nil
}

func (entity *productEntity) GetOneBySerialNumber(serialNumber string) (*model.Product, int, error) {
	logrus.Info("GetOneBySerialNumber")
	ctx, cancel := core.InitContext()
	defer cancel()
	var data model.Product
	err := entity.productRepo.FindOne(ctx, bson.M{"serialNumber": serialNumber}).Decode(&data)
	if err != nil {
		logrus.Error(err)
		return nil, http.StatusBadRequest, err
	}
	return &data, http.StatusOK, nil
}

func (entity *productEntity) CreateOne(form form.Product) (*model.Product, int, error) {
	logrus.Info("CreateOne")
	ctx, cancel := core.InitContext()
	defer cancel()
	data, _, _ := entity.GetOneBySerialNumber(form.SerialNumber)
	if data != nil {
		data.Name = form.Name
		data.NameEn = form.NameEn
		data.Description = form.Description
		data.SerialNumber = form.SerialNumber
		data.Price = form.Price
		data.CostPrice = form.CostPrice
		data.Unit = form.Unit
		data.Quantity = data.Quantity + form.Quantity
		data.UpdatedDate = time.Now()

		isReturnNewDoc := options.After
		opts := &options.FindOneAndUpdateOptions{
			ReturnDocument: &isReturnNewDoc,
		}
		err := entity.productRepo.FindOneAndUpdate(ctx, bson.M{"_id": data.Id}, bson.M{"$set": data}, opts).Decode(&data)
		if err != nil {
			logrus.Error(err)
			return nil, http.StatusBadRequest, err
		}
		_, _, err = entity.CreateLotOne(data.Id.Hex(), form)
		if err != nil {
			logrus.Error(err)
			return nil, http.StatusBadRequest, err
		}
		return data, http.StatusOK, nil
	} else {
		data := model.Product{}
		data.Id = primitive.NewObjectID()
		data.Name = form.Name
		data.NameEn = form.NameEn
		data.Description = form.Description
		data.SerialNumber = form.SerialNumber
		data.Unit = form.Unit
		data.Price = form.Price
		data.CostPrice = form.CostPrice
		data.Quantity = form.Quantity
		data.CreatedDate = time.Now()
		data.UpdatedDate = time.Now()
		_, err := entity.productRepo.InsertOne(ctx, data)
		if err != nil {
			logrus.Error(err)
			return nil, http.StatusBadRequest, err
		}
		_, _, err = entity.CreateLotOne(data.Id.Hex(), form)
		if err != nil {
			logrus.Error(err)
			return nil, http.StatusBadRequest, err
		}
		return &data, http.StatusOK, nil
	}
}

func (entity *productEntity) GetOneById(id string) (*model.Product, int, error) {
	logrus.Info("GetOneById")
	ctx, cancel := core.InitContext()
	defer cancel()
	var data model.Product
	objId, _ := primitive.ObjectIDFromHex(id)
	err := entity.productRepo.FindOne(ctx, bson.M{"_id": objId}).Decode(&data)
	if err != nil {
		logrus.Error(err)
		return nil, http.StatusBadRequest, err
	}
	return &data, http.StatusOK, nil
}

func (entity *productEntity) RemoveOneById(id string) (*model.Product, int, error) {
	logrus.Info("RemoveOneById")
	ctx, cancel := core.InitContext()
	defer cancel()
	var data model.Product
	objId, _ := primitive.ObjectIDFromHex(id)
	err := entity.productRepo.FindOne(ctx, bson.M{"_id": objId}).Decode(&data)
	if err != nil {
		logrus.Error(err)
		return nil, http.StatusBadRequest, err
	}
	_, err = entity.productRepo.DeleteOne(ctx, bson.M{"_id": objId})
	if err != nil {
		logrus.Error(err)
		return nil, http.StatusBadRequest, err
	}
	return &data, http.StatusOK, nil
}

func (entity *productEntity) UpdateOneById(id string, form form.UpdateProduct) (*model.Product, int, error) {
	logrus.Info("UpdateOneById")
	ctx, cancel := core.InitContext()
	defer cancel()
	objId, _ := primitive.ObjectIDFromHex(id)
	data, _, err := entity.GetOneById(id)
	if err != nil {
		logrus.Error(err)
		return nil, http.StatusNotFound, err
	}
	data.Name = form.Name
	data.NameEn = form.NameEn
	data.Description = form.Description
	data.Price = form.Price
	data.CostPrice = form.CostPrice
	data.Unit = form.Unit
	data.Quantity = form.Quantity
	data.UpdatedDate = time.Now()

	isReturnNewDoc := options.After
	opts := &options.FindOneAndUpdateOptions{
		ReturnDocument: &isReturnNewDoc,
	}
	err = entity.productRepo.FindOneAndUpdate(ctx, bson.M{"_id": objId}, bson.M{"$set": data}, opts).Decode(&data)
	if err != nil {
		logrus.Error(err)
		return nil, http.StatusBadRequest, err
	}
	return data, http.StatusOK, nil
}

func (entity *productEntity) RemoveQuantityById(id string, quantity int) (*model.Product, int, error) {
	logrus.Info("RemoveQuantityById")
	ctx, cancel := core.InitContext()
	defer cancel()
	objId, _ := primitive.ObjectIDFromHex(id)
	data, _, err := entity.GetOneById(id)
	if err != nil {
		logrus.Error(err)
		return nil, http.StatusNotFound, err
	}
	data.Quantity = data.Quantity - quantity
	data.UpdatedDate = time.Now()
	isReturnNewDoc := options.After
	opts := &options.FindOneAndUpdateOptions{
		ReturnDocument: &isReturnNewDoc,
	}
	err = entity.productRepo.FindOneAndUpdate(ctx, bson.M{"_id": objId}, bson.M{"$set": data}, opts).Decode(&data)
	if err != nil {
		logrus.Error(err)
		return nil, http.StatusBadRequest, err
	}
	return data, http.StatusOK, nil
}

func (entity *productEntity) CreateLotOne(productId string, form form.Product) (*model.ProductLot, int, error) {
	logrus.Info("CreateLotOne")
	ctx, cancel := core.InitContext()
	defer cancel()
	data := model.ProductLot{}
	data.Id = primitive.NewObjectID()
	data.ProductId, _ = primitive.ObjectIDFromHex(productId)
	data.LotNumber = form.LotNumber
	data.ExpireDate = form.ExpireDate
	data.Quantity = form.Quantity
	data.CostPrice = form.CostPrice
	data.CreatedDate = time.Now()
	data.UpdatedDate = time.Now()
	_, err := entity.lotRepo.InsertOne(ctx, data)
	if err != nil {
		logrus.Error(err)
		return nil, http.StatusBadRequest, err
	}
	return &data, http.StatusOK, nil
}

func (entity *productEntity) GetLotAllByProductId(productId string) ([]model.ProductLot, int, error) {
	logrus.Info("GetLotAllByProductId")
	var productLots []model.ProductLot
	ctx, cancel := core.InitContext()
	defer cancel()
	objId, _ := primitive.ObjectIDFromHex(productId)
	cursor, err := entity.lotRepo.Find(ctx, bson.M{"productId": objId})
	if err != nil {
		logrus.Error(err)
		return nil, http.StatusBadRequest, err
	}
	for cursor.Next(ctx) {
		var productLot model.ProductLot
		err = cursor.Decode(&productLot)
		if err != nil {
			logrus.Error(err)
		}
		productLots = append(productLots, productLot)
	}
	if productLots == nil {
		productLots = []model.ProductLot{}
	}
	return productLots, http.StatusOK, nil
}

func (entity *productEntity) GetLotOneById(id string) (*model.ProductLot, int, error) {
	logrus.Info("GetLotOneById")
	ctx, cancel := core.InitContext()
	defer cancel()
	var data model.ProductLot
	objId, _ := primitive.ObjectIDFromHex(id)
	err := entity.lotRepo.FindOne(ctx, bson.M{"_id": objId}).Decode(&data)
	if err != nil {
		logrus.Error(err)
		return nil, http.StatusBadRequest, err
	}
	return &data, http.StatusOK, nil
}

func (entity *productEntity) UpdateLotOneById(id string, form form.ProductLot) (*model.ProductLot, int, error) {
	logrus.Info("UpdateLotOneById")
	ctx, cancel := core.InitContext()
	defer cancel()
	objId, _ := primitive.ObjectIDFromHex(id)
	data, _, err := entity.GetLotOneById(id)
	if err != nil {
		logrus.Error(err)
		return nil, http.StatusNotFound, err
	}
	data.LotNumber = form.LotNumber
	data.ExpireDate = form.ExpireDate
	data.Quantity = form.Quantity
	data.CostPrice = form.CostPrice
	data.UpdatedDate = time.Now()

	isReturnNewDoc := options.After
	opts := &options.FindOneAndUpdateOptions{
		ReturnDocument: &isReturnNewDoc,
	}
	err = entity.lotRepo.FindOneAndUpdate(ctx, bson.M{"_id": objId}, bson.M{"$set": data}, opts).Decode(&data)
	if err != nil {
		logrus.Error(err)
		return nil, http.StatusBadRequest, err
	}
	return data, http.StatusOK, nil
}
