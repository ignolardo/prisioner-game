package lib

type Move int

func (self *Move) Symbol() string {
	switch *self {
	case Betray:
		return "\033[31m☒\033[0m"
	case Rely:
		return "\033[32m🗹\033[0m"
	default:
		return "_"
	}
}

type Player int

const (
	First Player = iota
	Second
	Both
)

func (self *Player) Opposite() Player {
	switch *self {
	case First:
		return Second
	case Second:
		return First
	case Both:
		return Both
	default:
		return *self
	}
}

const (
	Betray Move = iota
	Rely
)

type Strategy func(Player, *map[Player][]Move) Move

func Round(s1 Strategy, s2 Strategy, r *map[Player][]Move) (Player, int) {
	move1 := s1(First, r)
	move2 := s2(Second, r)

	(*r)[First] = append((*r)[First], move1)
	(*r)[Second] = append((*r)[Second], move2)

	switch {
	case move1 == Betray && move2 == Rely:
		return First, 5
	case move1 == Rely && move2 == Betray:
		return Second, 5
	case move1 == Betray && move2 == Betray:
		return Both, 1
	case move1 == Rely && move2 == Rely:
		return Both, 3
	default:
		return Both, 0
	}
}
