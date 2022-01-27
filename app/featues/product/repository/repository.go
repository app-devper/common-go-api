package repository

import (
	"devper/app/core"
	"devper/app/featues/product/form"
	"devper/app/featues/product/model"
	"devper/db"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"strings"
	"time"
)

var Entity IProduct

type productEntity struct {
	resource    *db.Resource
	productRepo *mongo.Collection
	lotRepo     *mongo.Collection
}

type IProduct interface {
	CreateIndex() (string, error)
	GetProductAll() ([]model.Product, error)
	GetProductBySerialNumber(serialNumber string) (*model.Product, error)
	GetProductById(id string) (*model.Product, error)
	CreateProduct(form form.Product) (*model.Product, error)
	RemoveProductById(id string) (*model.Product, error)
	UpdateProductById(id string, form form.UpdateProduct) (*model.Product, error)
	RemoveQuantityById(id string, quantity int) (*model.Product, error)
	AddQuantityById(id string, quantity int) (*model.Product, error)
	GetTotalCostPrice(id string, quantity int) float64

	CreateLot(productId string, form form.Product) (*model.ProductLot, error)
	GetLotAllByProductId(productId string) ([]model.ProductLot, error)
	GetLotById(id string) (*model.ProductLot, error)
	UpdateLotById(id string, form form.ProductLot) (*model.ProductLot, error)
}

func NewProductEntity(resource *db.Resource) IProduct {
	productRepo := resource.DB.Collection("products")
	lotRepo := resource.DB.Collection("product_lots")
	Entity = &productEntity{resource: resource, productRepo: productRepo, lotRepo: lotRepo}
	_, _ = Entity.CreateIndex()
	return Entity
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

func (entity *productEntity) GetProductAll() ([]model.Product, error) {
	logrus.Info("GetProductAll")
	ctx, cancel := core.InitContext()
	defer cancel()
	var products []model.Product
	cursor, err := entity.productRepo.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	for cursor.Next(ctx) {
		var user model.Product
		err = cursor.Decode(&user)
		if err != nil {
			logrus.Error(err)
			logrus.Info(cursor.Current)
		} else {
			products = append(products, user)
		}
	}
	if products == nil {
		products = []model.Product{}
	}
	return products, nil
}

func (entity *productEntity) GetProductBySerialNumber(serialNumber string) (*model.Product, error) {
	logrus.Info("GetProductBySerialNumber")
	ctx, cancel := core.InitContext()
	defer cancel()
	var data model.Product
	err := entity.productRepo.FindOne(ctx, bson.M{"serialNumber": serialNumber}).Decode(&data)
	if err != nil {
		return nil, err
	}
	return &data, nil
}

func (entity *productEntity) CreateProduct(form form.Product) (*model.Product, error) {
	logrus.Info("CreateProduct")
	ctx, cancel := core.InitContext()
	defer cancel()
	serialNumber := strings.TrimSpace(form.SerialNumber)
	data, _ := entity.GetProductBySerialNumber(serialNumber)
	if data != nil {
		data.Name = form.Name
		data.NameEn = form.NameEn
		data.Description = form.Description
		data.SerialNumber = serialNumber
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
			return nil, err
		}
		_, err = entity.CreateLot(data.Id.Hex(), form)
		if err != nil {
			return nil, err
		}
		return data, nil
	} else {
		data := model.Product{}
		data.Id = primitive.NewObjectID()
		data.Name = form.Name
		data.NameEn = form.NameEn
		data.Description = form.Description
		data.SerialNumber = serialNumber
		data.Unit = form.Unit
		data.Price = form.Price
		data.CostPrice = form.CostPrice
		data.Quantity = form.Quantity
		data.CreatedDate = time.Now()
		data.UpdatedDate = time.Now()
		_, err := entity.productRepo.InsertOne(ctx, data)
		if err != nil {
			return nil, err
		}
		_, err = entity.CreateLot(data.Id.Hex(), form)
		if err != nil {
			return nil, err
		}
		return &data, nil
	}
}

func (entity *productEntity) GetProductById(id string) (*model.Product, error) {
	logrus.Info("GetProductById")
	ctx, cancel := core.InitContext()
	defer cancel()
	var data model.Product
	objId, _ := primitive.ObjectIDFromHex(id)
	err := entity.productRepo.FindOne(ctx, bson.M{"_id": objId}).Decode(&data)
	if err != nil {
		return nil, err
	}
	return &data, nil
}

func (entity *productEntity) RemoveProductById(id string) (*model.Product, error) {
	logrus.Info("RemoveProductById")
	ctx, cancel := core.InitContext()
	defer cancel()
	var data model.Product
	objId, _ := primitive.ObjectIDFromHex(id)
	err := entity.productRepo.FindOne(ctx, bson.M{"_id": objId}).Decode(&data)
	if err != nil {
		return nil, err
	}
	_, err = entity.productRepo.DeleteOne(ctx, bson.M{"_id": objId})
	if err != nil {
		return nil, err
	}
	_, _ = entity.lotRepo.DeleteMany(ctx, bson.M{"productId": objId})
	return &data, nil
}

func (entity *productEntity) UpdateProductById(id string, form form.UpdateProduct) (*model.Product, error) {
	logrus.Info("UpdateProductById")
	ctx, cancel := core.InitContext()
	defer cancel()
	objId, _ := primitive.ObjectIDFromHex(id)
	data, err := entity.GetProductById(id)
	if err != nil {
		return nil, err
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
		return nil, err
	}
	return data, nil
}

func (entity *productEntity) RemoveQuantityById(id string, quantity int) (*model.Product, error) {
	logrus.Info("RemoveQuantityById")
	ctx, cancel := core.InitContext()
	defer cancel()
	objId, _ := primitive.ObjectIDFromHex(id)

	data, err := entity.GetProductById(id)
	if err != nil {
		return nil, err
	}
	data.Quantity = data.Quantity - quantity
	data.UpdatedDate = time.Now()

	isReturnNewDoc := options.After
	opts := &options.FindOneAndUpdateOptions{
		ReturnDocument: &isReturnNewDoc,
	}
	err = entity.productRepo.FindOneAndUpdate(ctx, bson.M{"_id": objId}, bson.M{"$set": data}, opts).Decode(&data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (entity *productEntity) AddQuantityById(id string, quantity int) (*model.Product, error) {
	logrus.Info("AddQuantityById")
	ctx, cancel := core.InitContext()
	defer cancel()
	objId, _ := primitive.ObjectIDFromHex(id)
	data, err := entity.GetProductById(id)
	if err != nil {
		return nil, err
	}
	data.Quantity = data.Quantity + quantity
	data.UpdatedDate = time.Now()
	isReturnNewDoc := options.After
	opts := &options.FindOneAndUpdateOptions{
		ReturnDocument: &isReturnNewDoc,
	}
	err = entity.productRepo.FindOneAndUpdate(ctx, bson.M{"_id": objId}, bson.M{"$set": data}, opts).Decode(&data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (entity *productEntity) GetTotalCostPrice(id string, quantity int) float64 {
	logrus.Info("GetTotalCostPrice")
	data, err := entity.GetProductById(id)
	if err != nil {
		return 0
	}
	return data.CostPrice * float64(quantity)
}

func (entity *productEntity) CreateLot(productId string, form form.Product) (*model.ProductLot, error) {
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
		return nil, err
	}
	return &data, nil
}

func (entity *productEntity) GetLotAllByProductId(productId string) ([]model.ProductLot, error) {
	logrus.Info("GetLotAllByProductId")
	var productLots []model.ProductLot
	ctx, cancel := core.InitContext()
	defer cancel()
	objId, _ := primitive.ObjectIDFromHex(productId)
	cursor, err := entity.lotRepo.Find(ctx, bson.M{"productId": objId})
	if err != nil {
		return nil, err
	}
	for cursor.Next(ctx) {
		var productLot model.ProductLot
		err = cursor.Decode(&productLot)
		if err != nil {
			logrus.Error(err)
			logrus.Info(cursor.Current)
		} else {
			productLots = append(productLots, productLot)
		}
	}
	if productLots == nil {
		productLots = []model.ProductLot{}
	}
	return productLots, nil
}

func (entity *productEntity) GetLotById(id string) (*model.ProductLot, error) {
	logrus.Info("GetLotById")
	ctx, cancel := core.InitContext()
	defer cancel()
	var data model.ProductLot
	objId, _ := primitive.ObjectIDFromHex(id)
	err := entity.lotRepo.FindOne(ctx, bson.M{"_id": objId}).Decode(&data)
	if err != nil {
		return nil, err
	}
	return &data, nil
}

func (entity *productEntity) UpdateLotById(id string, form form.ProductLot) (*model.ProductLot, error) {
	logrus.Info("UpdateLotById")
	ctx, cancel := core.InitContext()
	defer cancel()
	objId, _ := primitive.ObjectIDFromHex(id)
	data, err := entity.GetLotById(id)
	if err != nil {
		return nil, err
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
		return nil, err
	}
	return data, nil
}
