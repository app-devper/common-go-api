package repository

import (
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"mgo-gin/app/core"
	"mgo-gin/app/featues/category/form"
	"mgo-gin/app/featues/category/model"
	"mgo-gin/db"
	"net/http"
	"strings"
	"time"
)

var CategoryEntity ICategory

type categoryEntity struct {
	resource     *db.Resource
	categoryRepo *mongo.Collection
}

func (entity categoryEntity) UpdateDefaultCategoryById(id string) (*model.Category, int, error) {
	logrus.Info("UpdateDefaultCategoryById")
	ctx, cancel := core.InitContext()
	defer cancel()
	_, err := entity.categoryRepo.UpdateMany(ctx, bson.M{}, bson.M{"$set": bson.M{
		"default": false,
	}})
	if err != nil {
		logrus.Error(err)
		return nil, http.StatusBadRequest, err
	}

	objId, _ := primitive.ObjectIDFromHex(id)
	var data model.Category
	err = entity.categoryRepo.FindOne(ctx, bson.M{"_id": objId}).Decode(&data)
	if err != nil {
		logrus.Error(err)
		return nil, http.StatusBadRequest, err
	}
	data.Default = true
	data.UpdatedDate = time.Now()

	isReturnNewDoc := options.After
	opts := &options.FindOneAndUpdateOptions{
		ReturnDocument: &isReturnNewDoc,
	}
	err = entity.categoryRepo.FindOneAndUpdate(ctx, bson.M{"_id": objId}, bson.M{"$set": data}, opts).Decode(&data)
	if err != nil {
		logrus.Error(err)
		return nil, http.StatusBadRequest, err
	}
	return &data, http.StatusOK, nil
}

func (entity categoryEntity) GetCategoryAll() ([]model.Category, int, error) {
	logrus.Info("GetCategoryAll")
	ctx, cancel := core.InitContext()
	defer cancel()
	var items []model.Category
	cursor, err := entity.categoryRepo.Find(ctx, bson.M{})
	if err != nil {
		logrus.Error(err)
		return nil, http.StatusBadRequest, err
	}
	for cursor.Next(ctx) {
		var user model.Category
		err = cursor.Decode(&user)
		if err != nil {
			logrus.Error(err)
		}
		items = append(items, user)
	}
	if items == nil {
		items = []model.Category{}
	}
	return items, http.StatusOK, nil
}

func (entity categoryEntity) CreateCategory(form form.Category) (*model.Category, int, error) {
	logrus.Info("CreateCategory")
	ctx, cancel := core.InitContext()
	defer cancel()
	data := model.Category{
		Id:          primitive.NewObjectID(),
		Name:        form.Name,
		Value:       strings.ToUpper(form.Value),
		Description: form.Description,
		CreatedDate: time.Now(),
		UpdatedDate: time.Now(),
	}
	_, err := entity.categoryRepo.InsertOne(ctx, data)
	if err != nil {
		logrus.Error(err)
		return nil, http.StatusBadRequest, err
	}
	return &data, http.StatusOK, nil
}

func (entity categoryEntity) GetCategoryById(id string) (*model.Category, int, error) {
	logrus.Info("GetCategoryById")
	ctx, cancel := core.InitContext()
	defer cancel()
	objId, _ := primitive.ObjectIDFromHex(id)
	var data model.Category
	err := entity.categoryRepo.FindOne(ctx, bson.M{"_id": objId}).Decode(&data)
	if err != nil {
		logrus.Error(err)
		return nil, http.StatusBadRequest, err
	}
	return &data, http.StatusOK, nil
}

func (entity categoryEntity) RemoveCategoryById(id string) (*model.Category, int, error) {
	logrus.Info("RemoveCategoryById")
	ctx, cancel := core.InitContext()
	defer cancel()
	var data model.Category
	objId, _ := primitive.ObjectIDFromHex(id)
	err := entity.categoryRepo.FindOne(ctx, bson.M{"_id": objId}).Decode(&data)
	if err != nil {
		logrus.Error(err)
		return nil, http.StatusBadRequest, err
	}
	_, err = entity.categoryRepo.DeleteOne(ctx, bson.M{"_id": objId})
	if err != nil {
		logrus.Error(err)
		return nil, http.StatusBadRequest, err
	}
	return &data, http.StatusOK, nil
}

func (entity categoryEntity) UpdateCategoryById(id string, form form.Category) (*model.Category, int, error) {
	logrus.Info("UpdateCategoryById")
	ctx, cancel := core.InitContext()
	defer cancel()
	objId, _ := primitive.ObjectIDFromHex(id)
	var data model.Category
	err := entity.categoryRepo.FindOne(ctx, bson.M{"_id": objId}).Decode(&data)
	if err != nil {
		logrus.Error(err)
		return nil, http.StatusBadRequest, err
	}
	data.Name = form.Name
	data.Value = strings.ToUpper(form.Value)
	data.Description = form.Description
	data.UpdatedDate = time.Now()

	isReturnNewDoc := options.After
	opts := &options.FindOneAndUpdateOptions{
		ReturnDocument: &isReturnNewDoc,
	}
	err = entity.categoryRepo.FindOneAndUpdate(ctx, bson.M{"_id": objId}, bson.M{"$set": data}, opts).Decode(&data)
	if err != nil {
		logrus.Error(err)
		return nil, http.StatusBadRequest, err
	}
	return &data, http.StatusOK, nil
}

func (entity *categoryEntity) CreateIndex() (string, error) {
	ctx, cancel := core.InitContext()
	defer cancel()
	mod := mongo.IndexModel{
		Keys: bson.M{
			"value": 1,
		},
		Options: options.Index().SetUnique(true),
	}
	ind, err := entity.categoryRepo.Indexes().CreateOne(ctx, mod)
	return ind, err
}

type ICategory interface {
	CreateIndex() (string, error)
	GetCategoryAll() ([]model.Category, int, error)
	CreateCategory(form form.Category) (*model.Category, int, error)
	GetCategoryById(id string) (*model.Category, int, error)
	RemoveCategoryById(id string) (*model.Category, int, error)
	UpdateCategoryById(id string, form form.Category) (*model.Category, int, error)
	UpdateDefaultCategoryById(id string) (*model.Category, int, error)
}

func NewCategoryEntity(resource *db.Resource) ICategory {
	categoryRepo := resource.DB.Collection("categories")
	CategoryEntity = &categoryEntity{resource: resource, categoryRepo: categoryRepo}
	return CategoryEntity
}
