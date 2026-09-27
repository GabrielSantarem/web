package adapters_test

import (
	"database/sql"
	"errors"
	"testing"

	"codeberg.org/MrTomate/web/internal/adapters"
	"codeberg.org/MrTomate/web/internal/core"
	_ "modernc.org/sqlite"
)

func setupTestSQLite(t *testing.T) *adapters.SQLiteUserRepository {
	t.Helper()

	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("erro ao abrir sqlite em memoria: %v", err)
	}

	repo, err := adapters.NewSQLiteUserRepository(db)
	if err != nil {
		t.Fatalf("erro ao inicializar sqlite repository: %v", err)
	}

	t.Cleanup(func() {
		_ = db.Close()
	})

	return repo
}

func TestSQLiteUserRepository(t *testing.T) {
	t.Run("Create e FindByID salvam e recuperam usuario", func(t *testing.T) {
		repo := setupTestSQLite(t)

		user := core.NewUser("sqlite@test.com", "SQLite User")
		created, err := repo.Create(user)
		if err != nil {
			t.Fatalf("erro ao criar: %v", err)
		}
		if created.ID != "1" {
			t.Errorf("esperava id 1, veio %s", created.ID)
		}

		found, err := repo.FindByID("1")
		if err != nil {
			t.Fatalf("erro ao buscar: %v", err)
		}
		if found.Email != "sqlite@test.com" || found.Name != "SQLite User" {
			t.Errorf("dados incorretos: %+v", found)
		}
		if found.IsActive {
			t.Error("usuario deveria iniciar inativo")
		}
	})

	t.Run("Create rejeita email duplicado", func(t *testing.T) {
		repo := setupTestSQLite(t)

		u1 := core.NewUser("duplicado@test.com", "User 1")
		_, err := repo.Create(u1)
		if err != nil {
			t.Fatalf("primeira criacao falhou: %v", err)
		}

		u2 := core.NewUser("duplicado@test.com", "User 2")
		_, err = repo.Create(u2)
		if !errors.Is(err, core.ErrEmailAlreadyUsed) {
			t.Fatalf("esperava ErrEmailAlreadyUsed, veio %v", err)
		}
	})

	t.Run("FindByID inexistente retorna ErrUserNotFound", func(t *testing.T) {
		repo := setupTestSQLite(t)

		_, err := repo.FindByID("999")
		if !errors.Is(err, core.ErrUserNotFound) {
			t.Fatalf("esperava ErrUserNotFound, veio %v", err)
		}
	})

	t.Run("FindByEmail encontra usuario ou retorna ErrUserNotFound", func(t *testing.T) {
		repo := setupTestSQLite(t)

		u := core.NewUser("busca@email.com", "Busca Email")
		_, _ = repo.Create(u)

		found, err := repo.FindByEmail("busca@email.com")
		if err != nil {
			t.Fatalf("esperava encontrar usuario, veio erro: %v", err)
		}
		if found.Name != "Busca Email" {
			t.Errorf("esperava nome Busca Email, veio %s", found.Name)
		}

		_, err = repo.FindByEmail("naoexiste@email.com")
		if !errors.Is(err, core.ErrUserNotFound) {
			t.Fatalf("esperava ErrUserNotFound, veio %v", err)
		}
	})

	t.Run("Update altera dados e status de ativo", func(t *testing.T) {
		repo := setupTestSQLite(t)

		u, _ := repo.Create(core.NewUser("update@test.com", "Antes"))
		u.Name = "Depois"
		u.IsActive = true

		updated, err := repo.Update(u)
		if err != nil {
			t.Fatalf("erro ao atualizar: %v", err)
		}
		if updated.Name != "Depois" || !updated.IsActive {
			t.Errorf("atualizacao incorreta: %+v", updated)
		}

		// Tentar atualizar usuario inexistente
		inexistente := &core.User{ID: "9999", Name: "Nao Existe", Email: "no@test.com"}
		_, err = repo.Update(inexistente)
		if !errors.Is(err, core.ErrUserNotFound) {
			t.Fatalf("esperava ErrUserNotFound ao atualizar inexistente, veio %v", err)
		}
	})

	t.Run("FindAll lista todos os usuarios em ordem de insercao", func(t *testing.T) {
		repo := setupTestSQLite(t)

		_, _ = repo.Create(core.NewUser("u1@test.com", "User 1"))
		_, _ = repo.Create(core.NewUser("u2@test.com", "User 2"))

		all, err := repo.FindAll()
		if err != nil {
			t.Fatalf("erro ao listar todos: %v", err)
		}
		if len(all) != 2 {
			t.Fatalf("esperava 2 usuarios, vieram %d", len(all))
		}
	})

	t.Run("Delete remove usuario e falha se nao existir", func(t *testing.T) {
		repo := setupTestSQLite(t)

		u, _ := repo.Create(core.NewUser("del@test.com", "Delete Me"))

		if err := repo.Delete(u.ID); err != nil {
			t.Fatalf("erro ao deletar: %v", err)
		}

		_, err := repo.FindByID(u.ID)
		if !errors.Is(err, core.ErrUserNotFound) {
			t.Fatalf("usuario ainda existe apos delete: %v", err)
		}

		err = repo.Delete("9999")
		if !errors.Is(err, core.ErrUserNotFound) {
			t.Fatalf("esperava ErrUserNotFound ao deletar id inexistente, veio %v", err)
		}
	})

	t.Run("Integracao: UserService funciona perfeitamente com SQLite", func(t *testing.T) {
		repo := setupTestSQLite(t)
		service := core.NewUserService(repo)

		// 1. Criar via servico
		u, err := service.CreateUser(core.NewUser("servico@test.com", "Via Servico"))
		if err != nil {
			t.Fatalf("falha ao criar via servico: %v", err)
		}

		// 2. Buscar via servico
		found, err := service.GetUser(u.ID)
		if err != nil {
			t.Fatalf("falha ao buscar via servico: %v", err)
		}
		if found.Name != "Via Servico" {
			t.Errorf("nome incorreto: %s", found.Name)
		}

		// 3. Ativar via servico
		err = service.ActivateUser(u.ID)
		if err != nil {
			t.Fatalf("falha ao ativar via servico: %v", err)
		}

		// 4. Reativar deve acusar ErrUserAlreadyActive
		err = service.ActivateUser(u.ID)
		if !errors.Is(err, core.ErrUserAlreadyActive) {
			t.Fatalf("esperava ErrUserAlreadyActive, veio %v", err)
		}
	})
}
