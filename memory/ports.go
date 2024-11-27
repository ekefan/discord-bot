package memory

import (
	"errors"

	"github.com/ekefan/discord-bot/domain/challenge"
)

var (
	ErrInvalidChallenge   = errors.New("challenge is not valid")
	ErrSavingChallenge    = errors.New("can not create challenge")
	ErrInvalidChallengeId = errors.New("challenge id is not valid")
	ErrChallengeNotFound  = errors.New("challenge doesn't exist")
)

type ChallangeRespository interface {
	CreateChallenge(c *challenge.Challenge) error
	GetChallenge(id string) (*challenge.Challenge, error)
	DeleteChallenge(id string) error
}
