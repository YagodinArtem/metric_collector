package repository

type Repository[TYPE any, ID any] interface {
	Save(*TYPE) (bool, error)

	SaveAll([]TYPE) error

	FindByID(ID) (*TYPE, error)

	FindAll() ([]TYPE, error)

	Update(*TYPE) (bool, error)

	DeleteByID(ID) (bool, error)

	DeleteAll() (bool, error)
}
