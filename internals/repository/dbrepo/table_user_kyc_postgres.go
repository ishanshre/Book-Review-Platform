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
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	query := `
		select u.id, u.username, u.email, u.password, u.access_level,
			u.created_at, u.updated_at, u.last_login, k.id, k.user_id, k.first_name,
			k.last_name, k.gender, k.address, k.phone, k.profile_pic, k.dob, k.document_number,
			k.document_front, k.document_back, k.is_validated, k.updated_at
		from users as u
		join kycs as k ON u.id=k.user_id
		where u.id = $1
	`
	user := &models.User{}
	kyc := &models.Kyc{}
	row := m.DB.QueryRowContext(ctx, query, id)
	if err := row.Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Password,
		&user.AccessLevel,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.LastLogin,
		&kyc.ID,
		&kyc.UserID,
		&kyc.FirstName,
		&kyc.LastName,
		&kyc.Gender,
		&kyc.Address,
		&kyc.Phone,
		&kyc.ProfilePic,
		&kyc.DateOfBirth,
		&kyc.DocumentNumber,
		&kyc.DocumentFront,
		&kyc.DocumentBack,
		&kyc.IsValidated,
		&kyc.UpdatedAt,
	); err != nil {
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
	stmt := `UPDATE kycs SET profile_pic=$2 WHERE user_id=$1`
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
	stmt := `UPDATE kycs SET document_front=$2, document_back=$3 WHERE user_id=$1`
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
		WHERE user_id=$1`
	_, err := m.DB.ExecContext(
		ctx,
		stmt,
		update.UserID,
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
	return nil
}

func (m *postgresDBRepo) PublicKycUpdate(update *models.Kyc) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	stmt := `
		UPDATE kycs 
		SET first_name = $2, last_name = $3, gender = $4, phone = $5, address = $6, dob = $7, document_type = $8, document_number = $9, document_front = $10, document_back = $11, is_validated = $12, updated_at = $13
		WHERE id=$1`
	_, err := m.DB.ExecContext(
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
		update.DocumentFront,
		update.DocumentBack,
		update.IsValidated,
		update.UpdatedAt,
	)
	if err != nil {
		return err
	}
	return nil
}
