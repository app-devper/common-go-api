package repository

import (
	"errors"
	"github.com/jinzhu/copier"
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
	CreateOne(userForm form.User) (*model.User, int, error)
	RemoveOneById(id string) (*model.User, int, error)
	UpdateUserById(id string, userForm form.User) (*model.User, int, error)
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

func (entity *userEntity) CreateOne(userForm form.User) (*model.User, int, error) {
	logrus.Info("CreateOne")
	ctx, cancel := core.InitContext()
	defer cancel()
	user := model.User{
		Id:       primitive.NewObjectID(),
		Username: userForm.Username,
		Password: userForm.Password,
		Role:     constant.USER,
	}
	found, _, _ := entity.GetOneByUsername(user.Username)
	if found != nil {
		return nil, http.StatusConflict, errors.New("username is taken")
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

func (entity *userEntity) UpdateUserById(id string, userForm form.User) (*model.User, int, error) {
	logrus.Info("UpdateUserById")
	ctx, cancel := core.InitContext()
	defer cancel()
	objId, _ := primitive.ObjectIDFromHex(id)
	user, _, err := entity.GetOneById(id)
	if err != nil {
		return nil, http.StatusNotFound, err
	}
	err = copier.Copy(user, userForm) // this is why we need return a pointer: to copy value
	if err != nil {
		logrus.Error(err)
		return nil, http.StatusBadRequest, err
	}

	isReturnNewDoc := options.After
	opts := &options.FindOneAndUpdateOptions{
		ReturnDocument: &isReturnNewDoc,
	}
	err = entity.repo.FindOneAndUpdate(ctx, bson.M{"_id": objId}, bson.M{"$set": user}, opts).Decode(user)
	if err != nil {
		logrus.Error(err)
		return nil, http.StatusBadRequest, err
	}
	return user, http.StatusOK, nil
}
