package services

import (
	"gobot/backend/dao"
	"gobot/backend/forms"
)

func QueryAllExecution() (resultList []*forms.ExecutionForm, err error) {
	if err := dao.ReadTx(dao.MainDB, func(tx dao.Tx) error {
		executions, err := tx.SelectAllExecutions()
		for _, execution := range executions {
			project, err := tx.SelectProject(execution.ProjectId)
			if err != nil {
				return err
			}
			projectName := ""
			if project != nil {
				projectName = project.Name
			}
			resultList = append(resultList, &forms.ExecutionForm{
				Id:            execution.Id,
				ProjectId:     execution.ProjectId,
				ProjectName:   projectName,
				SubFlowId:     execution.SubFlowId,
				TriggerType:   execution.TriggerType,
				ExecuteResult: execution.ExecuteResult,
				ErrorMsg:      execution.ErrorMsg,
				StartTs:       execution.StartTs,
				EndTs:         execution.EndTs,
			})
		}
		return err
	}); err != nil {
		return nil, err
	}
	return resultList, nil
}

func RemoveExecution(id string) (err error) {
	if err := dao.WriteTx(dao.MainDB, func(tx dao.Tx) error {
		if err := tx.DeleteExecution(id); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return err
	}

	if err := dao.WriteTx(dao.LogDB, func(tx dao.Tx) error {
		if err := tx.DeleteExecutionLog(id); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return err
	}
	return nil
}

func SelectExecutionLog(id string) (logs string, err error) {
	if err := dao.ReadTx(dao.LogDB, func(tx dao.Tx) error {
		logs, err = tx.SelectExecutionLog(id)
		if err != nil {
			return err
		}
		return nil
	}); err != nil {
		return "", err
	}
	return logs, nil
}
