package adapters

import (
	"database/sql"
	"errors"
	"strconv"
	"strings"

	"codeberg.org/MrTomate/web/internal/core"
	_ "modernc.org/sqlite"
)

// SQLiteUserRepository implementa core.UserRepository utilizando SQLite como motor de persistencia.
type SQLiteUserRepository struct {
	db *sql.DB
}

// NewSQLiteUserRepository inicializa o adaptador e cria a tabela caso nao exista.
func NewSQLiteUserRepository(db *sql.DB) (*SQLiteUserRepository, error) {
	schema := `
	CREATE TABLE IF NOT EXISTS users (
		id TEXT PRIMARY KEY,
		name TEXT NOT NULL,
		email TEXT UNIQUE NOT NULL,
		is_active INTEGER NOT NULL DEFAULT 0
	);`

	if _, err := db.Exec(schema); err != nil {
		return nil, err
	}

	return &SQLiteUserRepository{db: db}, nil
}

// Create insere um novo usuario no banco SQLite.
func (r *SQLiteUserRepository) Create(u *core.User) (*core.User, error) {
	if u.ID == "" {
		var maxID sql.NullInt64
		_ = r.db.QueryRow("SELECT MAX(CAST(id AS INTEGER)) FROM users").Scan(&maxID)
		if maxID.Valid {
			u.ID = strconv.FormatInt(maxID.Int64+1, 10)
		} else {
			u.ID = "1"
		}
	}

	isActiveInt := 0
	if u.IsActive {
		isActiveInt = 1
	}

	query := `INSERT INTO users (id, name, email, is_active) VALUES (?, ?, ?, ?)`
	_, err := r.db.Exec(query, u.ID, u.Name, u.Email, isActiveInt)
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE") || strings.Contains(err.Error(), "constraint failed") {
			return nil, core.ErrEmailAlreadyUsed
		}
		return nil, err
	}

	return u, nil
}

// FindByID busca um usuario pela chave primaria.
func (r *SQLiteUserRepository) FindByID(id string) (*core.User, error) {
	query := `SELECT id, name, email, is_active FROM users WHERE id = ?`
	row := r.db.QueryRow(query, id)

	var u core.User
	var isActiveInt int
	err := row.Scan(&u.ID, &u.Name, &u.Email, &isActiveInt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, core.ErrUserNotFound
		}
		return nil, err
	}

	u.IsActive = isActiveInt == 1
	return &u, nil
}

// FindByEmail busca um usuario pelo endereco de email.
func (r *SQLiteUserRepository) FindByEmail(email string) (*core.User, error) {
	query := `SELECT id, name, email, is_active FROM users WHERE email = ?`
	row := r.db.QueryRow(query, email)

	var u core.User
	var isActiveInt int
	err := row.Scan(&u.ID, &u.Name, &u.Email, &isActiveInt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, core.ErrUserNotFound
		}
		return nil, err
	}

	u.IsActive = isActiveInt == 1
	return &u, nil
}

// Update atualiza os dados de um usuario existente.
func (r *SQLiteUserRepository) Update(u *core.User) (*core.User, error) {
	isActiveInt := 0
	if u.IsActive {
		isActiveInt = 1
	}

	query := `UPDATE users SET name = ?, email = ?, is_active = ? WHERE id = ?`
	res, err := r.db.Exec(query, u.Name, u.Email, isActiveInt, u.ID)
	if err != nil {
		return nil, err
	}

	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		return nil, core.ErrUserNotFound
	}

	return u, nil
}

// FindAll retorna todos os usuarios armazenados.
func (r *SQLiteUserRepository) FindAll() ([]*core.User, error) {
	query := `SELECT id, name, email, is_active FROM users ORDER BY CAST(id AS INTEGER) ASC`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*core.User
	for rows.Next() {
		var u core.User
		var isActiveInt int
		if err := rows.Scan(&u.ID, &u.Name, &u.Email, &isActiveInt); err != nil {
			return nil, err
		}
		u.IsActive = isActiveInt == 1
		users = append(users, &u)
	}

	return users, nil
}

// Delete remove um usuario pelo ID.
func (r *SQLiteUserRepository) Delete(id string) error {
	query := `DELETE FROM users WHERE id = ?`
	res, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		return core.ErrUserNotFound
	}

	return nil
}

// Metodos de compatibilidade com UserRepository
func (r *SQLiteUserRepository) CreateUser(u *core.User) (*core.User, error) { return r.Create(u) }
func (r *SQLiteUserRepository) UpdateUser(u *core.User) (*core.User, error) { return r.Update(u) }
func (r *SQLiteUserRepository) FindById(id string) (*core.User, error)       { return r.FindByID(id) }
