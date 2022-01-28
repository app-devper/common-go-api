package repository

import (
	"devper/app/core/constant"
	"devper/app/core/utils"
	"devper/app/featues/user/form"
	"devper/app/featues/user/model"
	"devper/db"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"strings"
	"time"
)

var Entity IUser

type userEntity struct {
	resource   *db.Resource
	userRepo   *mongo.Collection
	verifyRepo *mongo.Collection
}

type IUser interface {
	CreateIndex() (string, error)
	GetUserAll() ([]model.User, error)
	GetUserByUsername(username string) (*model.User, error)
	GetUserById(id string) (*model.User, error)
	CreateUser(form form.User) (*model.User, error)
	RemoveUserById(id string) (*model.User, error)
	UpdateUserById(id string, form form.UpdateUser) (*model.User, error)
	UpdateStatusById(id string, form form.UpdateStatus) (*model.User, error)
	UpdateRoleById(id string, form form.UpdateRole) (*model.User, error)
	ChangePassword(id string, form form.ChangePassword) (*model.User, error)
	SetPassword(id string, form form.SetPassword) (*model.User, error)

	CreateVerification(form form.Reference) (*model.UserReference, error)
	UpdateVerification(form form.VerifyChannel, expireTime time.Time) (*model.UserReference, error)
	ActiveVerification(userRefId string, expireTime time.Time) (*model.UserReference, error)
	RevokeVerification(userRefId string) (*model.UserReference, error)
	RemoveVerificationObjective(objective string) error
	GetVerificationById(userRefId string) (*model.UserReference, error)
}

func NewUserEntity(resource *db.Resource) IUser {
	userRepo := resource.DB.Collection("users")
	verifyRepo := resource.DB.Collection("verifications")
	Entity = &userEntity{resource: resource, userRepo: userRepo, verifyRepo: verifyRepo}
	_, _ = Entity.CreateIndex()
	return Entity
}

func (entity *userEntity) CreateIndex() (string, error) {
	ctx, cancel := utils.InitContext()
	defer cancel()
	mod := mongo.IndexModel{
		Keys: bson.M{
			"username": 1,
		},
		Options: options.Index().SetUnique(true),
	}
	ind, err := entity.userRepo.Indexes().CreateOne(ctx, mod)
	return ind, err
}

func (entity *userEntity) GetUserAll() ([]model.User, error) {
	logrus.Info("GetUserAll")
	var usersList []model.User
	ctx, cancel := utils.InitContext()
	defer cancel()
	cursor, err := entity.userRepo.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	for cursor.Next(ctx) {
		var user model.User
		err = cursor.Decode(&user)
		if err != nil {
			logrus.Error(err)
			logrus.Info(cursor.Current)
		} else {
			usersList = append(usersList, user)
		}
	}
	if usersList == nil {
		usersList = []model.User{}
	}
	return usersList, nil
}

func (entity *userEntity) GetUserByUsername(username string) (*model.User, error) {
	logrus.Info("GetUserByUsername")
	ctx, cancel := utils.InitContext()
	defer cancel()
	var user model.User
	err := entity.userRepo.FindOne(ctx, bson.M{"username": strings.TrimSpace(username)}).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (entity *userEntity) CreateUser(form form.User) (*model.User, error) {
	logrus.Info("CreateUser")
	ctx, cancel := utils.InitContext()
	defer cancel()

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
		Password:    utils.HashPassword(form.Password),
		Role:        constant.USER,
		Status:      constant.ACTIVE,
		CreatedBy:   createdBy,
		CreatedDate: time.Now(),
		UpdatedBy:   createdBy,
		UpdatedDate: time.Now(),
	}
	_, err := entity.userRepo.InsertOne(ctx, user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (entity *userEntity) GetUserById(id string) (*model.User, error) {
	logrus.Info("GetUserById")
	ctx, cancel := utils.InitContext()
	defer cancel()
	var user model.User
	objId, _ := primitive.ObjectIDFromHex(id)
	err := entity.userRepo.FindOne(ctx, bson.M{"_id": objId}).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (entity *userEntity) RemoveUserById(id string) (*model.User, error) {
	logrus.Info("RemoveUserById")
	ctx, cancel := utils.InitContext()
	defer cancel()
	var user model.User
	objId, _ := primitive.ObjectIDFromHex(id)
	err := entity.userRepo.FindOne(ctx, bson.M{"_id": objId}).Decode(&user)
	if err != nil {
		return nil, err
	}
	_, err = entity.userRepo.DeleteOne(ctx, bson.M{"_id": objId})
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (entity *userEntity) UpdateUserById(id string, form form.UpdateUser) (*model.User, error) {
	logrus.Info("UpdateUserById")
	ctx, cancel := utils.InitContext()
	defer cancel()
	objId, _ := primitive.ObjectIDFromHex(id)
	user, err := entity.GetUserById(id)
	if err != nil {
		return nil, err
	}

	user.FirstName = form.FirstName
	user.LastName = form.LastName
	user.Email = form.Email
	user.Phone = form.Phone
	user.UpdatedBy, _ = primitive.ObjectIDFromHex(form.UpdatedBy)
	user.UpdatedDate = time.Now()

	isReturnNewDoc := options.After
	opts := &options.FindOneAndUpdateOptions{
		ReturnDocument: &isReturnNewDoc,
	}
	err = entity.userRepo.FindOneAndUpdate(ctx, bson.M{"_id": objId}, bson.M{"$set": user}, opts).Decode(&user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (entity *userEntity) UpdateStatusById(id string, form form.UpdateStatus) (*model.User, error) {
	logrus.Info("UpdateStatusById")
	ctx, cancel := utils.InitContext()
	defer cancel()
	objId, _ := primitive.ObjectIDFromHex(id)
	user, err := entity.GetUserById(id)
	if err != nil {
		return nil, err
	}
	user.Status = form.Status
	user.UpdatedBy, _ = primitive.ObjectIDFromHex(form.UpdatedBy)
	user.UpdatedDate = time.Now()

	isReturnNewDoc := options.After
	opts := &options.FindOneAndUpdateOptions{
		ReturnDocument: &isReturnNewDoc,
	}
	err = entity.userRepo.FindOneAndUpdate(ctx, bson.M{"_id": objId}, bson.M{"$set": user}, opts).Decode(&user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (entity *userEntity) UpdateRoleById(id string, form form.UpdateRole) (*model.User, error) {
	logrus.Info("UpdateRoleById")
	ctx, cancel := utils.InitContext()
	defer cancel()
	objId, _ := primitive.ObjectIDFromHex(id)
	user, err := entity.GetUserById(id)
	if err != nil {
		return nil, err
	}

	user.Role = form.Role
	user.UpdatedBy, _ = primitive.ObjectIDFromHex(form.UpdatedBy)
	user.UpdatedDate = time.Now()

	isReturnNewDoc := options.After
	opts := &options.FindOneAndUpdateOptions{
		ReturnDocument: &isReturnNewDoc,
	}
	err = entity.userRepo.FindOneAndUpdate(ctx, bson.M{"_id": objId}, bson.M{"$set": user}, opts).Decode(&user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (entity *userEntity) ChangePassword(id string, form form.ChangePassword) (*model.User, error) {
	logrus.Info("ChangePassword")
	ctx, cancel := utils.InitContext()
	defer cancel()
	objId, _ := primitive.ObjectIDFromHex(id)
	user, err := entity.GetUserById(id)
	if err != nil {
		return nil, err
	}
	user.Password = utils.HashPassword(form.NewPassword)
	user.UpdatedBy = objId
	user.UpdatedDate = time.Now()
	isReturnNewDoc := options.After
	opts := &options.FindOneAndUpdateOptions{
		ReturnDocument: &isReturnNewDoc,
	}
	err = entity.userRepo.FindOneAndUpdate(ctx, bson.M{"_id": objId}, bson.M{"$set": user}, opts).Decode(&user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (entity *userEntity) SetPassword(id string, form form.SetPassword) (*model.User, error) {
	logrus.Info("SetPassword")
	ctx, cancel := utils.InitContext()
	defer cancel()
	objId, _ := primitive.ObjectIDFromHex(id)
	user, err := entity.GetUserById(id)
	if err != nil {
		return nil, err
	}
	user.Password = utils.HashPassword(form.Password)
	user.UpdatedBy = objId
	user.UpdatedDate = time.Now()
	isReturnNewDoc := options.After
	opts := &options.FindOneAndUpdateOptions{
		ReturnDocument: &isReturnNewDoc,
	}
	err = entity.userRepo.FindOneAndUpdate(ctx, bson.M{"_id": objId}, bson.M{"$set": user}, opts).Decode(&user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (entity *userEntity) CreateVerification(form form.Reference) (*model.UserReference, error) {
	logrus.Info("CreateVerification")
	ctx, cancel := utils.InitContext()
	defer cancel()

	var userRefId = primitive.NewObjectID()
	reference := model.UserReference{
		Id:          userRefId,
		UserId:      form.UserId,
		Type:        form.Type,
		Objective:   form.Objective,
		Channel:     form.Channel,
		ChannelInfo: form.ChannelInfo,
		CreatedDate: time.Now(),
		ExpireDate:  form.ExpireDate,
		Status:      form.Status,
	}
	_, err := entity.verifyRepo.InsertOne(ctx, reference)
	if err != nil {
		return nil, err
	}
	return &reference, nil
}

func (entity *userEntity) UpdateVerification(form form.VerifyChannel, expireTime time.Time) (*model.UserReference, error) {
	logrus.Info("UpdateVerification")
	ctx, cancel := utils.InitContext()
	defer cancel()
	var reference model.UserReference
	objId, _ := primitive.ObjectIDFromHex(form.UserRefId)
	err := entity.verifyRepo.FindOne(ctx, bson.M{"_id": objId}).Decode(&reference)
	if err != nil {
		return nil, err
	}
	reference.Channel = form.Channel
	reference.ChannelInfo = form.ChannelInfo
	reference.Code = utils.GenerateCode(6)
	reference.RefId = utils.GenerateRefId(4)
	reference.ExpireDate = expireTime
	reference.ValidPeriod = 5
	isReturnNewDoc := options.After
	opts := &options.FindOneAndUpdateOptions{
		ReturnDocument: &isReturnNewDoc,
	}
	err = entity.verifyRepo.FindOneAndUpdate(ctx, bson.M{"_id": objId}, bson.M{"$set": reference}, opts).Decode(&reference)
	if err != nil {
		return nil, err
	}
	return &reference, nil
}

func (entity *userEntity) ActiveVerification(userRefId string, expireTime time.Time) (*model.UserReference, error) {
	logrus.Info("ActiveVerification")
	ctx, cancel := utils.InitContext()
	defer cancel()
	var reference model.UserReference
	objId, _ := primitive.ObjectIDFromHex(userRefId)
	err := entity.verifyRepo.FindOne(ctx, bson.M{"_id": objId}).Decode(&reference)
	if err != nil {
		return nil, err
	}
	reference.Status = constant.ACTIVE
	reference.ExpireDate = expireTime
	isReturnNewDoc := options.After
	opts := &options.FindOneAndUpdateOptions{
		ReturnDocument: &isReturnNewDoc,
	}
	err = entity.verifyRepo.FindOneAndUpdate(ctx, bson.M{"_id": objId}, bson.M{"$set": reference}, opts).Decode(&reference)
	if err != nil {
		return nil, err
	}
	return &reference, nil
}

func (entity *userEntity) GetVerificationById(userRefId string) (*model.UserReference, error) {
	logrus.Info("GetVerificationById")
	ctx, cancel := utils.InitContext()
	defer cancel()
	var reference model.UserReference
	objId, _ := primitive.ObjectIDFromHex(userRefId)
	err := entity.verifyRepo.FindOne(ctx, bson.M{"_id": objId}).Decode(&reference)
	if err != nil {
		return nil, err
	}
	return &reference, nil
}

func (entity *userEntity) RevokeVerification(userRefId string) (*model.UserReference, error) {
	logrus.Info("RevokeVerification")
	ctx, cancel := utils.InitContext()
	defer cancel()
	var reference model.UserReference
	objId, _ := primitive.ObjectIDFromHex(userRefId)
	err := entity.verifyRepo.FindOne(ctx, bson.M{"_id": objId}).Decode(&reference)
	if err != nil {
		return nil, err
	}
	reference.ExpireDate = time.Now()
	isReturnNewDoc := options.After
	opts := &options.FindOneAndUpdateOptions{
		ReturnDocument: &isReturnNewDoc,
	}
	err = entity.verifyRepo.FindOneAndUpdate(ctx, bson.M{"_id": objId}, bson.M{"$set": reference}, opts).Decode(&reference)
	if err != nil {
		return nil, err
	}
	return &reference, nil
}

func (entity *userEntity) RemoveVerificationObjective(objective string) error {
	logrus.Info("RemoveVerificationObjective")
	ctx, cancel := utils.InitContext()
	defer cancel()
	_, err := entity.verifyRepo.DeleteMany(ctx, bson.M{"objective": objective})
	if err != nil {
		return err
	}
	return nil
}
