package core

// UserRepository define as operações de persistência de que o domínio precisa.
type UserRepository interface {
	// CreateUser cria um novo usuário no repositório.
	CreateUser(u *User) (*User, error)

	// UpdateUser atualiza um usuário no repositório.
	UpdateUser(u *User) (*User, error)

	// FindById retorna um usuário pelo ID.
	FindById(id string) (*User, error)
	// FindByEmail retorna um usuário pelo email.
	FindByEmail(email string) (*User, error)
}

// User representa um usuário no sistema.
// Nossa classe de domínio.
type User struct {
	ID       string
	Email    string
	Name     string
	IsActive bool
}

// NewUser cria um novo usuário com os dados fornecidos.
// Retorna um ponteiro para um novo User com os campos preenchidos.
func NewUser(email, name string) *User {
	return &User{
		Email: email,
		Name:  name,
	}
}

// UserService define as que sao expostas para o mundo externo.
type UserService interface {
	// CreateUser cria um novo usuário no sistema.
	CreateUser(u *User) (*User, error)
	// GetUser retorna um usuário pelo ID.
	GetUser(id string) (*User, error)
	// ActivateUser ativa um usuário pelo ID.
	ActivateUser(id string) error
}

// userService implementa as operações de serviço para usuários.
// Ele atua como intermediário entre o repositório e o mundo externo.
type userService struct {
	repo UserRepository
}

func NewUserService(r UserRepository) UserService {
	return &userService{repo: r}
}

func (s *userService) CreateUser(u *User) (*User, error) {
	if u.Email == "" {
		return nil, ErrInvalidInput
	}
	if u.Name == "" {
		return nil, ErrInvalidInput
	}

	return s.repo.CreateUser(u)
}

func (s *userService) GetUser(id string) (*User, error) {
	return s.repo.FindById(id)
}

func (s *userService) ActivateUser(id string) error {
	user, err := s.GetUser(id)
	if err != nil {
		return err
	}

	if user.IsActive {
		return ErrUserAlreadyActive
	}

	user.IsActive = true

	_, err = s.repo.UpdateUser(user)

	if err != nil {
		return err
	}

	return nil
}
