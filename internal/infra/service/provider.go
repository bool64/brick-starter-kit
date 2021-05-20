package service

import (
	"github.com/bool64/brick-template/internal/domain/greeting"
)

type GreetingMakerProvider interface {
	GreetingMaker() greeting.Maker
}
