package stream

import (
	"encoding/json"

	"github.com/sirupsen/logrus"
)

func payment[T any](msg []byte) {
	var result T
	err := json.Unmarshal(msg, &result)
	if err != nil {
		logrus.WithFields(logrus.Fields{"msg": "failed to unmarshal message for payment topic"}).Error(err.Error())
	}
}

func campaign[T any](msg []byte) {
	var result T
	err := json.Unmarshal(msg, &result)
	if err != nil {
		logrus.WithFields(logrus.Fields{"msg": "failed to unmarshal message for campaign topic"}).Error(err.Error())
	}
}

func refund[T any](msg []byte) {
	var result T
	err := json.Unmarshal(msg, &result)
	if err != nil {
		logrus.WithFields(logrus.Fields{"msg": "failed to unmarshal message for refund topic"}).Error(err.Error())
	}
}

func transfer[T any](msg []byte) {
	var result T
	err := json.Unmarshal(msg, &result)
	if err != nil {
		logrus.WithFields(logrus.Fields{"msg": "failed to unmarshal message for transfer topic"}).Error(err.Error())
	} else
}

func system[T any](msg []byte) {
	var result T
	err := json.Unmarshal(msg, &result)
	if err != nil {
		logrus.WithFields(logrus.Fields{"msg": "failed to unmarshal message for system topic"}).Error(err.Error())
	}
}
