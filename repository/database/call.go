package database

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/mutezebra/ClassroomRandomRollCallSystem/biz/model/api/call"
	"github.com/mutezebra/ClassroomRandomRollCallSystem/pkg/log"
	"github.com/pkg/errors"
	"strings"
)

type CallRepository struct {
	db *sql.DB
}

func NewCallRepository() *CallRepository {
	return &CallRepository{
		db: _db,
	}
}

func convert(uids []int64, extra ...int64) []interface{} {
	result := make([]interface{}, 0, len(uids)+len(extra))
	for _, e := range extra {
		result = append(result, e)
	}
	for _, uid := range uids {
		result = append(result, uid)
	}
	return result
}

func (repo *CallRepository) CreateCallEvent(ctx context.Context, event *call.CallEvent) error {
	query := `
       INSERT INTO call_event (id, call_event_name, class_id, class_name, caller_id, caller_name, start_time, end_time) 
       VALUES (?, ?, ?, (SELECT name FROM class WHERE id = ?), ?, (SELECT name FROM user WHERE id = ?), ?, ?)`

	if _, err := repo.db.ExecContext(ctx, query,
		event.GetID(), event.GetCallEventName(), event.GetClassID(),
		event.GetClassID(), event.GetCallerID(), event.GetCallerID(),
		event.GetStartTime(), event.GetEndTime()); err != nil {
		return errors.Wrap(err, "failed to create call event")
	}

	return nil
}

func (repo *CallRepository) EventExist(ctx context.Context, eventID int64) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM call_event WHERE id = ?)`
	var exist bool
	if err := repo.db.QueryRowContext(ctx, query, eventID).Scan(&exist); err != nil {
		return false, errors.Wrap(err, "failed to get event exist")
	}
	return exist, nil
}

func (repo *CallRepository) GetCallEventByID(ctx context.Context, eventID int64) (*call.CallEvent, error) {
	query := `
			SELECT id, call_event_name, class_id, class_name, caller_id, caller_name, start_time,end_time 
			FROM call_event WHERE id = ?`

	event := &call.CallEvent{}
	if err := repo.db.QueryRowContext(ctx, query, eventID).Scan(
		&event.ID, &event.CallEventName, &event.ClassID, &event.ClassName,
		&event.CallerID, &event.CallerName, &event.StartTime, &event.EndTime); err != nil {
		return nil, errors.Wrap(err, "failed to get call event by id")
	}
	return event, nil
}

func (repo *CallRepository) UpdateEventDone(ctx context.Context, eventID, classID int64, done, undo []int64) error {
	var tx *sql.Tx
	var err error
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	if tx, err = repo.db.BeginTx(ctx, nil); err != nil {
		return errors.Wrap(err, "failed to begin transaction")
	}

	var pre *sql.Stmt
	defer func() {
		_ = pre.Close()
	}()
	if pre, err = tx.Prepare("INSERT INTO call_event_with_user (call_event_id, uid, done) VALUES (?, ?, ?)"); err != nil {
		return errors.Wrap(err, "failed to prepare insert statement")
	}
	for _, uid := range done {
		if _, err = pre.ExecContext(ctx, eventID, uid, true); err != nil {
			return errors.Wrap(err, "failed to insert done uid")
		}
	}
	for _, uid := range undo {
		if _, err = pre.ExecContext(ctx, eventID, uid, false); err != nil {
			return errors.Wrap(err, "failed to insert undo uid")
		}
	}

	chars := strings.Repeat("?,", len(undo)-1)
	query := fmt.Sprintf("UPDATE user_with_class SET weight = weight + 1 WHERE class_id=? AND uid IN (%s?)", chars)
	log.LogrusObj.Info(query)

	// 将切片的元素作为参数传递给 ExecContext 方法
	if _, err = tx.ExecContext(ctx, query, convert(undo, classID)...); err != nil {
		return errors.Wrap(err, "failed to update weight")
	}

	if err = tx.Commit(); err != nil {
		return errors.Wrap(err, "failed to commit transaction")
	}
	return nil
}

func (repo *CallRepository) UpdateUndoneUIDS(ctx context.Context, eventID int64, uids []int64) error {
	var pre *sql.Stmt
	defer func() {
		_ = pre.Close()
	}()
	var err error
	if pre, err = repo.db.Prepare("INSERT INTO call_event_with_user (call_event_id, uid, done) VALUES (?, ?, 0)"); err != nil {
		return errors.Wrap(err, "failed to prepare insert statement")
	}

	for _, uid := range uids {
		if _, err = pre.ExecContext(ctx, eventID, uid); err != nil {
			return errors.Wrap(err, "failed to insert undo uid")
		}
	}
	return nil
}

func (repo *CallRepository) UpdateDoneUIDS(ctx context.Context, eventID int64, uids []int64) error {
	var pre *sql.Stmt
	defer func() {
		_ = pre.Close()
	}()
	var err error
	if pre, err = repo.db.Prepare("INSERT INTO call_event_with_user (call_event_id, uid, done) VALUES (?, ?, 1)"); err != nil {
		return errors.Wrap(err, "failed to prepare insert statement")
	}

	for _, uid := range uids {
		if _, err = pre.ExecContext(ctx, eventID, uid); err != nil {
			return errors.Wrap(err, "failed to insert undo uid")
		}
	}
	return nil
}

func (repo *CallRepository) GetUidAndWeight(ctx context.Context, classID int64) (map[int64]int, error) {
	query := `SELECT uid, weight FROM user_with_class WHERE class_id = ?`
	rows, err := repo.db.QueryContext(ctx, query, classID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get uid and weight")
	}
	defer rows.Close()

	uidAndWeight := make(map[int64]int)
	for rows.Next() {
		var uid int64
		var weight int
		if err = rows.Scan(&uid, &weight); err != nil {
			return nil, errors.Wrap(err, "failed to scan uid and weight")
		}
		uidAndWeight[uid] = weight
	}
	return uidAndWeight, nil
}

func (repo *CallRepository) GetClassUids(ctx context.Context, classID int64) ([]int64, error) {
	query := `SELECT uid FROM user_with_class WHERE class_id = ?`
	rows, err := repo.db.QueryContext(ctx, query, classID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get class uids")
	}
	defer rows.Close()

	var uids []int64
	for rows.Next() {
		var uid int64
		if err = rows.Scan(&uid); err != nil {
			return nil, errors.Wrap(err, "failed to scan uid")
		}
		uids = append(uids, uid)
	}
	return uids, nil
}

func (repo *CallRepository) GetCallEvents(ctx context.Context, classID int64) ([]*call.CallEvent, error) {
	query := `SELECT id, call_event_name, class_id, class_name, caller_id, caller_name, start_time, end_time FROM call_event WHERE class_id = ? LIMIT 5`
	rows, err := repo.db.QueryContext(ctx, query, classID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get call events")
	}
	defer rows.Close()

	var events []*call.CallEvent
	for rows.Next() {
		var event call.CallEvent
		if err = rows.Scan(&event.ID, &event.CallEventName, &event.ClassID, &event.ClassName,
			&event.CallerID, &event.CallerName, &event.StartTime, &event.EndTime); err != nil {
			return nil, errors.Wrap(err, "failed to scan call event")
		}
		events = append(events, &event)
	}
	return events, nil
}
