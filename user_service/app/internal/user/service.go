package user

import (
	"context"
	"errors"
	"fmt"

	"github.com/Cadet-Blue/backend-go/user_service/internal/apperror"
	"github.com/Cadet-Blue/backend-go/user_service/pkg/logging"
)

var _ Service = &service{}

type service struct {
	storage Storage
	logger  logging.Logger
}

func NewService(userStorage Storage, logger logging.Logger) (Service, error) {
	return &service{
		storage: userStorage,
		logger:  logger,
	}, nil
}

type Service interface {
	Create(ctx context.Context, dto CreateUserDTO) (string, error)
	CheckEmailExist(ctx context.Context, email string) bool
}

func (s service) Create(ctx context.Context, dto CreateUserDTO) (userUUID string, err error) {
	s.logger.Debug("check password and repeat password")
	if dto.Password != dto.RepeatPassword {
		return userUUID, apperror.BadRequestError("password does not match repeat password")
	}

	user := NewUser(dto)

	s.logger.Debug("generate password hash")
	err = user.GeneratePasswordHash()
	if err != nil {
		s.logger.Errorf("failed to create user due to error %v", err)
		return
	}

	userUUID, err = s.storage.Create(ctx, user)

	if err != nil {
		if errors.Is(err, apperror.ErrNotFound) {
			return userUUID, err
		}
		return userUUID, fmt.Errorf("failed to create user. error: %w", err)
	}

	return userUUID, nil
}

func (s service) CheckEmailExist(ctx context.Context, email string) bool {
	_, err := s.storage.FindByEmail(ctx, email)

	return err != nil
}
