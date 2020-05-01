package model

import (
	"github.com/romshark/messenger-sim/service/gateway/graph/model"
)

type Session struct {
	model.Session

	User *User `json:"user"`
}
