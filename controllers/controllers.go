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
			result := BlockIterator(block, p, m)
			if result == -1 {
				continue
			}
			return lib.Move(result)
		}

		return lib.Rely

	}

	return strategie

}

func BlockIterator(block goblockly.Block, p lib.Player, m *map[lib.Player][]lib.Move) int {
	values := block.Values

	results := []int{}

	for _, value := range values {
		for _, block := range value.Blocks {
			results = append(results, BlockIterator(block, p, m))
		}
	}

	switch block.Type {
	case "round":
		return RoundBlock(m)
	case "betray":
		return 0
	case "rely":
		return 1
	case "math_number":
		return MathNumberBlock(block)
	case "plmoves":
		return PlayerMovesBlock(p, m, results)
	case "opmoves":
		return OpponentMovesBlock(p, m, results)
	case "operations":
		return OperationBlock(block, results)
	case "logic_compare":
		return LogicCompareBlock(block, results)
	case "return_move":
		return ReturnMoveBlock(results)
	case "controls_if":
		return IfBlock(block, p, m, results)
	}
	return -1
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

func RoundBlock(m *map[lib.Player][]lib.Move) int {
	return len((*m)[0])
}

func MathNumberBlock(block goblockly.Block) int {
	if value, err := strconv.Atoi(block.Fields[0].Value); err == nil {
		return value
	}
	return 0
}

func LogicCompareBlock(block goblockly.Block, r []int) int {
	r0 := (r)[0]
	r1 := (r)[1]
	switch block.Fields[0].Value {
	case "EQ":
		if r0 == r1 {
			return 1
		} else {
			return 0
		}
	case "NEQ":
		if r0 != r1 {
			return 1
		} else {
			return 0
		}
	case "LT":
		if r0 < r1 {
			return 1
		} else {
			return 0
		}
	case "LTE":
		if r0 <= r1 {
			return 1
		} else {
			return 0
		}
	case "GT":
		if r0 > r1 {
			return 1
		} else {
			return 0
		}
	case "GTE":
		if r0 >= r1 {
			return 1
		} else {
			return 0
		}
	default:
		return 0
	}
}

func OperationBlock(block goblockly.Block, r []int) int {
	r0 := (r)[0]
	r1 := (r)[1]

	if r1 != 0 {
		switch block.Fields[0].Value {
		case "sum":
			return r0 + r1
		case "substraction":
			return r0 - r1
		case "multiplication":
			return r0 * r1
		case "division":
			return r0 / r1
		}
	}

	return 0
}

func PlayerMovesBlock(p lib.Player, m *map[lib.Player][]lib.Move, r []int) int {
	return int((*m)[p][(r)[0]])
}

func OpponentMovesBlock(p lib.Player, m *map[lib.Player][]lib.Move, r []int) int {
	return int((*m)[p.Opposite()][(r)[0]])
}

func ReturnMoveBlock(r []int) int {
	if len((r)) > 0 {
		return (r)[0]
	}
	return -1
}

func IfBlock(block goblockly.Block, p lib.Player, m *map[lib.Player][]lib.Move, r []int) int {
	state := (r)[0]
	statements := block.Statements[0].Blocks

	if state == 1 {
		if statements[0].Type == "return_move" {
			return BlockIterator(statements[0], p, m)
		}
	}

	return -1
}
