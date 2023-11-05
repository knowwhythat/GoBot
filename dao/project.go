package dao

import (
	"encoding/json"
	"gobot/constants"
	"gobot/models"
	"time"

	uuid "github.com/google/uuid"
)

func (t *Tx) SelectAllProject() (projects []*models.Project, err error) {
	b := t.tx.Bucket([]byte(constants.ProjectBucket))
	if b == nil {
		return projects, nil
	}
	c := b.Cursor()
	for k, v := c.First(); k != nil; k, v = c.Next() {
		var project *models.Project
		if err = json.Unmarshal(v, &project); err != nil {
			return nil, err
		}
		projects = append(projects, project)
	}
	return projects, nil
}

func (t *Tx) SelectProject(id string) (project *models.Project, err error) {
	b := t.tx.Bucket([]byte(constants.ProjectBucket))
	if b == nil {
		return nil, nil
	}
	value := b.Get([]byte(id))
	if len(value) == 0 {
		return nil, nil
	}
	if err = json.Unmarshal(value, &project); err != nil {
		return nil, err
	}
	return project, nil
}

func (t *Tx) InsertProject(project *models.Project) (err error) {
	if project.Id == uuid.Nil {
		project.Id = uuid.New()
	}
	if project.CreateTs.IsZero() {
		project.CreateTs = time.Now()
	}
	if project.UpdateTs.IsZero() {
		project.UpdateTs = time.Now()
	}

	value, err := json.Marshal(&project)
	if err != nil {
		return err
	}
	b, err := t.tx.CreateBucketIfNotExists([]byte(constants.ProjectBucket))
	if err != nil {
		return err
	}
	if err = b.Put([]byte(project.Id.String()), value); err != nil {
		return err
	}
	return nil
}

func (t *Tx) UpdateProject(project *models.Project) (err error) {
	b := t.tx.Bucket([]byte(constants.ProjectBucket))
	if b == nil {
		return nil
	}
	project.UpdateTs = time.Now()
	value, _ := json.Marshal(&project)
	if err = b.Put([]byte(project.Id.String()), value); err != nil {
		return err
	}
	return nil
}

func (t *Tx) DeleteProject(id string) (err error) {
	b := t.tx.Bucket([]byte(constants.ProjectBucket))
	if b == nil {
		return nil
	}
	if err = b.Delete([]byte(id)); err != nil {
		return err
	}
	return nil
}
