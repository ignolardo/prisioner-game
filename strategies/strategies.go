package strategies

import (
	"math/rand"
	. "prisioner-game/lib"
)

func Random(p Player, r *map[Player][]Move) Move {
	return Move(rand.Intn(2))
}

func AlwaysBetray(p Player, r *map[Player][]Move) Move {
	return Betray
}

func AlwaysRely(p Player, r *map[Player][]Move) Move {
	return Rely
}

func TitForTat(p Player, r *map[Player][]Move) Move {

	oppositeMoves := (*r)[p.Opposite()]

	if len(oppositeMoves) == 0 {
		//return Move(rand.Intn(2))
		return Rely
	}

	return oppositeMoves[len(oppositeMoves)-1]

}

func George(p Player, r *map[Player][]Move) Move {

	oppositeMoves := (*r)[p.Opposite()]
	playerMoves := (*r)[p]

	switch {
	case len(playerMoves) == 0:
		return Rely
	case playerMoves[len(playerMoves)-1] == Betray:
		return Betray
	case oppositeMoves[len(oppositeMoves)-1] == Betray:
		return Betray
	default:
		return Rely
	}

}
