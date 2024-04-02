package controllers

import (
	"encoding/xml"
	"fmt"
	"prisioner-game/lib"
	"prisioner-game/strategies"
	"strconv"

	"github.com/fixermark/goblockly"
	"github.com/gofiber/fiber/v2"
)

type RoundBody struct {
	Size uint
	Xml  string
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
	/* var body RoundBody
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

	return c.JSON(response) */

	var body RoundBody
	if err := c.BodyParser(&body); err != nil {
		c.Status(400).SendString(err.Error())
	}

	var blocks goblockly.BlockXml
	if err := xml.Unmarshal([]byte(body.Xml), &blocks); err != nil {
		return c.Status(400).SendString(err.Error())
	}
	var strategie lib.Strategy

	if len(blocks.Blocks) == 0 {
		strategie = strategies.Random
	} else {
		strategie = ParseStrategie(GetBlocks(blocks.Blocks[0]))
	}

	moves_record := make(map[lib.Player][]lib.Move)
	scores_record := make(map[lib.Player][]int)

	for range body.Size {
		p, s := lib.Round(strategie, strategies.Random, &moves_record)

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
	var xmlBody XMLStruct
	if err := c.BodyParser(&xmlBody); err != nil {
		c.Status(400).SendString(err.Error())
	}

	var blocks goblockly.BlockXml
	if err := xml.Unmarshal([]byte(xmlBody.Xml), &blocks); err != nil {
		return c.Status(400).SendString(err.Error())
	}

	PrintBlocks(blocks.Blocks)

	//strategie := ParseStrategie(blocks.Blocks)

	return c.XML(blocks)
}

func GetBlocks(firstBlock goblockly.Block) []goblockly.Block {
	blocks := []goblockly.Block{firstBlock}
	lastBlock := firstBlock

	for lastBlock.Next != nil {
		blocks = append(blocks, *lastBlock.Next)
		lastBlock = *lastBlock.Next
	}

	return blocks
}

func ParseStrategie(blocks []goblockly.Block) lib.Strategy {

	strategie := func(p lib.Player, m *map[lib.Player][]lib.Move) lib.Move {
		for _, block := range blocks {
			if block.Type == "controls_if" {
				if IfBlock(block, p, m) {
					fmt.Println("IF YES")
					if block.Statements[0].Blocks[0].Type == "return_move" {
						switch BlockIterator(block.Statements[0].Blocks[0], p, m) {
						case 0:
							return lib.Betray
						case 1:
							return lib.Rely
						}
					}
				} else {
					fmt.Println("IF NOT")
				}
			}
			if block.Type == "return_move" {
				switch BlockIterator(block, p, m) {
				case 0:
					return lib.Betray
				case 1:
					return lib.Rely
				}
			}
		}

		return lib.Rely

	}

	return strategie

}

func IfBlock(block goblockly.Block, p lib.Player, m *map[lib.Player][]lib.Move) bool {
	value_blocks := block.Values[0].Blocks
	results := []bool{}

	switch BlockIterator(value_blocks[0], p, m) {
	case 1:
		results = append(results, true)
	case 0:
		results = append(results, false)
	}

	fmt.Println(results)

	for _, b := range results {
		if !b {
			return false
		}
	}

	return true

}

func BlockIterator(block goblockly.Block, p lib.Player, m *map[lib.Player][]lib.Move) int {
	values := block.Values

	if len(values) == 0 {
		return BasicBlocks(block, p, m)
	}

	results := []int{}
	for _, value := range values {
		for _, block := range value.Blocks {
			results = append(results, BlockIterator(block, p, m))
		}
	}

	////////////////////////////////////////////////////////////
	// MAKE ALL THE OPTIONS NOT ONLY EQUAL OPERATION
	////////////////////////////////////////////////////////////
	if block.Type == "logic_compare" {
		if results[0] == results[1] {
			fmt.Println("BOTH SIDES ARE EQUAL")
			return 1
		} else {
			return 0
		}
	}

	if block.Type == "operations" {
		switch block.Fields[0].Value {
		case "sum":
			return results[0] + results[1]
		case "substraction":
			return results[0] - results[1]
		case "multiplication":
			return results[0] * results[1]
		case "division":
			return results[0] / results[1]
		}
	}

	if block.Type == "return_move" {
		return results[0]
	}

	if block.Type == "plmoves" {
		return int((*m)[p][results[0]])
	}

	if block.Type == "opmoves" {
		return int((*m)[p.Opposite()][results[0]])
	}

	return results[0]
}

func LogicCompare(block goblockly.Block) bool {
	/* a_value := block.Values[0].Blocks
	b_value := block.Values[0].Blocks

	for _,a_block := range a_value {

	} */

	return true
}

func BasicBlocks(block goblockly.Block, p lib.Player, m *map[lib.Player][]lib.Move) int {
	if block.Type == "round" {
		return len((*m)[0])
	}
	if block.Type == "betray" {
		//fmt.Println("BETRAY BLOCK FOUNDED")
		return 0
	}
	if block.Type == "rely" {
		//fmt.Println("RELY BLOCK FOUNDED")
		return 1
	}
	if block.Type == "math_number" {
		if value, err := strconv.Atoi(block.Fields[0].Value); err == nil {
			return value
		}
	}
	return 0
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
