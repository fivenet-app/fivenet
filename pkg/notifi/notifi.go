package notifi

import "database/sql"

type INotifi interface {
	Add(userId int32, title string, content string, ty string)
}

type Notifi struct {
	db *sql.DB
}

func New(db *sql.DB) *Notifi {
	return &Notifi{
		db: db,
	}
}

func (n *Notifi) Add(userId int32, title string, content string, ty string) {

}
