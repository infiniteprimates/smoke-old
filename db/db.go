package db

type Db struct {
	user
}

func New() *Db {
	return &Db{
		user: user{},
	}
}
