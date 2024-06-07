package repository

type Repository[T any] interface {
	Create(record T) (T, error)
	Read(uuid string) (T, error)
	Update(uuid string, newRecord T) (T, error)
	Delete(uuid string) error
	Exists(uuid string) bool
}

type Identifiable interface {
	GetUuid() string
}
