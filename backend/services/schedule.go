package services

import (
	"context"
	"errors"
	"gobot/backend/dao"
	"gobot/backend/log"
	"gobot/backend/models"
	"time"

	"github.com/google/uuid"
	"github.com/reugn/go-quartz/logger"
	"github.com/reugn/go-quartz/quartz"
)

var sched quartz.Scheduler

type ScheduleJob struct {
	*models.Schedule
}

func (s ScheduleJob) Execute(ctx context.Context) error {
	project, err := QueryProjectById(s.ProjectId)
	if err != nil {
		return err
	}
	if project.IsFlow {
		return errors.New("暂不支持流程图类型的项目运行")
	} else {
		return RunSubFlow(ctx, project.Id, "main")
	}
}

func (s ScheduleJob) Description() string {
	return s.Desc
}

func InitSchedule(ctx context.Context) {
	logger.SetDefault(logger.NewSimpleLogger(nil, logger.LevelOff))
	sched = quartz.NewStdScheduler()
	sched.Start(ctx)
	if err := dao.ReadTx(dao.MainDB, func(tx dao.Tx) error {
		allSchedules, err := tx.SelectAllSchedules()
		if err != nil {
			return err
		}
		for _, schedule := range allSchedules {
			if schedule.Enabled {
				cronTrigger, _ := quartz.NewCronTrigger(schedule.Cron)
				if err = sched.ScheduleJob(quartz.NewJobDetail(ScheduleJob{schedule}, quartz.NewJobKey(schedule.Id)),
					cronTrigger); err != nil {
					log.Logger.Error("定时任务" + schedule.Id + "启动失败")
				}
			}
		}
		return nil
	}); err != nil {
		log.Logger.Error("定时任务启动失败")
	} else {
		log.Logger.Info("定时任务启动成功")
	}
}

func QuerySchedulePage() (resultList []*models.Schedule, err error) {
	if err := dao.ReadTx(dao.MainDB, func(tx dao.Tx) error {
		resultList, err = tx.SelectAllSchedules()
		return err
	}); err != nil {
		return nil, err
	}
	return resultList, nil
}

func ToggleScheduleById(id string, state bool) (err error) {
	if err := dao.WriteTx(dao.MainDB, func(tx dao.Tx) error {
		result, err := tx.SelectSchedule(id)
		if err != nil {
			return err
		}
		result.Enabled = state
		if err := tx.UpdateSchedule(result); err != nil {
			return err
		}
		if result.Enabled {
			cronTrigger, _ := quartz.NewCronTrigger(result.Cron)
			if err = sched.ScheduleJob(quartz.NewJobDetail(ScheduleJob{result}, quartz.NewJobKey(result.Id)),
				cronTrigger); err != nil {
				log.Logger.Error("定时任务" + result.Id + "启动失败")
				return err
			}
		} else {
			if err = sched.DeleteJob(quartz.NewJobKey(result.Id)); err != nil {
				return err
			}
		}
		return nil
	}); err != nil {
		return err
	}
	return nil
}

func AddOrUpdateSchedule(schedule models.Schedule) (err error) {
	if err = quartz.ValidateCronExpression(schedule.Cron); err != nil {
		return err
	}
	if schedule.Id == "" {
		schedule.Id = uuid.New().String()
		if err := dao.WriteTx(dao.MainDB, func(tx dao.Tx) error {
			if err := tx.InsertSchedule(&schedule); err != nil {
				return err
			}

			return nil
		}); err != nil {
			return err
		}
	} else {
		if err := dao.WriteTx(dao.MainDB, func(tx dao.Tx) error {
			old, err := tx.SelectSchedule(schedule.Id)
			if err != nil {
				return err
			}
			if old.Enabled {
				if err = sched.DeleteJob(quartz.NewJobKey(old.Id)); err != nil {
					return err
				}
			}
			if err := tx.UpdateSchedule(&schedule); err != nil {
				return err
			}
			return nil
		}); err != nil {
			return err
		}
	}
	if schedule.Enabled {
		cronTrigger, _ := quartz.NewCronTrigger(schedule.Cron)
		if err = sched.ScheduleJob(quartz.NewJobDetail(ScheduleJob{&schedule}, quartz.NewJobKey(schedule.Id)),
			cronTrigger); err != nil {
			log.Logger.Error("定时任务" + schedule.Id + "启动失败")
		}
	}
	return nil
}

func RemoveSchedule(id string) (err error) {
	if err := dao.WriteTx(dao.MainDB, func(tx dao.Tx) error {
		if err := tx.DeleteSchedule(id); err != nil {
			return err
		}
		if err = sched.DeleteJob(quartz.NewJobKey(id)); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return err
	}
	return nil
}

func GetNextTriggerTime(cron string) (string, error) {
	trigger, err := quartz.NewCronTrigger(cron)
	if err != nil {
		return "", err
	}
	result := ""
	now := time.Now().UnixNano()
	for i := 0; i < 5; i++ {
		nextFireTime, err := trigger.NextFireTime(now)
		if err != nil {
			return "", err
		}
		t := time.Unix(0, nextFireTime)
		result += t.UTC().Format("2006-01-02 15:04:05") + "  "
		now = nextFireTime
	}

	return result, nil
}
