package database

import (
	"database/sql"
	"github.com/mutezebra/ClassroomRandomRollCallSystem/biz/model/api/user"
	"github.com/mutezebra/ClassroomRandomRollCallSystem/pkg/pack"
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

func (repo *UserRepository) FindUserByUID(uids []int64) (users []*user.BaseUser, err error) {
	var pre *sql.Stmt
	defer func() {
		pack.LogError(pre.Close())
	}()
	if pre, err = repo.db.Prepare("SELECT id,name,student_number,avatar FROM user WHERE id=?"); err != nil {
		return nil, errors.Wrap(err, "failed when prepare query")
	}
	users = make([]*user.BaseUser, 0, len(uids))
	for _, uid := range uids {
		var u user.BaseUser
		if err = pre.QueryRow(uid).Scan(&u.UID, &u.Name, &u.StudentNumber, &u.Avatar); err != nil {
			return nil, errors.Wrap(err, "failed when scan user")
		}
		users = append(users, &u)
	}
	return users, nil
}
