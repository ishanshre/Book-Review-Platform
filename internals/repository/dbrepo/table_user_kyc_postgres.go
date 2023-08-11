package dbrepo

import (
	"context"
	"time"

	"github.com/ishanshre/Book-Review-Platform/internals/models"
)

func (m *postgresDBRepo) GetKycByUserID(user_id int) (*models.Kyc, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	query := `
		SELECT * FROM kycs WHERE user_id = $1
	`
	kyc := &models.Kyc{}
	if err := m.DB.QueryRowContext(ctx, query, user_id).Scan(
		&kyc.ID,
		&kyc.UserID,
		&kyc.FirstName,
		&kyc.LastName,
		&kyc.Gender,
		&kyc.Address,
		&kyc.Phone,
		&kyc.ProfilePic,
		&kyc.DateOfBirth,
		&kyc.DocumentType,
		&kyc.DocumentNumber,
		&kyc.DocumentFront,
		&kyc.DocumentBack,
		&kyc.IsValidated,
		&kyc.UpdatedAt,
	); err != nil {
		return nil, err
	}
	return kyc, nil
}

func (m *postgresDBRepo) GetUserWithKyc(id int) (*models.UserKycData, error) {
	_, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	user, err := m.GetUserByID(id)
	if err != nil {
		return nil, err
	}
	kyc, err := m.GetKycByUserID(user.ID)
	if err != nil {
		return nil, err
	}
	userWithKyc := &models.UserKycData{}
	userWithKyc.User = user
	userWithKyc.Kyc = kyc
	return userWithKyc, nil
}
