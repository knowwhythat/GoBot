package dao

import (
	"encoding/json"
	"gobot/backend/constants"
	"gobot/backend/models"
	"sort"
	"strings"
	"time"

	uuid "github.com/google/uuid"
)

// SelectAllSchedules 查询所有定时调度
func (t *Tx) SelectAllSchedules() (schedules []*models.Schedule, err error) {
	b := t.tx.Bucket([]byte(constants.ScheduleBucket))
	if b == nil {
		return schedules, nil
	}
	c := b.Cursor()
	for k, v := c.First(); k != nil; k, v = c.Next() {
		var schedule *models.Schedule
		if err = json.Unmarshal(v, &schedule); err != nil {
			return nil, err
		}
		schedules = append(schedules, schedule)
	}
	sort.Slice(schedules, func(i, j int) bool {
		return schedules[i].UpdateTs.After(schedules[j].UpdateTs)
	})
	return schedules, nil
}

// SelectSchedule 根据 ID 查询定时调度
func (t *Tx) SelectSchedule(id string) (schedule *models.Schedule, err error) {
	b := t.tx.Bucket([]byte(constants.ScheduleBucket))
	if b == nil {
		return nil, nil
	}
	value := b.Get([]byte(id))
	if len(value) == 0 {
		return nil, nil
	}
	if err = json.Unmarshal(value, &schedule); err != nil {
		return nil, err
	}
	return schedule, nil
}

// InsertSchedule 插入新定时调度
func (t *Tx) InsertSchedule(schedule *models.Schedule) (err error) {
	if strings.TrimSpace(schedule.Id) == "" {
		schedule.Id = uuid.New().String()
	}
	if schedule.CreateTs.IsZero() {
		schedule.CreateTs = time.Now()
	}
	if schedule.UpdateTs.IsZero() {
		schedule.UpdateTs = time.Now()
	}

	value, err := json.Marshal(&schedule)
	if err != nil {
		return err
	}
	b, err := t.tx.CreateBucketIfNotExists([]byte(constants.ScheduleBucket))
	if err != nil {
		return err
	}
	if err = b.Put([]byte(schedule.Id), value); err != nil {
		return err
	}
	return nil
}

// UpdateSchedule 更新定时调度
func (t *Tx) UpdateSchedule(schedule *models.Schedule) (err error) {
	b := t.tx.Bucket([]byte(constants.ScheduleBucket))
	if b == nil {
		return nil
	}
	schedule.UpdateTs = time.Now()
	value, _ := json.Marshal(&schedule)
	if err = b.Put([]byte(schedule.Id), value); err != nil {
		return err
	}
	return nil
}

// DeleteSchedule 通过 ID 删除定时调度
func (t *Tx) DeleteSchedule(id string) (err error) {
	b := t.tx.Bucket([]byte(constants.ScheduleBucket))
	if b == nil {
		return nil
	}
	if err = b.Delete([]byte(id)); err != nil {
		return err
	}
	return nil
}

// DeleteAllSchedulesWhereSpiderId 根据项目 ID 删除所有定时调度
func (t *Tx) DeleteAllSchedulesWhereProjectId(projectId string) (err error) {
	b := t.tx.Bucket([]byte(constants.ScheduleBucket))
	if b == nil {
		return nil
	}
	c := b.Cursor()
	for k, v := c.First(); k != nil; k, v = c.Next() {
		var schedule *models.Schedule
		if err = json.Unmarshal(v, &schedule); err != nil {
			return err
		}
		if schedule.ProjectId == projectId {
			if err = b.Delete([]byte(schedule.Id)); err != nil {
				return err
			}
		}
	}
	return nil
}
