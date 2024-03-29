package controllers

import (
	"encoding/json"
	"prisioner-game/lib"
	"prisioner-game/strategies"

	"github.com/gofiber/fiber/v2"
)

type RoundBody struct {
	Size          uint
	Strategie_one []string
	Strategie_two []string
}

type RoundResult struct {
	Moves  map[lib.Player][]lib.Move
	Scores map[lib.Player][]int
}

func Root(c *fiber.Ctx) error {
	return c.SendString("Root")
}

func GetRound(c *fiber.Ctx) error {
	var body RoundBody
	err := json.Unmarshal(c.Body(), &body)
	if err != nil {
		return fiber.ErrBadRequest
	}

	moves_record := make(map[lib.Player][]lib.Move)
	scores_record := make(map[lib.Player][]int)

	for range body.Size {
		p, s := lib.Round(strategies.Random, strategies.TitForTat, &moves_record)

		if p == lib.Both {
			scores_record[lib.First] = append(scores_record[lib.First], s)
			scores_record[lib.Second] = append(scores_record[lib.Second], s)
		} else {
			scores_record[p] = append(scores_record[p], s)
			scores_record[p.Opposite()] = append(scores_record[p.Opposite()], 0)
		}
	}

	response := RoundResult{moves_record, scores_record}

	return c.JSON(response)
}
