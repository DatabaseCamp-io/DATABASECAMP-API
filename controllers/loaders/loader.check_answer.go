package loaders

import (
	"DatabaseCamp/models/general"
	"DatabaseCamp/repositories"
	"sync"
)

type checkAnswerLoader struct {
	learningRepo repositories.ILearningRepository
	ChoicesDB    interface{}
	ActivityDB   *general.ActivityDB
}

func NewCheckAnswerLoader(learningRepo repositories.ILearningRepository) *checkAnswerLoader {
	return &checkAnswerLoader{learningRepo: learningRepo}
}

func (c *checkAnswerLoader) Load(activityID int, activityTypeID int, getChoicesFunc func(activityID int, activityTypeID int) (interface{}, error)) error {
	var wg sync.WaitGroup
	var err error
	concurrent := general.Concurrent{Wg: &wg, Err: &err}
	wg.Add(2)
	go c.loadActivityAsync(&concurrent, activityID)
	go c.getChioceAsync(&concurrent, activityID, activityTypeID, getChoicesFunc)
	wg.Wait()
	return err
}

func (c *checkAnswerLoader) loadActivityAsync(concurrent *general.Concurrent, activityID int) {
	defer concurrent.Wg.Done()
	var err error
	c.ActivityDB, err = c.learningRepo.GetActivity(activityID)
	if err != nil {
		*concurrent.Err = err
	}
}

func (c *checkAnswerLoader) getChioceAsync(
	concurrent *general.Concurrent,
	activityID int,
	activityTypeID int,
	getChoicesFunc func(activityID int, activityTypeID int) (interface{}, error),
) {
	defer concurrent.Wg.Done()
	var err error
	c.ChoicesDB, err = getChoicesFunc(activityID, activityTypeID)
	if err != nil {
		*concurrent.Err = err
	}
}
