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
	GetProductAll() ([]model.Product, int, error)
	GetProductBySerialNumber(serialNumber string) (*model.Product, int, error)
	GetProductById(id string) (*model.Product, int, error)
	CreateProduct(form form.Product) (*model.Product, int, error)
	RemoveProductById(id string) (*model.Product, int, error)
	UpdateProductById(id string, form form.UpdateProduct) (*model.Product, int, error)
	RemoveQuantityById(id string, quantity int) (*model.Product, int, error)
	AddQuantityById(id string, quantity int) (*model.Product, int, error)
	GetTotalCostPrice(id string, quantity int) float64

	CreateLot(productId string, form form.Product) (*model.ProductLot, int, error)
	GetLotAllByProductId(productId string) ([]model.ProductLot, int, error)
	GetLotById(id string) (*model.ProductLot, int, error)
	UpdateLotById(id string, form form.ProductLot) (*model.ProductLot, int, error)
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

func (entity *productEntity) GetProductAll() ([]model.Product, int, error) {
	logrus.Info("GetProductAll")
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

func (entity *productEntity) GetProductBySerialNumber(serialNumber string) (*model.Product, int, error) {
	logrus.Info("GetProductBySerialNumber")
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

func (entity *productEntity) CreateProduct(form form.Product) (*model.Product, int, error) {
	logrus.Info("CreateProduct")
	ctx, cancel := core.InitContext()
	defer cancel()
	data, _, _ := entity.GetProductBySerialNumber(form.SerialNumber)
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
		_, _, err = entity.CreateLot(data.Id.Hex(), form)
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
		_, _, err = entity.CreateLot(data.Id.Hex(), form)
		if err != nil {
			logrus.Error(err)
			return nil, http.StatusBadRequest, err
		}
		return &data, http.StatusOK, nil
	}
}

func (entity *productEntity) GetProductById(id string) (*model.Product, int, error) {
	logrus.Info("GetProductById")
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

func (entity *productEntity) RemoveProductById(id string) (*model.Product, int, error) {
	logrus.Info("RemoveProductById")
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
	_, _ = entity.lotRepo.DeleteMany(ctx, bson.M{"productId": objId})
	return &data, http.StatusOK, nil
}

func (entity *productEntity) UpdateProductById(id string, form form.UpdateProduct) (*model.Product, int, error) {
	logrus.Info("UpdateProductById")
	ctx, cancel := core.InitContext()
	defer cancel()
	objId, _ := primitive.ObjectIDFromHex(id)
	data, _, err := entity.GetProductById(id)
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

	data, _, err := entity.GetProductById(id)
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

func (entity *productEntity) AddQuantityById(id string, quantity int) (*model.Product, int, error) {
	logrus.Info("AddQuantityById")
	ctx, cancel := core.InitContext()
	defer cancel()
	objId, _ := primitive.ObjectIDFromHex(id)
	data, _, err := entity.GetProductById(id)
	if err != nil {
		logrus.Error(err)
		return nil, http.StatusNotFound, err
	}
	data.Quantity = data.Quantity + quantity
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

func (entity *productEntity) GetTotalCostPrice(id string, quantity int) float64 {
	logrus.Info("GetTotalCostPrice")
	data, _, err := entity.GetProductById(id)
	if err != nil {
		logrus.Error(err)
		return 0
	}
	return data.CostPrice * float64(quantity)
}

func (entity *productEntity) CreateLot(productId string, form form.Product) (*model.ProductLot, int, error) {
	logrus.Info("CreateLot")
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

func (entity *productEntity) GetLotById(id string) (*model.ProductLot, int, error) {
	logrus.Info("GetLotById")
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

func (entity *productEntity) UpdateLotById(id string, form form.ProductLot) (*model.ProductLot, int, error) {
	logrus.Info("UpdateLotById")
	ctx, cancel := core.InitContext()
	defer cancel()
	objId, _ := primitive.ObjectIDFromHex(id)
	data, _, err := entity.GetLotById(id)
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
