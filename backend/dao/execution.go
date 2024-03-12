package dao

import (
	"encoding/json"
	"gobot/backend/constants"
	"gobot/backend/models"
	"time"

	uuid "github.com/google/uuid"
)

func (t *Tx) SelectAllExecutions() (executions []*models.Execution, err error) {
	b := t.tx.Bucket([]byte(constants.ExecutionBucket))
	if b == nil {
		return executions, nil
	}
	c := b.Cursor()
	for k, v := c.First(); k != nil; k, v = c.Next() {
		var execution *models.Execution
		if err = json.Unmarshal(v, &execution); err != nil {
			return nil, err
		}
		executions = append(executions, execution)
	}
	return executions, nil
}

func (t *Tx) SelectExecution(id uuid.UUID) (execution *models.Execution, err error) {
	b := t.tx.Bucket([]byte(constants.ExecutionBucket))
	if b == nil {
		return nil, nil
	}
	value := b.Get([]byte(id.String()))
	if len(value) == 0 {
		return nil, nil
	}
	if err = json.Unmarshal(value, &execution); err != nil {
		return nil, err
	}
	return execution, nil
}

func (t *Tx) InsertExecution(execution *models.Execution) (err error) {
	if execution.Id == uuid.Nil {
		execution.Id = uuid.New()
	}
	if execution.StartTs.IsZero() {
		execution.StartTs = time.Now()
	}
	if execution.EndTs.IsZero() {
		execution.EndTs = time.Now()
	}

	value, err := json.Marshal(&execution)
	if err != nil {
		return err
	}
	b, err := t.tx.CreateBucketIfNotExists([]byte(constants.ExecutionBucket))
	if err != nil {
		return err
	}
	if err = b.Put([]byte(execution.Id.String()), value); err != nil {
		return err
	}
	return nil
}

func (t *Tx) UpdateExecution(execution *models.Execution) (err error) {
	b := t.tx.Bucket([]byte(constants.ExecutionBucket))
	if b == nil {
		return nil
	}
	execution.EndTs = time.Now()
	value, _ := json.Marshal(&execution)
	if err = b.Put([]byte(execution.Id.String()), value); err != nil {
		return err
	}
	return nil
}

func (t *Tx) DeleteExecution(id uuid.UUID) (err error) {
	b := t.tx.Bucket([]byte(constants.ExecutionBucket))
	if b == nil {
		return nil
	}
	if err = b.Delete([]byte(id.String())); err != nil {
		return err
	}
	return nil
}

func (t *Tx) DeleteAllExecutionsWhereProjectId(ProjectId uuid.UUID) (err error) {
	b := t.tx.Bucket([]byte(constants.ExecutionBucket))
	if b == nil {
		return nil
	}
	c := b.Cursor()
	for k, v := c.First(); k != nil; k, v = c.Next() {
		var execution *models.Execution
		if err = json.Unmarshal(v, &execution); err != nil {
			return err
		}
		if execution.ProjectId == ProjectId {
			if err = b.Delete([]byte(execution.Id.String())); err != nil {
				return err
			}
		}
	}
	return nil
}

func (t *Tx) InsertExecutionLog(executionLogs string, executionId uuid.UUID) (err error) {
	b, err := t.tx.CreateBucketIfNotExists([]byte(constants.ExecutionLogBucket))
	if err != nil {
		return err
	}
	if err = b.Put([]byte(executionId.String()), []byte(executionLogs)); err != nil {
		return err
	}
	return nil
}

func (t *Tx) SelectExecutionLog(id uuid.UUID) (executionLog string, err error) {
	b := t.tx.Bucket([]byte(constants.ExecutionLogBucket))
	if b == nil {
		return "", nil
	}
	value := b.Get([]byte(id.String()))
	if len(value) == 0 {
		return "", nil
	}

	return string(value), nil
}
