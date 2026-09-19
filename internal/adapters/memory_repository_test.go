package adapters_test

import (
	"fmt"
	"sync"
	"testing"

	"codeberg.org/MrTomate/web/internal/adapters"
	"codeberg.org/MrTomate/web/internal/core"
)

func TestMemoryRepository(t *testing.T) {
	repo := adapters.NewMemoryRepository[*core.User]()

	t.Run("cria e busca entidade por ID", func(t *testing.T) {
		u := &core.User{ID: "u-1", Name: "Alice", Email: "alice@example.com"}
		created, err := repo.Create(u)
		if err != nil {
			t.Fatalf("esperava sucesso na criacao, erro: %v", err)
		}
		if created.GetID() != "u-1" {
			t.Errorf("ID esperado u-1, veio %s", created.GetID())
		}

		found, err := repo.FindByID("u-1")
		if err != nil {
			t.Fatalf("esperava encontrar entidade, erro: %v", err)
		}
		if found.Name != "Alice" {
			t.Errorf("esperava nome Alice, veio %s", found.Name)
		}
	})

	t.Run("busca de ID inexistente retorna erro", func(t *testing.T) {
		_, err := repo.FindByID("inexistente")
		if err != adapters.ErrEntityNotFound {
			t.Fatalf("esperava ErrEntityNotFound, veio %v", err)
		}
	})

	t.Run("atualiza entidade existente", func(t *testing.T) {
		u := &core.User{ID: "u-1", Name: "Alice Silva", Email: "alice@example.com"}
		updated, err := repo.Update(u)
		if err != nil {
			t.Fatalf("esperava sucesso no update, erro: %v", err)
		}
		if updated.Name != "Alice Silva" {
			t.Errorf("esperava Alice Silva, veio %s", updated.Name)
		}
	})

	t.Run("atualizar entidade inexistente retorna erro", func(t *testing.T) {
		u := &core.User{ID: "u-nao-existe", Name: "Ninguem"}
		_, err := repo.Update(u)
		if err != adapters.ErrEntityNotFound {
			t.Fatalf("esperava ErrEntityNotFound, veio %v", err)
		}
	})

	t.Run("retorna todas as entidades", func(t *testing.T) {
		u2 := &core.User{ID: "u-2", Name: "Bob", Email: "bob@example.com"}
		_, _ = repo.Create(u2)

		all, err := repo.FindAll()
		if err != nil {
			t.Fatalf("erro ao buscar todas: %v", err)
		}
		if len(all) != 2 {
			t.Errorf("esperava 2 itens, obteve %d", len(all))
		}
	})

	t.Run("remove entidade com sucesso", func(t *testing.T) {
		err := repo.Delete("u-1")
		if err != nil {
			t.Fatalf("erro ao deletar: %v", err)
		}

		_, err = repo.FindByID("u-1")
		if err != adapters.ErrEntityNotFound {
			t.Fatalf("esperava ErrEntityNotFound apos delete, veio %v", err)
		}
	})

	t.Run("deletar ID inexistente retorna erro", func(t *testing.T) {
		err := repo.Delete("inexistente")
		if err != adapters.ErrEntityNotFound {
			t.Fatalf("esperava ErrEntityNotFound, veio %v", err)
		}
	})

	t.Run("seguranca em operacoes concorrentes", func(t *testing.T) {
		cRepo := adapters.NewMemoryRepository[*core.User]()
		var wg sync.WaitGroup

		for i := range 50 {
			wg.Add(1)
			go func(idx int) {
				defer wg.Done()
				user := &core.User{
					ID:    fmt.Sprintf("conc-%d", idx),
					Name:  fmt.Sprintf("User %d", idx),
					Email: fmt.Sprintf("u%d@test.com", idx),
				}
				_, _ = cRepo.Create(user)
				_, _ = cRepo.FindByID(user.ID)
			}(i)
		}

		wg.Wait()

		all, err := cRepo.FindAll()
		if err != nil {
			t.Fatalf("erro ao ler apos concorrencia: %v", err)
		}
		if len(all) != 50 {
			t.Errorf("esperava 50 usuarios salvos concorrentemente, veio %d", len(all))
		}
	})
}
