package visual

import (
	"fmt"
	. "prisioner-game/lib"
	"prisioner-game/strategies"
	"strconv"
)

type Color string

const (
	Red    Color = "\033[31m"
	Green  Color = "\033[32m"
	Blue   Color = "\033[34m"
	Yellow Color = "\033[33m"
	Purple Color = "\033[35m"
	Cyan   Color = "\033[36m"
	Gray   Color = "\033[37m"
	White  Color = "\033[97m"
)

func main() {

	// Define players strategies

	var s1 Strategy = strategies.Random

	var s2 Strategy = strategies.TitForTat

	// Number of rounds

	rounds_number := 50

	// Record and Score lists

	record := make(map[Player][]Move)

	score := make(map[Player][]int)

	// Start playing rounds and recording the results

	for i := 0; i < rounds_number; i++ {

		player, points := Round(s1, s2, &record)

		if player == 2 {
			score[First] = append(score[First], points)
			score[Second] = append(score[Second], points)
		} else {
			score[player] = append(score[player], points)
			score[player.Opposite()] = append(score[player.Opposite()], 0)
		}
	}

	// Printing the results in terminal

	fmt.Println("\nMoves record")
	fmt.Println("")

	for n := range rounds_number {
		len_loop := func(i int) int {
			if i == 0 {
				return 1
			}
			count := 0
			for i != 0 {
				i /= 10
				count++
			}
			return count
		}
		rn := color(strconv.Itoa(n+1), Blue)
		zeros := len_loop(rounds_number) - len_loop(n+1)
		for range zeros {
			rn += " "
		}
		fmt.Printf("%s (%d) %s   %s (%d)\n", rn, score[0][n], record[0][n].Symbol(), record[1][n].Symbol(), score[1][n])
	}

	sum1 := 0

	for _, points := range score[0] {
		sum1 += points
	}

	sum2 := 0

	for _, points := range score[1] {
		sum2 += points
	}

	fmt.Println("\nPlayer 1 score:", sum1)
	fmt.Println("\nPlayer 2 score:", sum2)
}

func color(s string, c Color) string {

	var color_str string

	color_str += string(c)
	color_str += s
	color_str += "\033[0m"

	return color_str
}
