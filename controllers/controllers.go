package controllers

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"prisioner-game/lib"
	"prisioner-game/strategies"

	"github.com/fixermark/goblockly"
	"github.com/gofiber/fiber/v2"
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
	var xmlBody XMLStruct
	if err := c.BodyParser(&xmlBody); err != nil {
		c.Status(400).SendString(err.Error())
	}
	var blocks goblockly.BlockXml
	if err := xml.Unmarshal([]byte(xmlBody.Xml), &blocks); err != nil {
		return c.Status(400).SendString(err.Error())
	}
	//fmt.Println(c.Query("xml"))
	PrintBlocks(blocks.Blocks)
	//return c.JSON(XMLStruct{Xml: c.Query("xml")})
	return c.XML(blocks)
}

func PrintBlocks(blocks []goblockly.Block) {
	fmt.Println("----------------------------")
	for _, block := range blocks {
		fmt.Printf("Type: %s Id: %s\n", block.Type, block.Id)
		PrintNextBlock(block)
	}
}

func PrintNextBlock(block goblockly.Block) {
	if block.Next == nil {
		fmt.Printf("Block %s has no next blocks\n", block.Id)
	} else {
		fmt.Printf("Next Of %s : %s\n", block.Id, block.Next.Type)
		PrintNextBlock(*block.Next)
	}
}
