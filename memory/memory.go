package memory

import (
	"log/slog"
	"sync"

	"github.com/ekefan/discord-bot/domain/challenge"
)

type InMemory struct {
	challenges map[string]challenge.Challenge
	sync.Mutex
}

func NewInMemory() ChallangeRespository {
	return &InMemory{
		challenges: make(map[string]challenge.Challenge),
	}
}

func (im *InMemory) CreateChallenge(c *challenge.Challenge) error {
	if c == nil {
		return ErrInvalidChallenge
	}
	id, err := c.GetChallengeID()
	if err != nil {
		slog.Error("error getting challenge id", "details", err.Error())
		return ErrSavingChallenge
	}
	im.Mutex.Lock()
	defer im.Mutex.Unlock()
	im.challenges[id] = *c
	return nil
}

func (im *InMemory) GetChallenge(id string) (*challenge.Challenge, error) {
	if id == "" {
		return nil, ErrInvalidChallengeId
	}
	im.Mutex.Lock()
	defer im.Mutex.Unlock()
	challenge, ok := im.challenges[id]
	if !ok {
		return nil, ErrChallengeNotFound
	}
	return &challenge, nil
}

func (im *InMemory) DeleteChallenge(id string) error {
	if id == "" {
		return ErrInvalidChallengeId
	}
	im.Mutex.Lock()
	defer im.Mutex.Unlock()
	if _, ok := im.challenges[id]; !ok {
		return ErrChallengeNotFound
	}
	delete(im.challenges, id)
	return nil
}
