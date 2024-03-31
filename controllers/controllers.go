package controllers

import (
	"encoding/json"
	"fmt"
	"prisioner-game/lib"
	"prisioner-game/strategies"

	"github.com/gofiber/fiber/v2"
)

type Block string

const (
	EQUAL Block = "EQUAL"
)

type RoundBody struct {
	Size          uint
	Strategie_one []string
	Strategie_two []string
}

type XMLStruct struct {
	Xml string
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

func XML(c *fiber.Ctx) error {
	//var body XMLStruct
	//c.BodyParser(&body)
	/* var blocks goblockly.BlockXml
	if err := xml.Unmarshal([]byte(c.Query("xml")), &blocks); err != nil {
		return c.Status(400).SendString("Bad Request")
	} */
	fmt.Println(c.Query("xml"))
	return c.JSON(XMLStruct{Xml: c.Query("xml")})

}

func Equal(s1 string, s2 string, p lib.Player, m *map[lib.Player][]lib.Move) bool {
	switch s1 {
	case "ROUND":
		s1 = string(len((*m)[0]))
	}

	switch s2 {
	case "ROUND":
		s1 = string(len((*m)[0]))
	}

	return s1 == s2
}
