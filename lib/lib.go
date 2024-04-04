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

/* type Blocks string

const (
	Something Blocks = "Hello"
) */

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

func SingleRound(s1 Strategy, s2 Strategy, r *map[Player][]Move, s *map[Player][]int) {
	move1 := s1(First, r)
	move2 := s2(Second, r)

	(*r)[First] = append((*r)[First], move1)
	(*r)[Second] = append((*r)[Second], move2)

	switch {
	case move1 == Betray && move2 == Rely:
		(*s)[First] = append((*s)[First], 5)
		(*s)[Second] = append((*s)[Second], 0)
	case move1 == Rely && move2 == Betray:
		(*s)[First] = append((*s)[First], 0)
		(*s)[Second] = append((*s)[Second], 5)
	case move1 == Betray && move2 == Betray:
		(*s)[First] = append((*s)[First], 1)
		(*s)[Second] = append((*s)[Second], 1)
	case move1 == Rely && move2 == Rely:
		(*s)[First] = append((*s)[First], 3)
		(*s)[Second] = append((*s)[Second], 3)
	}
}

func MultipleRounds(rounds uint, s1 Strategy, s2 Strategy, r *map[Player][]Move, s *map[Player][]int) {
	for range rounds {
		SingleRound(s1, s2, r, s)
	}
}
