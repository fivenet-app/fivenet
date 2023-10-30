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
	tree *quadtree.Quadtree
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

	tree := quadtree.New(orb.Bound{Min: orb.Point{-9_000, -9_000}, Max: orb.Point{11_000, 11_000}})
	for k := range codes {
		if err := tree.Add(codes[k]); err != nil {
			return nil, fmt.Errorf("failed to add postal to postals map: %w", err)
		}
	}

	return &Postals{
		tree: tree,
	}, nil
}

func (p *Postals) Closest(x, y float64) *Postal {
	nearest := p.tree.Find(orb.Point{x, y})

	return nearest.(*Postal)
}
