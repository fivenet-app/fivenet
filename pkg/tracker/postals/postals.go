package postals

type Postals struct {
}

func New() (*Postals, error) {
	// TODO load postals.json file

	return &Postals{}, nil
}

func (p *Postals) ClosestPostal(x, y float64) (int64, bool) {

	// TODO

	return 0, false
}
