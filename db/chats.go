package db

import (
	"context"
	apitypes "lenscape-chat/api-types"
	"lenscape-chat/db/types"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func (db *DB) CreateChannel(channel *apitypes.ChannelCreation) (interface{}, error) {
	coll := db.client.Database("lenscape-chat").Collection("channels")
	findResult := coll.FindOne(context.Background(), bson.D{{"sender", channel.Sender}, {"receiver", channel.Receiver}})

	channelDetails := &types.Channel{}
	if err := findResult.Decode(channelDetails); err != nil && err != mongo.ErrNoDocuments {
		return nil, err
	}

	if channelDetails.ChannelID != primitive.NilObjectID {
		return channelDetails.ChannelID, nil
	}

	res, err := coll.InsertOne(context.TODO(), channel)
	if err != nil {
		return nil, err
	}

	msgsColl := db.client.Database("lenscape-chat").Collection("channel_messages")
	_, err = msgsColl.InsertOne(context.TODO(), bson.D{
		{"channel_id", res.InsertedID},
		// {"messages", []apitypes.ChatMessage{}},
	})
	if err != nil {
		return nil, err
	}
	return res.InsertedID, nil
}

func (db *DB) SendMessage(channel string, message *apitypes.ChatMessage) error {
	collection := db.client.Database("lenscape-chat").Collection("channel_messages")
	channelId, err := primitive.ObjectIDFromHex(channel)
	if err != nil {
		return err
	}

	change := bson.M{"$push": bson.M{"messages": message}}

	message.CreatedAt = time.Now()
	_, err = collection.UpdateOne(context.Background(),
		bson.M{"channel_id": channelId},
		change,
	)
	if err != nil {
		return err
	}
	return nil
}
