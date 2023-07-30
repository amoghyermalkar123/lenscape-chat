package types

import "go.mongodb.org/mongo-driver/bson/primitive"

type Channel struct {
	ChannelID   primitive.ObjectID `json:"channel_id" bson:"_id"`
	SenderEmail string             `json:"sender" bson:"sender"`
	RecvEmail   string             `json:"receiver" bson:"receiver"`
}
