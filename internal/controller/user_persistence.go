package controller

type PersistUser interface {
	Persist(user *CreateUserE) (string, error)
}

type ReadUser interface {
	Read(uuid string) (*ReadUserE, error)
}

type UserPersistence interface {
	PersistUser
	ReadUser
}

type CreateUserE struct {
	firstName string
	lastName  string
	age       int
	address   string
	email     string
}

type ReadUserE struct {
	id        string
	firstName string
	lastName  string
	age       int
	address   string
	email     string
}
