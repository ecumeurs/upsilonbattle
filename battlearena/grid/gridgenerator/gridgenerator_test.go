package gridgenerator

import (
	"os"
	"testing"

	"github.com/ecumeurs/upsilontools/tools"
)

func TestGridGeneratorFlat(t *testing.T) {
	gg := GridGenerator{}
	gg.Width = tools.NewIntRange(20, 50)
	gg.Length = tools.NewIntRange(20, 50)
	gg.Height = tools.NewIntRange(10, 15)
	gg.GenerateObstrcution = true
	gg.Type = Flat
	gg.ObstructionRate = tools.NewIntRange(10, 50)

	gr := gg.Generate()
	res := gr.GenerateHTML()
	// store res to file result.svg
	os.WriteFile("resultFlat.html", []byte(res), 0644)
}

func TestGridGeneratorHill(t *testing.T) {
	gg := GridGenerator{}
	gg.Width = tools.NewIntRange(20, 50)
	gg.Length = tools.NewIntRange(20, 50)
	gg.Height = tools.NewIntRange(10, 15)
	gg.GenerateObstrcution = true
	gg.Type = Hill
	gg.ObstructionRate = tools.NewIntRange(10, 50)

	gr := gg.Generate()
	res := gr.GenerateHTML()
	// store res to file result.svg
	os.WriteFile("resultHill.html", []byte(res), 0644)
}

func TestGridGeneratorRiver(t *testing.T) {
	gg := GridGenerator{}
	gg.Width = tools.NewIntRange(20, 50)
	gg.Length = tools.NewIntRange(20, 50)
	gg.Height = tools.NewIntRange(10, 15)
	gg.GenerateObstrcution = true
	gg.Type = River
	gg.ObstructionRate = tools.NewIntRange(10, 50)

	gr := gg.Generate()
	res := gr.GenerateHTML()
	// store res to file result.svg
	os.WriteFile("resultRiver.html", []byte(res), 0644)
}
