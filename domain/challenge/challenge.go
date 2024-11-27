// challenge package holds the challenge entity
//
// A challenge entity refers to a single rock paper and scissor challenge
package challenge

import (
	"errors"

	"github.com/ekefan/discord-bot/domain"
)

// Challenge Errors
var (
	ErrInvalidChallengeID   = errors.New("challenge id must not be empty")
	ErrInvalidPlayer        = errors.New("player must have a valid choice and id")
	ErrOpponentExists       = errors.New("an opponent for this challenge exists")
	ErrChallengerNotOpposed = errors.New("an opponent has not opposed this challenge")
)

// Challenge is an instance of a Rock Paper Scissor Challenge
type Challenge struct {
	id         string
	challenger *domain.Player
	opponent   *domain.Player
	result     *domain.ChallengeResult
}

// TODO: Write tests for these functions

// NewChallenge Factory create new Challenges
func NewChallenge(challengeId string, challenger *domain.Player) (*Challenge, error) {
	if challengeId == "" {
		return nil, ErrInvalidChallengeID
	}
	if challenger == nil || !challenger.Valid() {
		return nil, ErrInvalidPlayer
	}
	return &Challenge{
		id:         challengeId,
		challenger: challenger,
	}, nil
}

func (c *Challenge) GetChallengeID() (string, error) {
	if c.id == "" {
		return "", ErrInvalidChallengeID
	}
	return c.id, nil
}

// SetOpponent sets an opponent for a challenge
// If an opponent doesn't already exists
func (c *Challenge) SetOpponent(opponent *domain.Player) error {
	if c.opponent != nil {
		return ErrOpponentExists
	}
	c.opponent = opponent
	return nil
}

// DeterminChallengeResult determines challenge result from the choice of
// challenger and opponent
func (c *Challenge) DetermineChallengeResult() error {
	if c.opponent == nil {
		return ErrChallengerNotOpposed
	}
	choiceMap := map[domain.RpsChoice]domain.RpsChoice{
		domain.Rock:    domain.Scissor,
		domain.Paper:   domain.Rock,
		domain.Scissor: domain.Paper,
	}
	if c.challenger.Choice == c.opponent.Choice {
		c.result = &domain.ChallengeResult{
			Winner:      c.challenger,
			Looser:      c.opponent,
			OutcomeDraw: true,
		}
	} else if choiceMap[c.challenger.Choice] == c.opponent.Choice {
		c.result = &domain.ChallengeResult{
			Winner: c.challenger,
			Looser: c.opponent,
		}
	} else {
		c.result = &domain.ChallengeResult{
			Winner: c.opponent,
			Looser: c.challenger,
		}
	}
	return nil
}

// GetResultMsg returns the result of a challenge in a formatted message
func (c *Challenge) GetResultMsg() (string, error) {
	resultMsg, err := c.result.FormatResult()
	if err != nil {
		return "", err
	}
	return resultMsg, nil
}
