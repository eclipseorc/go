package models

type User struct {
	Id   int64  `json:id, xorm:"id"`
	Name string `json:name, xorm:"name"`
}

func (u *User) TableName() string {
	return "user"
}

func (u *User) Add(user User) error {
	Db(0).
	return nil
}

func (u *User) List() error  {
	return nil
}