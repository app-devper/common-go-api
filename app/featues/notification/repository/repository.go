package repository

import (
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"mgo-gin/app/core"
	"mgo-gin/app/featues/notification/form"
	"mgo-gin/app/featues/notification/model"
	"mgo-gin/db"
	"net/http"
)

var NotificationEntity INotification

type notificationEntity struct {
	resource *db.Resource
	repo     *mongo.Collection
}

type INotification interface {
	Subscription(form form.Subscription) (*model.Subscription, int, error)
	GetOneByUserId(userId string) (*model.Subscription, int, error)
}

func NewNotificationEntity(resource *db.Resource) INotification {
	pushRepo := resource.DB.Collection("push_devices")
	NotificationEntity = &notificationEntity{resource: resource, repo: pushRepo}
	return NotificationEntity
}

func (entity *notificationEntity) Subscription(form form.Subscription) (*model.Subscription, int, error) {
	logrus.Info("Subscription")
	ctx, cancel := core.InitContext()
	defer cancel()
	userId, _ := primitive.ObjectIDFromHex(form.UserId)
	subscription := model.Subscription{
		UserId:      userId,
		DeviceToken: form.DeviceToken,
		Channel:     form.Channel,
	}
	found, _, _ := entity.GetOneByUserId(form.UserId)
	if found != nil {
		subscription.Id = found.Id
		isReturnNewDoc := options.After
		opts := &options.FindOneAndUpdateOptions{
			ReturnDocument: &isReturnNewDoc,
		}
		err := entity.repo.FindOneAndUpdate(ctx, bson.M{"userId": userId}, bson.M{"$set": subscription}, opts).Decode(&subscription)
		if err != nil {
			logrus.Error(err)
			return nil, http.StatusBadRequest, err
		}
		return &subscription, http.StatusOK, nil
	} else {
		subscription.Id = primitive.NewObjectID()
		_, err := entity.repo.InsertOne(ctx, subscription)
		if err != nil {
			logrus.Error(err)
			return nil, http.StatusBadRequest, err
		}
		return &subscription, http.StatusOK, nil
	}
}

func (entity *notificationEntity) GetOneByUserId(userId string) (*model.Subscription, int, error) {
	logrus.Info("GetOneByUserId")
	ctx, cancel := core.InitContext()
	defer cancel()
	objId, _ := primitive.ObjectIDFromHex(userId)
	var user model.Subscription
	err := entity.repo.FindOne(ctx, bson.M{"userId": objId}).Decode(&user)
	if err != nil {
		logrus.Error(err)
		return nil, http.StatusBadRequest, err
	}
	return &user, http.StatusOK, nil
}
