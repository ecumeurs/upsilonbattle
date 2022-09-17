package gridgenerator

import (
	"os"
	"testing"

	"github.com/ecumeurs/upsilonbattle/tools"
)

func TestGridGeneratorFlat(t *testing.T) {
	gg := GridGenerator{}
	gg.Width = tools.NewIntRange(20, 50, 1)
	gg.Length = tools.NewIntRange(20, 50, 1)
	gg.Height = tools.NewIntRange(10, 15, 1)
	gg.GenerateObstrcution = true
	gg.Type = Flat
	gg.ObstructionRate = tools.NewIntRange(10, 50, 1)

	gr := gg.Generate()
	res := gr.GenerateHTML()
	// store res to file result.svg
	os.WriteFile("result.html", []byte(res), 0644)
	t.Fail()
}

func TestGridGeneratorHill(t *testing.T) {
	gg := GridGenerator{}
	gg.Width = tools.NewIntRange(20, 50, 1)
	gg.Length = tools.NewIntRange(20, 50, 1)
	gg.Height = tools.NewIntRange(10, 15, 1)
	gg.GenerateObstrcution = true
	gg.Type = Hill
	gg.ObstructionRate = tools.NewIntRange(10, 50, 1)

	gr := gg.Generate()
	res := gr.GenerateHTML()
	// store res to file result.svg
	os.WriteFile("result.html", []byte(res), 0644)
	t.Fail()
}

func TestGridGeneratorRiver(t *testing.T) {
	gg := GridGenerator{}
	gg.Width = tools.NewIntRange(20, 50, 1)
	gg.Length = tools.NewIntRange(20, 50, 1)
	gg.Height = tools.NewIntRange(10, 15, 1)
	gg.GenerateObstrcution = true
	gg.Type = River
	gg.ObstructionRate = tools.NewIntRange(10, 50, 1)

	gr := gg.Generate()
	res := gr.GenerateHTML()
	// store res to file result.svg
	os.WriteFile("result.html", []byte(res), 0644)
	t.Fail()
}
