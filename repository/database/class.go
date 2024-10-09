package database

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/mutezebra/ClassroomRandomRollCallSystem/biz/model/api/class"
	"github.com/mutezebra/ClassroomRandomRollCallSystem/biz/model/api/user"
	"github.com/mutezebra/ClassroomRandomRollCallSystem/pkg/excel"
	"github.com/pkg/errors"
)

type ClassRepository struct {
	db *sql.DB
}

func NewClassRepository() *ClassRepository {
	return &ClassRepository{
		db: _db,
	}
}

func (repo *ClassRepository) CreateClass(ctx context.Context, req *class.Class, uid int64) error {
	tx, err := repo.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return errors.Wrap(err, "failed when begin tx")
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	query := "INSERT INTO class(id,name,user_count,invitation_code) VALUES (?,?,?,?)"
	if _, err = repo.db.Exec(query, req.GetID(), req.GetName(), req.GetUserCount(), req.GetInvitationCode()); err != nil {
		return errors.Wrap(err, "failed when create a class")
	}
	query = "INSERT INTO class_owner(uid, class_id) VALUES (?,?)"
	if _, err = repo.db.Exec(query, uid, req.GetID()); err != nil {
		return errors.Wrap(err, "failed when insert item to `class_owner`")
	}
	if err = tx.Commit(); err != nil {
		return errors.Wrap(err, "failed when tx commit")
	}

	return nil
}

func (repo *ClassRepository) ClassExistByID(ctx context.Context, classID int64) (exist bool, err error) {
	query := "SELECT EXISTS(SELECT 1 from class WHERE id=?)"
	if err = repo.db.QueryRow(query, classID).Scan(&exist); err != nil && !errors.Is(err, sql.ErrNoRows) {
		return false, errors.Wrap(err, fmt.Sprintf("failed when query class %d whether exist", classID))
	}
	return exist, nil
}

func (repo *ClassRepository) IsClassOwner(ctx context.Context, uid, classID int64) (is bool, err error) {
	query := "SELECT EXISTS(SELECT 1 from class_owner WHERE uid=? AND class_id=?)"
	if err = repo.db.QueryRow(query, uid, classID).Scan(&is); err != nil && !errors.Is(err, sql.ErrNoRows) {
		return false, errors.Wrap(err, fmt.Sprintf("failed when query whether %d is %d`s owner", uid, classID))
	}
	return is, nil
}

func (repo *ClassRepository) WhetherUserInClass(ctx context.Context, uid, classID int64) (in bool, err error) {
	query := "SELECT EXISTS(SELECT 1 from user_with_class WHERE uid=? AND class_id=?) OR EXISTS(SELECT 1 from class_owner WHERE uid=? AND class_id=?) "
	if err = repo.db.QueryRow(query, uid, classID, uid, classID).Scan(&in); err != nil && !errors.Is(err, sql.ErrNoRows) {
		return false, errors.Wrap(err, fmt.Sprintf("failed when query whether %d have in %d", uid, classID))
	}
	return in, nil
}

func (repo *ClassRepository) FindClassInvitationCode(ctx context.Context, classID int64) (code string, err error) {
	query := "SELECT invitation_code from class WHERE id=?"
	if err = repo.db.QueryRow(query, classID).Scan(&code); err != nil {
		return "", errors.Wrap(err, fmt.Sprintf("failed when get %d`s invivation_code", classID))
	}
	return code, nil
}

func (repo *ClassRepository) JoinClass(ctx context.Context, uid, classID int64) error {
	tx, err := repo.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return errors.Wrap(err, "failed when begin tx")
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	query := "INSERT INTO user_with_class(uid, class_id,weight) VALUES (?,?,100)"
	if _, err = tx.Exec(query, uid, classID); err != nil {
		return errors.Wrap(err, fmt.Sprintf("failed when %d try join %d", uid, classID))
	}
	query = "UPDATE class SET user_count=user_count+1 WHERE id=?"
	if _, err = tx.Exec(query, classID); err != nil {
		return errors.Wrap(err, "failed when ")
	}
	if err = tx.Commit(); err != nil {
		return errors.Wrap(err, "failed when tx commit")
	}

	return nil
}

func (repo *ClassRepository) ClassList(ctx context.Context, uid int64) (classes []*class.Class, err error) {
	classes = make([]*class.Class, 0)
	query := "SELECT id,name,user_count FROM class WHERE id IN (SELECT class_id from user_with_class WHERE uid=? UNION SELECT class_id FROM class_owner WHERE uid=?)"
	var rows *sql.Rows
	if rows, err = repo.db.Query(query, uid, uid); err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, errors.Wrap(err, fmt.Sprintf("failed when try find %d`s classes", uid))
	}
	defer rows.Close()

	for rows.Next() {
		var c class.Class
		if err = rows.Scan(&c.ID, &c.Name, &c.UserCount); err != nil {
			return nil, errors.Wrap(err, fmt.Sprintf("failed when scan class to classes"))
		}
		classes = append(classes, &c)
	}
	return classes, nil
}

func (repo *ClassRepository) ClassStudentList(ctx context.Context, classID int64) (users []*user.BaseUser, err error) {
	users = make([]*user.BaseUser, 0)
	query := "SELECT student_number,name,avatar,id FROM user WHERE id in (SELECT uid FROM user_with_class WHERE class_id=?)"
	var rows *sql.Rows
	if rows, err = repo.db.Query(query, classID); err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, errors.Wrap(err, fmt.Sprintf("failed when query %d`s students", classID))
	}
	defer rows.Close()

	for rows.Next() {
		var u user.BaseUser
		if err = rows.Scan(&u.StudentNumber, &u.Name, &u.Avatar, &u.UID); err != nil {
			return nil, errors.Wrap(err, "failed when scan u to user")
		}
		users = append(users, &u)
	}
	return users, nil
}

func (repo *ClassRepository) GetTeacherInfo(ctx context.Context, classID int64) (us *user.BaseUser, err error) {
	us = &user.BaseUser{}
	query := "SELECT name,avatar,id FROM user WHERE id=(SELECT uid FROM class_owner WHERE class_id=?)"
	if err = repo.db.QueryRow(query, classID).Scan(&us.Name, &us.Avatar, &us.Avatar); err != nil {
		return us, errors.Wrap(err, "failed when try to find class`s teacher")
	}
	return us, nil
}

func (repo *ClassRepository) ImportUserAndCreateClass(ctx context.Context, classID, uid int64, className, iCode string, pwd string, users []*excel.ImportUser) error {
	tx, err := repo.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return errors.Wrap(err, "failed when begin tx")
	}
	defer func() {
		if err != nil {
			newerr := tx.Rollback()
			if newerr != nil {
				fmt.Println("rollback failed")
			}
		}
	}()

	query := "INSERT INTO class(id,name,user_count,invitation_code) VALUES (?,?,?,?)"
	if _, err = repo.db.Exec(query, classID, className, len(users), iCode); err != nil {
		return errors.Wrap(err, "failed when create a class")
	}
	query = "INSERT INTO class_owner(uid, class_id) VALUES (?,?)"
	if _, err = repo.db.Exec(query, uid, classID); err != nil {
		return errors.Wrap(err, "failed when insert item to `class_owner`")
	}

	// create and join
	for _, u := range users {
		if _, err = repo.db.Exec("INSERT INTO user(id,name,student_number,avatar,phone_number,password_digest) values (?,?,?,?,?,?)",
			u.UID, u.Name, u.StudentNumber, "avatar",
			u.PhoneNumber, pwd,
		); err != nil {
			return errors.Wrap(err, "insert item to user failed")
		}

		query = "INSERT INTO user_with_class(uid, class_id,weight) VALUES (?,?,100)"
		if _, err = tx.Exec(query, u.UID, classID); err != nil {
			return errors.Wrap(err, fmt.Sprintf("failed when %d try join %d", u.UID, classID))
		}
	}

	if err = tx.Commit(); err != nil {
		return errors.Wrap(err, "failed when tx commit")
	}
	return nil
}
