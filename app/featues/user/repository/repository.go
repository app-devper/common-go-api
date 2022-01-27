package repository

import (
	"devper/app/core"
	bcrypt2 "devper/app/core/bcrypt"
	"devper/app/featues/user/form"
	"devper/app/featues/user/model"
	"devper/db"
	"devper/utils/constant"
	"errors"
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
	GetUserByRefId(userRefId string, objective string) (*model.User, error)
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
	_, err := Entity.CreateIndex()
	if err != nil {
		logrus.Error(err)
	}
	return Entity
}

func (entity *userEntity) CreateIndex() (string, error) {
	ctx, cancel := core.InitContext()
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

func (entity *userEntity) GetUserByRefId(userRefId string, objective string) (*model.User, error) {
	logrus.Info("GetUserByRefId")
	userRef, _ := entity.GetVerificationById(userRefId)
	if userRef == nil {
		return nil, errors.New("user ref invalid")
	}
	if userRef.Status != constant.ACTIVE {
		return nil, errors.New("user ref not active")
	}
	if userRef.Objective != objective {
		return nil, errors.New("wrong objective")
	}
	if userRef.ExpireDate.Before(time.Now()) {
		return nil, errors.New("token invalid")
	}
	user, _ := entity.GetUserById(userRef.UserId.Hex())
	if user == nil {
		return nil, errors.New("user invalid")
	}
	if user.Status != constant.ACTIVE {
		return nil, errors.New("user not active")
	}
	return user, nil
}

func (entity *userEntity) GetUserAll() ([]model.User, error) {
	logrus.Info("GetUserAll")
	var usersList []model.User
	ctx, cancel := core.InitContext()
	defer cancel()
	cursor, err := entity.userRepo.Find(ctx, bson.M{})
	if err != nil {
		logrus.Error(err)
		return nil, err
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
	return usersList, nil
}

func (entity *userEntity) GetUserByUsername(username string) (*model.User, error) {
	logrus.Info("GetUserByUsername")
	ctx, cancel := core.InitContext()
	defer cancel()
	var user model.User
	err := entity.userRepo.FindOne(ctx, bson.M{"username": strings.TrimSpace(username)}).Decode(&user)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	if user.Status != constant.ACTIVE {
		return nil, errors.New("user not active")
	}
	return &user, nil
}

func (entity *userEntity) CreateUser(form form.User) (*model.User, error) {
	logrus.Info("CreateUser")
	ctx, cancel := core.InitContext()
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
		Password:    bcrypt2.HashPassword(form.Password),
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
		return nil, err
	}
	return &user, nil
}

func (entity *userEntity) GetUserById(id string) (*model.User, error) {
	logrus.Info("GetUserById")
	ctx, cancel := core.InitContext()
	defer cancel()
	var user model.User
	objId, _ := primitive.ObjectIDFromHex(id)
	err := entity.userRepo.FindOne(ctx, bson.M{"_id": objId}).Decode(&user)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	return &user, nil
}

func (entity *userEntity) RemoveUserById(id string) (*model.User, error) {
	logrus.Info("RemoveUserById")
	ctx, cancel := core.InitContext()
	defer cancel()
	var user model.User
	objId, _ := primitive.ObjectIDFromHex(id)
	err := entity.userRepo.FindOne(ctx, bson.M{"_id": objId}).Decode(&user)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	_, err = entity.userRepo.DeleteOne(ctx, bson.M{"_id": objId})
	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	return &user, nil
}

func (entity *userEntity) UpdateUserById(id string, form form.UpdateUser) (*model.User, error) {
	logrus.Info("UpdateUserById")
	ctx, cancel := core.InitContext()
	defer cancel()
	objId, _ := primitive.ObjectIDFromHex(id)
	user, err := entity.GetUserById(id)
	if err != nil {
		logrus.Error(err)
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
		logrus.Error(err)
		return nil, err
	}
	return user, nil
}

func (entity *userEntity) UpdateStatusById(id string, form form.UpdateStatus) (*model.User, error) {
	logrus.Info("UpdateStatusById")
	ctx, cancel := core.InitContext()
	defer cancel()
	objId, _ := primitive.ObjectIDFromHex(id)
	user, err := entity.GetUserById(id)
	if err != nil {
		logrus.Error(err)
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
		logrus.Error(err)
		return nil, err
	}
	return user, nil
}

func (entity *userEntity) UpdateRoleById(id string, form form.UpdateRole) (*model.User, error) {
	logrus.Info("UpdateRoleById")
	ctx, cancel := core.InitContext()
	defer cancel()
	objId, _ := primitive.ObjectIDFromHex(id)
	user, err := entity.GetUserById(id)
	if err != nil {
		logrus.Error(err)
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
		logrus.Error(err)
		return nil, err
	}
	return user, nil
}

func (entity *userEntity) ChangePassword(id string, form form.ChangePassword) (*model.User, error) {
	logrus.Info("ChangePassword")
	ctx, cancel := core.InitContext()
	defer cancel()
	objId, _ := primitive.ObjectIDFromHex(id)
	user, err := entity.GetUserById(id)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	user.Password = bcrypt2.HashPassword(form.NewPassword)
	user.UpdatedBy = objId
	user.UpdatedDate = time.Now()
	isReturnNewDoc := options.After
	opts := &options.FindOneAndUpdateOptions{
		ReturnDocument: &isReturnNewDoc,
	}
	err = entity.userRepo.FindOneAndUpdate(ctx, bson.M{"_id": objId}, bson.M{"$set": user}, opts).Decode(&user)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	return user, nil
}

func (entity *userEntity) SetPassword(id string, form form.SetPassword) (*model.User, error) {
	logrus.Info("SetPassword")
	ctx, cancel := core.InitContext()
	defer cancel()
	objId, _ := primitive.ObjectIDFromHex(id)
	user, err := entity.GetUserById(id)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	user.Password = bcrypt2.HashPassword(form.Password)
	user.UpdatedBy = objId
	user.UpdatedDate = time.Now()
	isReturnNewDoc := options.After
	opts := &options.FindOneAndUpdateOptions{
		ReturnDocument: &isReturnNewDoc,
	}
	err = entity.userRepo.FindOneAndUpdate(ctx, bson.M{"_id": objId}, bson.M{"$set": user}, opts).Decode(&user)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	return user, nil
}

func (entity *userEntity) CreateVerification(form form.Reference) (*model.UserReference, error) {
	logrus.Info("CreateVerification")
	ctx, cancel := core.InitContext()
	defer cancel()

	var userRefId = primitive.NewObjectID()
	reference := model.UserReference{
		Id:          userRefId,
		UserId:      form.UserId,
		Objective:   form.Objective,
		Channel:     form.Channel,
		ChannelInfo: form.ChannelInfo,
		CreatedDate: time.Now(),
		ExpireDate:  form.ExpireDate,
		Status:      form.Status,
	}
	_, err := entity.verifyRepo.InsertOne(ctx, reference)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	return &reference, nil
}

func (entity *userEntity) UpdateVerification(form form.VerifyChannel, expireTime time.Time) (*model.UserReference, error) {
	logrus.Info("UpdateVerification")
	ctx, cancel := core.InitContext()
	defer cancel()
	var reference model.UserReference
	objId, _ := primitive.ObjectIDFromHex(form.UserRefId)
	err := entity.verifyRepo.FindOne(ctx, bson.M{"_id": objId}).Decode(&reference)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	reference.Channel = form.Channel
	reference.ChannelInfo = form.ChannelInfo
	reference.Code, _ = bcrypt2.GenerateCode(6)
	reference.RefId, _ = bcrypt2.GenerateRefId(4)
	reference.ExpireDate = expireTime
	reference.ValidPeriod = 5
	isReturnNewDoc := options.After
	opts := &options.FindOneAndUpdateOptions{
		ReturnDocument: &isReturnNewDoc,
	}
	err = entity.verifyRepo.FindOneAndUpdate(ctx, bson.M{"_id": objId}, bson.M{"$set": reference}, opts).Decode(&reference)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	return &reference, nil
}

func (entity *userEntity) ActiveVerification(userRefId string, expireTime time.Time) (*model.UserReference, error) {
	logrus.Info("ActiveVerification")
	ctx, cancel := core.InitContext()
	defer cancel()
	var reference model.UserReference
	objId, _ := primitive.ObjectIDFromHex(userRefId)
	err := entity.verifyRepo.FindOne(ctx, bson.M{"_id": objId}).Decode(&reference)
	if err != nil {
		logrus.Error(err)
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
		logrus.Error(err)
		return nil, err
	}
	return &reference, nil
}

func (entity *userEntity) GetVerificationById(userRefId string) (*model.UserReference, error) {
	logrus.Info("GetVerificationById")
	ctx, cancel := core.InitContext()
	defer cancel()
	var reference model.UserReference
	objId, _ := primitive.ObjectIDFromHex(userRefId)
	err := entity.verifyRepo.FindOne(ctx, bson.M{"_id": objId}).Decode(&reference)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	return &reference, nil
}

func (entity *userEntity) RevokeVerification(userRefId string) (*model.UserReference, error) {
	logrus.Info("RevokeVerification")
	ctx, cancel := core.InitContext()
	defer cancel()
	var reference model.UserReference
	objId, _ := primitive.ObjectIDFromHex(userRefId)
	err := entity.verifyRepo.FindOne(ctx, bson.M{"_id": objId}).Decode(&reference)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	reference.ExpireDate = time.Now()
	isReturnNewDoc := options.After
	opts := &options.FindOneAndUpdateOptions{
		ReturnDocument: &isReturnNewDoc,
	}
	err = entity.verifyRepo.FindOneAndUpdate(ctx, bson.M{"_id": objId}, bson.M{"$set": reference}, opts).Decode(&reference)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	return &reference, nil
}

func (entity *userEntity) RemoveVerificationObjective(objective string) error {
	logrus.Info("RemoveVerificationObjective")
	ctx, cancel := core.InitContext()
	defer cancel()
	_, err := entity.verifyRepo.DeleteMany(ctx, bson.M{"objective": objective})
	if err != nil {
		logrus.Error(err)
		return err
	}
	return nil
}
