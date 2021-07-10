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
	"mgo-gin/utils/bcrypt"
	"mgo-gin/utils/constant"
	"net/http"
	"time"
)

var UserEntity IUser

type userEntity struct {
	resource   *db.Resource
	userRepo   *mongo.Collection
	verifyRepo *mongo.Collection
}

type IUser interface {
	GetAll() ([]model.User, int, error)
	GetOneByUsername(username string) (*model.User, int, error)
	GetOneById(id string) (*model.User, int, error)
	CreateOne(form form.User) (*model.User, int, error)
	RemoveOneById(id string) (*model.User, int, error)
	UpdateUserById(id string, form form.User) (*model.User, int, error)
	ChangePassword(id string, form form.ChangePassword) (*model.User, int, error)
	SetPassword(id string, form form.SetPassword) (*model.User, int, error)
	CreateVerification(id string, objective string) (*model.UserReference, int, error)
	UpdateVerification(form form.VerifyRequest) (*model.UserReference, int, error)
	ActiveVerification(userRefId string) (*model.UserReference, int, error)
	RemoveVerification(userRefId string) (*model.UserReference, int, error)
	GetVerificationById(userRefId string) (*model.UserReference, int, error)
}

func NewUserEntity(resource *db.Resource) IUser {
	userRepo := resource.DB.Collection("users")
	verifyRepo := resource.DB.Collection("verifications")
	UserEntity = &userEntity{resource: resource, userRepo: userRepo, verifyRepo: verifyRepo}
	return UserEntity
}

func (entity *userEntity) GetAll() ([]model.User, int, error) {
	logrus.Info("GetAll")
	var usersList []model.User
	ctx, cancel := core.InitContext()
	defer cancel()
	cursor, err := entity.userRepo.Find(ctx, bson.M{})
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
	err := entity.userRepo.FindOne(ctx, bson.M{"username": username}).Decode(&user)
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
		Password:    bcrypt.HashPassword(form.Password),
		Role:        constant.USER,
		Status:      constant.ACTIVE,
		CreatedBy:   createdBy,
		CreatedDate: time.Now(),
		UpdatedBy:   createdBy,
		UpdatedDate: time.Now(),
	}
	_, err := entity.userRepo.InsertOne(ctx, user)
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
	err := entity.userRepo.FindOne(ctx, bson.M{"_id": objId}).Decode(&user)
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
	err := entity.userRepo.FindOne(ctx, bson.M{"_id": objId}).Decode(&user)
	if err != nil {
		logrus.Error(err)
		return nil, http.StatusBadRequest, err
	}
	_, err = entity.userRepo.DeleteOne(ctx, bson.M{"_id": objId})
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
	user.Password = bcrypt.HashPassword(form.Password)
	user.UpdatedBy, _ = primitive.ObjectIDFromHex(form.UpdatedBy)
	user.UpdatedDate = time.Now()

	isReturnNewDoc := options.After
	opts := &options.FindOneAndUpdateOptions{
		ReturnDocument: &isReturnNewDoc,
	}
	err = entity.userRepo.FindOneAndUpdate(ctx, bson.M{"_id": objId}, bson.M{"$set": user}, opts).Decode(&user)
	if err != nil {
		logrus.Error(err)
		return nil, http.StatusBadRequest, err
	}
	return user, http.StatusOK, nil
}

func (entity *userEntity) ChangePassword(id string, form form.ChangePassword) (*model.User, int, error) {
	logrus.Info("ChangePassword")
	ctx, cancel := core.InitContext()
	defer cancel()
	objId, _ := primitive.ObjectIDFromHex(id)
	user, _, err := entity.GetOneById(id)
	if err != nil {
		logrus.Error(err)
		return nil, http.StatusNotFound, err
	}
	user.Password = bcrypt.HashPassword(form.NewPassword)
	user.UpdatedBy = objId
	user.UpdatedDate = time.Now()
	isReturnNewDoc := options.After
	opts := &options.FindOneAndUpdateOptions{
		ReturnDocument: &isReturnNewDoc,
	}
	err = entity.userRepo.FindOneAndUpdate(ctx, bson.M{"_id": objId}, bson.M{"$set": user}, opts).Decode(&user)
	if err != nil {
		logrus.Error(err)
		return nil, http.StatusBadRequest, err
	}
	return user, http.StatusOK, nil
}

func (entity *userEntity) SetPassword(id string, form form.SetPassword) (*model.User, int, error) {
	logrus.Info("SetPassword")
	ctx, cancel := core.InitContext()
	defer cancel()
	objId, _ := primitive.ObjectIDFromHex(id)
	user, _, err := entity.GetOneById(id)
	if err != nil {
		logrus.Error(err)
		return nil, http.StatusNotFound, err
	}
	user.Password = bcrypt.HashPassword(form.Password)
	user.UpdatedBy = objId
	user.UpdatedDate = time.Now()
	isReturnNewDoc := options.After
	opts := &options.FindOneAndUpdateOptions{
		ReturnDocument: &isReturnNewDoc,
	}
	err = entity.userRepo.FindOneAndUpdate(ctx, bson.M{"_id": objId}, bson.M{"$set": user}, opts).Decode(&user)
	if err != nil {
		logrus.Error(err)
		return nil, http.StatusBadRequest, err
	}
	return user, http.StatusOK, nil
}

func (entity *userEntity) CreateVerification(id string, objective string) (*model.UserReference, int, error) {
	logrus.Info("CreateVerification")
	ctx, cancel := core.InitContext()
	defer cancel()
	var userId, _ = primitive.ObjectIDFromHex(id)
	_, err := entity.verifyRepo.DeleteMany(ctx, bson.M{"userId": userId})
	var userRefId = primitive.NewObjectID()
	user := model.UserReference{
		Id:          userRefId,
		UserId:      userId,
		Objective:   objective,
		CreatedDate: time.Now(),
		Status:      constant.INACTIVE,
	}
	_, err = entity.verifyRepo.InsertOne(ctx, user)
	if err != nil {
		logrus.Error(err)
		return nil, http.StatusBadRequest, err
	}
	return &user, http.StatusOK, nil
}

func (entity *userEntity) UpdateVerification(form form.VerifyRequest) (*model.UserReference, int, error) {
	logrus.Info("UpdateVerification")
	ctx, cancel := core.InitContext()
	defer cancel()
	var user model.UserReference
	objId, _ := primitive.ObjectIDFromHex(form.UserRefId)
	err := entity.verifyRepo.FindOne(ctx, bson.M{"_id": objId}).Decode(&user)
	if err != nil {
		logrus.Error(err)
		return nil, http.StatusBadRequest, err
	}
	user.Channel = form.Channel
	user.ChannelInfo = form.ChannelInfo
	user.Code, _ = bcrypt.GenerateCode(6)
	user.RefId, _ = bcrypt.GenerateRefId(4)
	user.ExpireDate = time.Now().Add(5 * time.Minute)
	user.ValidPeriod = 5
	isReturnNewDoc := options.After
	opts := &options.FindOneAndUpdateOptions{
		ReturnDocument: &isReturnNewDoc,
	}
	err = entity.verifyRepo.FindOneAndUpdate(ctx, bson.M{"_id": objId}, bson.M{"$set": user}, opts).Decode(&user)
	if err != nil {
		logrus.Error(err)
		return nil, http.StatusBadRequest, err
	}
	return &user, http.StatusOK, nil
}

func (entity *userEntity) ActiveVerification(userRefId string) (*model.UserReference, int, error) {
	logrus.Info("ActiveVerification")
	ctx, cancel := core.InitContext()
	defer cancel()
	var user model.UserReference
	objId, _ := primitive.ObjectIDFromHex(userRefId)
	err := entity.verifyRepo.FindOne(ctx, bson.M{"_id": objId}).Decode(&user)
	if err != nil {
		logrus.Error(err)
		return nil, http.StatusBadRequest, err
	}
	user.Status = constant.ACTIVE
	isReturnNewDoc := options.After
	opts := &options.FindOneAndUpdateOptions{
		ReturnDocument: &isReturnNewDoc,
	}
	err = entity.verifyRepo.FindOneAndUpdate(ctx, bson.M{"_id": objId}, bson.M{"$set": user}, opts).Decode(&user)
	if err != nil {
		logrus.Error(err)
		return nil, http.StatusBadRequest, err
	}
	return &user, http.StatusOK, nil
}

func (entity *userEntity) GetVerificationById(userRefId string) (*model.UserReference, int, error) {
	logrus.Info("GetVerificationById")
	ctx, cancel := core.InitContext()
	defer cancel()
	var user model.UserReference
	objId, _ := primitive.ObjectIDFromHex(userRefId)
	err := entity.verifyRepo.FindOne(ctx, bson.M{"_id": objId}).Decode(&user)
	if err != nil {
		logrus.Error(err)
		return nil, http.StatusBadRequest, err
	}
	return &user, http.StatusOK, nil
}

func (entity *userEntity) RemoveVerification(userRefId string) (*model.UserReference, int, error) {
	logrus.Info("RemoveVerification")
	ctx, cancel := core.InitContext()
	defer cancel()
	var user model.UserReference
	objId, _ := primitive.ObjectIDFromHex(userRefId)
	err := entity.verifyRepo.FindOne(ctx, bson.M{"_id": objId}).Decode(&user)
	if err != nil {
		logrus.Error(err)
		return nil, http.StatusBadRequest, err
	}
	_, _ = entity.verifyRepo.DeleteOne(ctx, bson.M{"_id": objId})
	return &user, http.StatusOK, nil
}
