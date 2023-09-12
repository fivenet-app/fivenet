package postals

import (
	"fmt"
	"io"
	"os"

	"github.com/galexrt/fivenet/pkg/config"
	jsoniter "github.com/json-iterator/go"
	"github.com/paulmach/orb"
	"github.com/paulmach/orb/quadtree"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

type Postals struct {
	tree  *quadtree.Quadtree
	codes map[string]*Postal
}

func New(cfg *config.Config) (*Postals, error) {
	file, err := os.Open(cfg.Game.Livemap.PostalsFile)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	buf, _ := io.ReadAll(file)

	var codes []*Postal
	if err := json.Unmarshal(buf, &codes); err != nil {
		return nil, err
	}

	codesMap := map[string]*Postal{}
	tree := quadtree.New(orb.Bound{Min: orb.Point{-4000, -4000}, Max: orb.Point{8000, 6000}})
	for _, code := range codes {
		point := orb.Point{code.X, code.Y}
		tree.Add(point)
		codesMap[getCoordsKey(point)] = code
	}

	return &Postals{
		tree:  tree,
		codes: codesMap,
	}, nil
}

func (p *Postals) Closest(x, y float64) *Postal {
	nearest := p.tree.Find(orb.Point{x, y})
	postal, ok := p.codes[getCoordsKey(nearest.Point())]
	if !ok {
		return nil
	}

	return postal
}

func getCoordsKey(point orb.Point) string {
	return fmt.Sprintf("%f-%f", point.X(), point.Y())
}
