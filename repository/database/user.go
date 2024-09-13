package database

import (
	"database/sql"
	"github.com/mutezebra/ClassroomRandomRollCallSystem/biz/model/api/user"
	"github.com/pkg/errors"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository() *UserRepository {
	return &UserRepository{
		db: _db,
	}
}

func (repo *UserRepository) CreateUser(user *user.User) error {
	_, err := repo.db.Exec("INSERT INTO user(id,name,student_number,avatar,phone_number,password_digest) values (?,?,?,?,?,?)",
		user.ID, user.Name, user.StudentNumber, user.Avatar,
		user.PhoneNumber, user.PasswordDigest,
	)
	if err != nil {
		return errors.Wrap(err, "insert item to user failed")
	}
	return nil
}

func (repo *UserRepository) PhoneNumberExist(phoneNumber string) (bool, error) {
	var exist bool
	query := "SELECT EXISTS(SELECT 1 from user WHERE phone_number=?)"
	err := repo.db.QueryRow(query, phoneNumber).Scan(&exist)
	if err != nil {
		return false, errors.Wrap(err, "query user whether exist base phone_number")
	}
	return exist, err
}

func (repo *UserRepository) FindUserPassword(phoneNumber string) (string, error) {
	var pwd string
	query := "SELECT password_digest from user WHERE phone_number=? LIMIT 1"
	err := repo.db.QueryRow(query, phoneNumber).Scan(&pwd)
	if err != nil {
		return "", errors.Wrap(err, "failed when query pwd base phone_number")
	}
	return pwd, nil
}

func (repo *UserRepository) FindUIDByPhoneNumber(phoneNumber string) (int64, error) {
	var uid int64
	query := "SELECT id from user WHERE phone_number=? LIMIT 1"
	err := repo.db.QueryRow(query, phoneNumber).Scan(&uid)
	if err != nil {
		return 0, errors.Wrap(err, "failed when find user id by phone_number")
	}
	return uid, nil
}
