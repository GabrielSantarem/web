package core

// Entity define a restrição genérica para qualquer entidade de domínio
// que possua um identificador único.
type Entity interface {
	GetID() string
}

// Repository define o contrato genérico base de persistência (porta de saída),
// permitindo operações CRUD independentes de infraestrutura ou modelo de banco.
type Repository[T Entity] interface {
	Create(entity T) (T, error)
	Update(entity T) (T, error)
	FindByID(id string) (T, error)
	FindAll() ([]T, error)
	Delete(id string) error
}
