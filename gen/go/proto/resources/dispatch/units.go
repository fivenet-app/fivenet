package dispatch

import "dario.cat/mergo"

func (x *Unit) Update(in *Unit) {
	if x.Id != in.Id {
		return
	}

	if err := mergo.Merge(x, in); err != nil {
		return
	}
}
