package dbrepo

import (
	"context"
	"errors"
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

// UpdateProfilePic updates user profile pic
func (m *postgresDBRepo) UpdateProfilePic(path string, id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	stmt := `UPDATE kycs SET profile_pic=$2 WHERE id=$1`
	_, err := m.DB.ExecContext(ctx, stmt, id, path)
	if err != nil {
		return err
	}
	return nil
}

// UpdateDocument updates user profile pic
func (m *postgresDBRepo) UpdateDocument(front_path, back_path string, id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	stmt := `UPDATE kycs SET document_front=$2, document_back=$3 WHERE id=$1`
	_, err := m.DB.ExecContext(ctx, stmt, id, front_path, back_path)
	if err != nil {
		return err
	}
	return nil
}

func (m *postgresDBRepo) AdminKycUpdate(update *models.Kyc) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	stmt := `
		UPDATE kycs 
		SET first_name = $2, last_name = $3, gender = $4, phone = $5, address = $6, dob = $7, is_validated = $8, document_type = $9, document_number = $10, updated_at = $11
		WHERE id=$1`
	res, err := m.DB.ExecContext(
		ctx,
		stmt,
		update.ID,
		update.FirstName,
		update.LastName,
		update.Gender,
		update.Phone,
		update.Address,
		update.DateOfBirth,
		update.IsValidated,
		update.DocumentType,
		update.DocumentNumber,
		update.UpdatedAt,
	)
	if err != nil {
		return err
	}
	affected, _ := res.RowsAffected()
	if affected == 0 {
		return errors.New("not updateed")
	}
	return nil
}

func (m *postgresDBRepo) PublicKycUpdate(update *models.Kyc) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	stmt := `
		UPDATE kycs 
		SET first_name = $2, last_name = $3, gender = $4, phone = $5, address = $6, dob = $7, document_type = $8, document_number = $9, updated_at = $10
		WHERE id=$1`
	res, err := m.DB.ExecContext(
		ctx,
		stmt,
		update.ID,
		update.FirstName,
		update.LastName,
		update.Gender,
		update.Phone,
		update.Address,
		update.DateOfBirth,
		update.DocumentType,
		update.DocumentNumber,
		update.UpdatedAt,
	)
	if err != nil {
		return err
	}
	affected, _ := res.RowsAffected()
	if affected == 0 {
		return errors.New("not updateed")
	}
	return nil
}
