package services

import (
	"crypto/sha1"
	"errors"
	"fmt"
	"github.com/architectv/networking-course-project/backend/pkg/models"
	"github.com/architectv/networking-course-project/backend/pkg/repositories"
	"time"

	"github.com/dgrijalva/jwt-go"
)

const (
	salt       = "hjqrhjqw124617ajfhajs"
	signingKey = "qrkjk#4#%35FSFJlja#4353KSFjH"
	tokenTTL   = 12 * time.Hour
)

type tokenClaims struct {
	jwt.StandardClaims
	UserId int `json:"id"`
}

type UserService struct {
	repo repositories.User
}

func NewUserService(repo repositories.User) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) Get(id int) *models.ApiResponse {
	r := &models.ApiResponse{}
	user, err := s.repo.GetById(id)
	if err != nil {
		r.Error(StatusInternalServerError, err.Error())
		return r
	}

	r.Set(StatusOK, "OK", Map{"user": user})
	return r
}

func (s *UserService) GetAll() ([]*models.User, error) {
	return s.repo.GetAll()
}

func (s *UserService) Update(id int, profile *models.UpdateUser) *models.ApiResponse {
	r := &models.ApiResponse{}
	if profile.Nickname != nil {
		if err := s.checkByNickname(*profile.Nickname); err == nil {
			r.Error(StatusConflict, "User already exists")
			return r
		}
	}

	err := s.repo.Update(id, profile)
	if err != nil {
		r.Error(StatusInternalServerError, err.Error())
		return r
	}

	r.Set(StatusOK, "OK", Map{})
	return r
}

func (s *UserService) Create(user *models.User) *models.ApiResponse {
	r := &models.ApiResponse{}
	if err := s.checkByNickname(user.Nickname); err == nil {
		r.Error(StatusConflict, "User already exists")
		return r
	}

	// TODO: generatePasswordHash(password)
	// user.Password = generatePasswordHash(user.Password)
	id, err := s.repo.Create(user)
	if err != nil {
		r.Error(StatusInternalServerError, err.Error())
		return r
	}

	r.Set(StatusOK, "OK", Map{userId: id})
	return r
}

func (s *UserService) checkByNickname(nickname string) error {
	_, err := s.repo.GetByNickname(nickname)
	return err
}

func (s *UserService) GenerateToken(nickname, password string) *models.ApiResponse {
	// TODO: generatePasswordHash(password)
	r := &models.ApiResponse{}
	user, err := s.repo.Get(nickname, password)
	if err != nil {
		r.Error(StatusConflict, err.Error())
		return r
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(tokenTTL).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		user.Id,
	})

	encToken, err := token.SignedString([]byte(signingKey))
	if err != nil {
		r.Error(StatusConflict, err.Error())
		return r
	}

	r.Set(StatusOK, "OK", Map{"token": encToken})
	return r
}

func (s *UserService) ParseToken(accessToken string) (int, error) {
	token, err := jwt.ParseWithClaims(accessToken, &tokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}

		return []byte(signingKey), nil
	})
	if err != nil {
		return 0, err
	}

	claims, ok := token.Claims.(*tokenClaims)
	if !ok {
		return 0, errors.New("token claims are not of type *tokenClaims")
	}

	err = s.repo.FindToken(accessToken)
	// logrus.Println(err.Error())
	if err == nil {
		return 0, errors.New("Invalid token")
	}

	return claims.UserId, nil
}

func generatePasswordHash(password string) string {
	hash := sha1.New()
	hash.Write([]byte(password))

	return fmt.Sprintf("%x", hash.Sum([]byte(salt)))
}

func (s *UserService) SignOut(token string) *models.ApiResponse {
	r := &models.ApiResponse{}
	_, err := s.repo.SignOut(token)
	if err != nil {
		r.Error(StatusInternalServerError, err.Error())
		return r
	}

	r.Set(StatusOK, "OK", nil)
	return r
}
