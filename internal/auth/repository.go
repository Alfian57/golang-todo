package auth

import (
	"context"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AuthRepository struct {
	db *gorm.DB
}

func NewAuthRepository(db *gorm.DB) AuthRepository {
	return AuthRepository{
		db: db,
	}
}

func (repository AuthRepository) FindUserByUsername(ctx context.Context, username string) (User, error) {
	user, err := gorm.G[User](repository.db).Where("username = ? ", username).First(ctx)
	return user, err
}

func (repository AuthRepository) FindUserByID(ctx context.Context, userID uuid.UUID) (User, error) {
	user, err := gorm.G[User](repository.db).Where("id = ?", userID).First(ctx)
	return user, err
}

func (repository AuthRepository) CreateUser(ctx context.Context, username string, password string) (User, error) {
	user := User{
		Username: username,
		Password: password,
	}

	result := gorm.WithResult()
	err := gorm.G[User](repository.db, result).Create(ctx, &user)

	return user, err
}

func (repository AuthRepository) CreateRefreshToken(ctx context.Context, token string, userID uuid.UUID, expiredAt time.Time) (RefreshToken, error) {
	refreshToken := RefreshToken{
		Token:     token,
		UserID:    userID,
		ExpiredAt: expiredAt,
	}

	result := gorm.WithResult()
	err := gorm.G[RefreshToken](repository.db, result).Create(ctx, &refreshToken)

	return refreshToken, err
}

func (repository AuthRepository) FindRefreshTokenByToken(ctx context.Context, token string) (RefreshToken, error) {
	refreshToken, err := gorm.G[RefreshToken](repository.db).Where("token = ?", token).First(ctx)
	return refreshToken, err
}

func (repository AuthRepository) DeleteRefreshTokenByToken(ctx context.Context, token string) (int, error) {
	rowsAffected, err := gorm.G[RefreshToken](repository.db).Where("token = ?", token).Delete(ctx)
	return rowsAffected, err
}
