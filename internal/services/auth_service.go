package services

import (
	"errors"
	"os"
	"time"

	"github.com/EthangarciaDev/backend-AquiEstoy/internal/models"
	"github.com/golang-jwt/jwt/v5"
)

// Claims estructura para el payload del JWT
type Claims struct {
	UserID uint   `json:"user_id"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

type AuthService struct {
	userService *UserService
	jwtSecret   []byte
}

// NewAuthService crea una nueva instancia del servicio de autenticación
func NewAuthService(userService *UserService) *AuthService {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "default-secret-change-in-production" // Valor por defecto, cambiar en producción
	}

	return &AuthService{
		userService: userService,
		jwtSecret:   []byte(secret),
	}
}

// Register registra un nuevo usuario
func (s *AuthService) Register(req *models.UserRegisterRequest) (*models.UserResponse, error) {
	// Verificar si el usuario ya existe
	existingUser, _ := s.userService.GetUserByEmail(req.Email)
	if existingUser != nil {
		return nil, errors.New("el email ya está registrado")
	}

	// Crear nuevo usuario
	user := &models.User{
		Email:    req.Email,
		Password: req.Password,
		Name:     req.Name,
	}

	err := s.userService.CreateUser(user)
	if err != nil {
		return nil, err
	}

	// Generar token JWT
	token, err := s.GenerateToken(user.ID, user.Email)
	if err != nil {
		return nil, err
	}

	return &models.UserResponse{
		ID:    user.ID,
		Email: user.Email,
		Name:  user.Name,
		Token: token,
	}, nil
}

// Login autentica un usuario existente
func (s *AuthService) Login(req *models.UserLoginRequest) (*models.UserResponse, error) {
	// Buscar usuario por email
	user, err := s.userService.GetUserByEmail(req.Email)
	if err != nil {
		return nil, errors.New("credenciales inválidas")
	}

	// Verificar password
	if !s.userService.ValidatePassword(user.Password, req.Password) {
		return nil, errors.New("credenciales inválidas")
	}

	// Generar token JWT
	token, err := s.GenerateToken(user.ID, user.Email)
	if err != nil {
		return nil, err
	}

	return &models.UserResponse{
		ID:    user.ID,
		Email: user.Email,
		Name:  user.Name,
		Token: token,
	}, nil
}

// GenerateToken genera un token JWT para el usuario
func (s *AuthService) GenerateToken(userID uint, email string) (string, error) {
	claims := &Claims{
		UserID: userID,
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.jwtSecret)
}

// ValidateToken valida un token JWT
func (s *AuthService) ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return s.jwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("token inválido")
}