package domain

type Player struct {
	ID     string    `json:"id"`
	Choice RpsChoice `json:"choice"`
}

// Valid returns false when player id is empty or choice is not either rock, paper or scissor
func (p *Player) Valid() bool {
	if p == nil || p.ID == "" {
		return false
	}
	switch p.Choice {
	case Rock:
		return true
	case Paper:
		return true
	case Scissor:
		return true
	default:
		return false
	}
}
