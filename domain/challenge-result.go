package domain

import (
	"errors"
	"fmt"
)

var (
	ErrInvalidChallengeResult = errors.New("a challenge result must have a valid winner looser and result choice")
)

type ChallengeResult struct {
	Winner      *Player
	Looser      *Player
	OutcomeDraw bool
}

//TODO: write test for formatResultMsg

// FormatResultMsg returns a formatted msg detailing the result of a challenge
// when no winner or looser has been set an invalid challenge result is returned
func (cr *ChallengeResult) FormatResult() (string, error) {
	if cr.Winner == nil || cr.Looser == nil {
		return "", ErrInvalidChallengeResult
	}

	if cr.OutcomeDraw {
		return fmt.Sprintf("<@%v> and <@%v> draw with **%v**", cr.Winner.ID, cr.Looser.ID, cr.Looser.Choice), nil
	}
	return fmt.Sprintf("<@%v> wins the challenge with **%v** beating <@%s>'s **%v**", cr.Winner.ID, cr.Winner.Choice, cr.Looser.ID, cr.Looser.Choice), nil
}
