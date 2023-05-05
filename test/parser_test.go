package test

import (
    "strings"
    "testing"

	"github.com/Illyakravchuk/lab_3/painter"
    "github.com/Illyakravchuk/lab_3/painter/lang"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

func TestParseFuncWhiteCommand(t *testing.T) {
    parser := &lang.Parser{}
    command := "white\nupdate"
    op := painter.OperationFunc(painter.WhiteFill)

    ops, err := parser.Parse(strings.NewReader(command))

    require.NoError(t, err)
    assert.IsType(t, op, ops[0])
}

func TestParseFuncGreenCommand(t *testing.T) {
    parser := &lang.Parser{}
    command := "green\nupdate"
    op := painter.OperationFunc(painter.GreenFill)

    ops, err := parser.Parse(strings.NewReader(command))

    require.NoError(t, err)
    assert.IsType(t, op, ops[0])
}

func TestParseFuncResetCommand(t *testing.T) {
    parser := &lang.Parser{}
    command := "reset\nupdate"
    op := painter.OperationFunc(painter.ResetScreen)

    ops, err := parser.Parse(strings.NewReader(command))

    require.NoError(t, err)
    assert.IsType(t, op, ops[0])
}

func TestParseStructBgRectCommand(t *testing.T) {
    parser := &lang.Parser{}
    command := "bgrect 0.25 0.25 0.75 0.75\nupdate "
    op := &painter.BgRectangle{X1: 200, Y1: 200, X2: 600, Y2: 600}

    ops, err := parser.Parse(strings.NewReader(command))

    require.NoError(t, err)
    assert.IsType(t, op, ops[1])
    assert.Equal(t, op, ops[1])
}

func TestParseStructFigureCommand(t *testing.T) {
    parser := &lang.Parser{}
    command := "figure 0.5 0.5\nupdate"
    op := &painter.Figure{X: 400, Y: 400}

    ops, err := parser.Parse(strings.NewReader(command))

    require.NoError(t, err)
    assert.IsType(t, op, ops[1])
    assert.Equal(t, op, ops[1])
}

func TestParseStructMoveCommand(t *testing.T) {
    parser := &lang.Parser{}
    command := "move 0.3 0.3\nupdate"
    op := &painter.Move{X: 240, Y: 240, Figures: []*painter.Figure(nil)}

    ops, err := parser.Parse(strings.NewReader(command))

    require.NoError(t, err)
    assert.IsType(t, op, ops[1])
    assert.Equal(t, op, ops[1])
}


func TestParseStructInvalidCommand(t *testing.T) {
    parser := &lang.Parser{}
    command := "invalid command"
    
    _, err := parser.Parse(strings.NewReader(command))

    assert.Error(t, err)
}
