//
// Code generated by go-jet DO NOT EDIT.
//
// WARNING: Changes to this file may cause incorrect behavior
// and will be lost if the code is regenerated
//

package model

import (
	"time"
)

type GksphoneJobMessage struct {
	ID      int32     `sql:"primary_key" json:"id"`
	Name    *string   `json:"name"`
	Number  *string   `json:"number"`
	Message *string   `json:"message"`
	Photo   *string   `json:"photo"`
	Gps     *string   `json:"gps"`
	Owner   int32     `json:"owner"`
	Jobm    *string   `json:"jobm"`
	Anon    *string   `json:"anon"`
	Time    time.Time `json:"time"`
}
