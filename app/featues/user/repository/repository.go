package repository

import (
	"errors"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"mgo-gin/app/core"
	"mgo-gin/app/featues/user/form"
	"mgo-gin/app/featues/user/model"
	"mgo-gin/db"
	"mgo-gin/utils/constant"
	"net/http"
	"time"
)

var UserEntity IUser

type userEntity struct {
	resource *db.Resource
	repo     *mongo.Collection
}

type IUser interface {
	GetAll() ([]model.User, int, error)
	GetOneByUsername(username string) (*model.User, int, error)
	GetOneById(id string) (*model.User, int, error)
	CreateOne(form form.User) (*model.User, int, error)
	RemoveOneById(id string) (*model.User, int, error)
	UpdateUserById(id string, form form.User) (*model.User, int, error)
}

func NewUserEntity(resource *db.Resource) IUser {
	userRepo := resource.DB.Collection("users")
	UserEntity = &userEntity{resource: resource, repo: userRepo}
	return UserEntity
}

func (entity *userEntity) GetAll() ([]model.User, int, error) {
	logrus.Info("GetAll")
	var usersList []model.User
	ctx, cancel := core.InitContext()
	defer cancel()
	cursor, err := entity.repo.Find(ctx, bson.M{})
	if err != nil {
		logrus.Error(err)
		return nil, http.StatusBadRequest, err
	}
	for cursor.Next(ctx) {
		var user model.User
		err = cursor.Decode(&user)
		if err != nil {
			logrus.Error(err)
		}
		usersList = append(usersList, user)
	}
	if usersList == nil {
		usersList = []model.User{}
	}
	return usersList, http.StatusOK, nil
}

func (entity *userEntity) GetOneByUsername(username string) (*model.User, int, error) {
	logrus.Info("GetOneByUsername")
	ctx, cancel := core.InitContext()
	defer cancel()
	var user model.User
	err := entity.repo.FindOne(ctx, bson.M{"username": username}).Decode(&user)
	if err != nil {
		logrus.Error(err)
		return nil, http.StatusBadRequest, err
	}
	return &user, http.StatusOK, nil
}

func (entity *userEntity) CreateOne(form form.User) (*model.User, int, error) {
	logrus.Info("CreateOne")
	ctx, cancel := core.InitContext()
	defer cancel()
	found, _, _ := entity.GetOneByUsername(form.Username)
	if found != nil {
		err := errors.New("username is taken")
		logrus.Error(err)
		return nil, http.StatusConflict, err
	}

	var userId = primitive.NewObjectID()
	var createdBy = userId
	if form.CreatedBy != "" {
		createdBy, _ = primitive.ObjectIDFromHex(form.CreatedBy)
	}
	user := model.User{
		Id:          userId,
		FirstName:   form.FirstName,
		LastName:    form.LastName,
		Username:    form.Username,
		Password:    form.Password,
		Role:        constant.USER,
		Status:      constant.ACTIVE,
		CreatedBy:   createdBy,
		CreatedDate: time.Now(),
		UpdatedBy:   createdBy,
		UpdatedDate: time.Now(),
	}
	_, err := entity.repo.InsertOne(ctx, user)
	if err != nil {
		logrus.Error(err)
		return nil, http.StatusBadRequest, err
	}
	return &user, http.StatusOK, nil
}

func (entity *userEntity) GetOneById(id string) (*model.User, int, error) {
	logrus.Info("GetOneById")
	ctx, cancel := core.InitContext()
	defer cancel()
	var user model.User
	objId, _ := primitive.ObjectIDFromHex(id)
	err := entity.repo.FindOne(ctx, bson.M{"_id": objId}).Decode(&user)
	if err != nil {
		logrus.Error(err)
		return nil, http.StatusBadRequest, err
	}
	return &user, http.StatusOK, nil
}

func (entity *userEntity) RemoveOneById(id string) (*model.User, int, error) {
	logrus.Info("RemoveOneById")
	ctx, cancel := core.InitContext()
	defer cancel()
	var user model.User
	objId, _ := primitive.ObjectIDFromHex(id)
	err := entity.repo.FindOne(ctx, bson.M{"_id": objId}).Decode(&user)
	if err != nil {
		logrus.Error(err)
		return nil, http.StatusBadRequest, err
	}
	_, err = entity.repo.DeleteOne(ctx, bson.M{"_id": objId})
	if err != nil {
		logrus.Error(err)
		return nil, http.StatusBadRequest, err
	}
	return &user, http.StatusOK, nil
}

func (entity *userEntity) UpdateUserById(id string, form form.User) (*model.User, int, error) {
	logrus.Info("UpdateUserById")
	ctx, cancel := core.InitContext()
	defer cancel()
	objId, _ := primitive.ObjectIDFromHex(id)
	user, _, err := entity.GetOneById(id)
	if err != nil {
		logrus.Error(err)
		return nil, http.StatusNotFound, err
	}
	user.FirstName = form.FirstName
	user.LastName = form.LastName
	user.Username = form.Username
	user.Password = form.Password
	user.UpdatedBy, _ = primitive.ObjectIDFromHex(form.UpdatedBy)
	user.UpdatedDate = time.Now()

	isReturnNewDoc := options.After
	opts := &options.FindOneAndUpdateOptions{
		ReturnDocument: &isReturnNewDoc,
	}
	err = entity.repo.FindOneAndUpdate(ctx, bson.M{"_id": objId}, bson.M{"$set": user}, opts).Decode(&user)
	if err != nil {
		logrus.Error(err)
		return nil, http.StatusBadRequest, err
	}
	return user, http.StatusOK, nil
}
