package dispatch

import "dario.cat/mergo"

func (x *Unit) Update(in *Unit) {
	if x.Id != in.Id {
		return
	}

	err := mergo.Merge(x, in)
	if err != nil {
		return
	}
}
