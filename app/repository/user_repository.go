package repository

import (
	"errors"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"mgo-gin/app/form"
	"mgo-gin/app/model"
	"mgo-gin/db"
	"mgo-gin/utils/bcrypt"
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
}

func NewUserEntity(resource *db.Resource) IUser {
	userRepo := resource.DB.Collection("users")
	UserEntity = &userEntity{resource: resource, repo: userRepo}
	return UserEntity
}

func (entity *userEntity) GetAll() ([]model.User, int, error) {
	logrus.Info("GetAll")
	var usersList []model.User
	ctx, cancel := initContext()
	defer cancel()

	cursor, err := entity.repo.Find(ctx, bson.M{})
	if err != nil {
		logrus.Error(err)
		return []model.User{}, http.StatusBadRequest, err
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
	ctx, cancel := initContext()
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
	ctx, cancel := initContext()
	defer cancel()
	user := model.User{
		Id:       primitive.NewObjectID(),
		Username: userForm.Username,
		Password: bcrypt.HashPassword(userForm.Password),
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
	ctx, cancel := initContext()
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
