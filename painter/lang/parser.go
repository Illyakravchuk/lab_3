package lang

import (
	"bufio"
	"errors"
	"io"
	"strconv"
	"strings"

	"github.com/Illyakravchuk/lab_3/painter"
)

// Parser уміє прочитати дані з вхідного io.Reader та повернути список операцій представлені вхідним скриптом.
type Parser struct {
	bgRect  *painter.BgRectangle
	backOp  painter.Operation
	move    []painter.Operation
	figures []*painter.Figure
	res     []painter.Operation
	update  painter.Operation
	updated bool
}

func (p *Parser) Parse(in io.Reader) ([]painter.Operation, error) {
	p.update = nil
	p.res = nil
	if p.backOp == nil {
		p.backOp = painter.OperationFunc(painter.WhiteFill)
	}

	scanner := bufio.NewScanner(in)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		commandLine := scanner.Text()
		if len(commandLine) == 0 {
			continue
		}
		err := p.parse(commandLine) 
		if err != nil {
			return nil, err
		}
		if p.updated {
			p.updated = false
			p.res = p.genOperations()
		}
	}

	// Append remaining operations to result
	if p.backOp != nil {
		p.res = append(p.res, p.backOp)
	}
	if p.bgRect != nil {
		p.res = append(p.res, p.bgRect)
	}
	if p.move != nil {
		p.res = append(p.res, p.move...)
	}
	for _, figure := range p.figures {
		p.res = append(p.res, figure)
	}

	return p.res, scanner.Err()
}

func (p *Parser) genOperations() []painter.Operation {
	var res []painter.Operation

	if p.backOp != nil {
		res = append(res, p.backOp)
	}
	if p.bgRect != nil {
		res = append(res, p.bgRect)
	}
	if p.move != nil {
		res = append(res, p.move...)
		p.move = nil
	}
	for _, figure := range p.figures {
		res = append(res, figure)
	}
	if p.update != nil {
		res = append(res, p.update)
		p.update = nil
	}

	return res
}

func (p *Parser) parse(commandLine string) error {
	fields := strings.Fields(commandLine)
	operation := fields[0]
	var args []int

	for i := 1; i < len(fields); i++ {
		arg, err := strconv.ParseFloat(fields[i], 64)
		if err != nil {
			return err
		}
		if arg > 0 && arg < 1 {
			arg = arg * 800.0
		}
		args = append(args, int(arg))
	}
	switch operation {
	case "white":
		p.backOp = painter.OperationFunc(painter.WhiteFill)
	case "green":
		p.backOp = painter.OperationFunc(painter.GreenFill)
	case "update":
		p.updated = true
		p.update = painter.UpdateOp
	case "bgrect":
		p.bgRect = &painter.BgRectangle{X1: args[0], Y1: args[1], X2: args[2], Y2: args[3]}
	case "figure":
		figure := &painter.Figure{X: args[0], Y: args[1]}
		p.figures = append(p.figures, figure)
	case "move":
		moveOp := &painter.Move{X: args[0], Y: args[1], Figures: p.figures}
		p.move = moveOp
	case "reset":
		p.figures = nil
		p.bgRect = nil
		p.move = nil
		p.backOp = painter.OperationFunc(painter.ResetScreen)
	default:
		return errors.New("Failed")
	}
	return nil
}