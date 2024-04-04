package controllers

import (
	"encoding/xml"
	"errors"
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

func GetRounds(c *fiber.Ctx) error {
	// BODY STRUCT VARIABLE
	var body RoundBody

	// PARSE BODY REQUEST INTO BODY VARIABLE
	if err := c.BodyParser(&body); err != nil {
		c.Status(400).SendString(err.Error())
	}

	// BLOCKS STRUCT VARIABLE
	blockXml, err := ParseXML(body.Xml)

	if err != nil {
		return c.Status(400).SendString(err.Error())
	}

	// PARSE BLOCKS INTO STRATEGIE
	strategie := ParseStrategie(blockXml.Blocks)

	// PLAYERS MOVES RECORD
	moves_record := make(map[lib.Player][]lib.Move)

	// PLAYERS SCORES
	scores_record := make(map[lib.Player][]int)

	// GENERATE THE REQUESTED ROUNDS
	lib.MultipleRounds(body.Size, strategie, strategies.Random, &moves_record, &scores_record)

	// RESPONSE VARIABLE
	response := RoundResult{moves_record, scores_record}

	// RETURNING JSON RESPONSE
	return c.JSON(response)

}

// XML PARSER
func ParseXML(strXml string) (goblockly.BlockXml, error) {
	// BLOCKS STRUCT VARIABLE
	var blocks goblockly.BlockXml

	// PARSE XML INTO BLOCKS VARIABLE
	if err := xml.Unmarshal([]byte(strXml), &blocks); err != nil {
		return goblockly.BlockXml{}, errors.New("cannot parse xml")
	}

	return blocks, nil
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

	if len(blocks) == 0 {
		return strategies.Random
	}

	blocks = GetBlocks(blocks[0])

	strategie := func(p lib.Player, m *map[lib.Player][]lib.Move) lib.Move {
		for _, block := range blocks {
			if block.Type == "controls_if" && len(block.Values[0].Blocks) > 0 && len(block.Statements[0].Blocks) > 0 {
				if IfBlock(block, p, m) {
					//fmt.Println("IF BLOCK IS CORRECT")
					if block.Statements[0].Blocks[0].Type == "return_move" {
						switch BlockIterator(block.Statements[0].Blocks[0], p, m) {
						case 0:
							return lib.Betray
						case 1:
							return lib.Rely
						}
					}
				}
			}
			if block.Type == "return_move" && len(block.Values[0].Blocks) > 0 {
				//fmt.Println("RETURN BLOCK DETECTED")
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
		//fmt.Println("LOGIC COMPARE DETECTED")
		switch block.Fields[0].Value {
		case "EQ":
			if results[0] == results[1] {
				//fmt.Println("BOTH SIDES ARE EQUAL")
				return 1
			} else {
				return 0
			}
		case "NEQ":
			if results[0] != results[1] {
				//fmt.Println("BOTH SIDES ARE NOT EQUAL")
				return 1
			} else {
				return 0
			}
		case "LT":
			if results[0] < results[1] {
				//fmt.Println("LEFT SIDE IS LESS THAN RIGHT SIDE")
				return 1
			} else {
				return 0
			}
		case "LTE":
			if results[0] <= results[1] {
				//fmt.Println("LEFT SIDE IS LESS OR EQUAL THAN RIGHT SIDE")
				return 1
			} else {
				return 0
			}
		case "GT":
			if results[0] > results[1] {
				//fmt.Println("LEFT SIDE IS GREATER THAN RIGHT SIDE")
				return 1
			} else {
				return 0
			}
		case "GTE":
			if results[0] >= results[1] {
				//fmt.Println("LEFT SIDE IS GREATER OR EQUAL THAN RIGHT SIDE")
				return 1
			} else {
				return 0
			}
		default:
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
