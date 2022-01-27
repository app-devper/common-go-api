package repository

import (
	"devper/app/core/utils"
	"devper/app/featues/notification/form"
	"devper/app/featues/notification/model"
	"devper/db"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Entity INotification

type notificationEntity struct {
	resource *db.Resource
	repo     *mongo.Collection
}

type INotification interface {
	Subscription(form form.Subscription) (*model.Subscription, error)
	GetOneByUserId(userId string) (*model.Subscription, error)
}

func NewNotificationEntity(resource *db.Resource) INotification {
	pushRepo := resource.DB.Collection("push_devices")
	Entity = &notificationEntity{resource: resource, repo: pushRepo}
	return Entity
}

func (entity *notificationEntity) Subscription(form form.Subscription) (*model.Subscription, error) {
	logrus.Info("Subscription")
	ctx, cancel := utils.InitContext()
	defer cancel()
	userId, _ := primitive.ObjectIDFromHex(form.UserId)
	subscription := model.Subscription{
		UserId:      userId,
		DeviceToken: form.DeviceToken,
		Channel:     form.Channel,
	}
	found, _ := entity.GetOneByUserId(form.UserId)
	if found != nil {
		subscription.Id = found.Id
		isReturnNewDoc := options.After
		opts := &options.FindOneAndUpdateOptions{
			ReturnDocument: &isReturnNewDoc,
		}
		err := entity.repo.FindOneAndUpdate(ctx, bson.M{"userId": userId}, bson.M{"$set": subscription}, opts).Decode(&subscription)
		if err != nil {
			return nil, err
		}
		return &subscription, nil
	} else {
		subscription.Id = primitive.NewObjectID()
		_, err := entity.repo.InsertOne(ctx, subscription)
		if err != nil {
			return nil, err
		}
		return &subscription, nil
	}
}

func (entity *notificationEntity) GetOneByUserId(userId string) (*model.Subscription, error) {
	logrus.Info("GetOneByUserId")
	ctx, cancel := utils.InitContext()
	defer cancel()
	objId, _ := primitive.ObjectIDFromHex(userId)
	var user model.Subscription
	err := entity.repo.FindOne(ctx, bson.M{"userId": objId}).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
