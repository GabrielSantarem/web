package adapters

import (
	"errors"
	"sync"

	"codeberg.org/MrTomate/web/internal/core"
)

var (
	ErrEntityNotFound = errors.New("entidade não encontrada")
	ErrEntityNil      = errors.New("entidade não pode ser nula")
)

// MemoryRepository implementa um repositório genérico em memória (adaptador de saída)
// seguro para concorrência utilizando sync.RWMutex.
type MemoryRepository[T core.Entity] struct {
	mu    sync.RWMutex
	items map[string]T
}

// NewMemoryRepository cria uma nova instância de repositório genérico.
func NewMemoryRepository[T core.Entity]() *MemoryRepository[T] {
	return &MemoryRepository[T]{
		items: make(map[string]T),
	}
}

// Create armazena uma nova entidade.
func (r *MemoryRepository[T]) Create(entity T) (T, error) {
	var zero T
	if any(entity) == nil {
		return zero, ErrEntityNil
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	r.items[entity.GetID()] = entity
	return entity, nil
}

// Update substitui os dados de uma entidade existente.
func (r *MemoryRepository[T]) Update(entity T) (T, error) {
	var zero T
	if any(entity) == nil {
		return zero, ErrEntityNil
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.items[entity.GetID()]; !exists {
		return zero, ErrEntityNotFound
	}

	r.items[entity.GetID()] = entity
	return entity, nil
}

// FindByID busca uma entidade pela chave identificadora.
func (r *MemoryRepository[T]) FindByID(id string) (T, error) {
	var zero T
	r.mu.RLock()
	defer r.mu.RUnlock()

	item, exists := r.items[id]
	if !exists {
		return zero, ErrEntityNotFound
	}

	return item, nil
}

// FindAll retorna todas as entidades armazenadas.
func (r *MemoryRepository[T]) FindAll() ([]T, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	all := make([]T, 0, len(r.items))
	for _, item := range r.items {
		all = append(all, item)
	}
	return all, nil
}

// Delete remove uma entidade pelo ID.
func (r *MemoryRepository[T]) Delete(id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.items[id]; !exists {
		return ErrEntityNotFound
	}

	delete(r.items, id)
	return nil
}
