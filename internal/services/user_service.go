package services

import (
	"errors"

	"github.com/EthangarciaDev/backend-AquiEstoy/internal/config"
	"github.com/EthangarciaDev/backend-AquiEstoy/internal/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserService struct {
	db *gorm.DB
}

// NewUserService crea una nueva instancia del servicio de usuarios
func NewUserService() *UserService {
	return &UserService{
		db: config.GetDB(),
	}
}

// CreateUser crea un nuevo usuario en la base de datos
func (s *UserService) CreateUser(user *models.User) error {
	// Hash del password antes de guardar
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hashedPassword)

	// Crear el usuario en la base de datos
	result := s.db.Create(user)
	return result.Error
}

// GetUserByEmail busca un usuario por email
func (s *UserService) GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	result := s.db.Where("email = ?", email).First(&user)
	
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, errors.New("usuario no encontrado")
	}
	
	return &user, result.Error
}

// GetUserByID busca un usuario por ID
func (s *UserService) GetUserByID(id uint) (*models.User, error) {
	var user models.User
	result := s.db.First(&user, id)
	
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, errors.New("usuario no encontrado")
	}
	
	return &user, result.Error
}

// ValidatePassword verifica si el password proporcionado coincide con el hash almacenado
func (s *UserService) ValidatePassword(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}

// MigrateDB ejecuta las migraciones de la base de datos
func (s *UserService) MigrateDB() error {
	return s.db.AutoMigrate(&models.User{})
}