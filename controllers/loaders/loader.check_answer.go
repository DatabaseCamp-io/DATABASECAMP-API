package loaders

// loader.check_answer.go
/**
 * 	This file is a part of controller, used to load concurrency activity answer data
 */

import (
	"DatabaseCamp/models/general"
	"DatabaseCamp/models/storages"
	"DatabaseCamp/repositories"
	"sync"
)

/**
 * This class load concurrency activity answer data
 */
type checkAnswerLoader struct {
	learningRepo repositories.ILearningRepository // repository for load activity data

	choicesDB  interface{}          // choices of the activity from the database
	activityDB *storages.ActivityDB // activity data from the database
}

/**
 * Constructor creates a new checkAnswerLoader instance
 *
 * @param   learningRepo Learning Repository for load learning data
 *
 * @return 	instance of checkAnswerLoader
 */
func NewCheckAnswerLoader(learningRepo repositories.ILearningRepository) *checkAnswerLoader {
	return &checkAnswerLoader{learningRepo: learningRepo}
}

/**
 * Getter for getting choicesDB
 *
 * @return choicesDB
 */
func (c *checkAnswerLoader) GetChoicesDB() interface{} {
	return c.choicesDB
}

/**
 * Getter for getting activityDB

 * * @return activityDB
 */
func (c *checkAnswerLoader) GetActivityDB() *storages.ActivityDB {
	return c.activityDB
}

/**
 * Load concurrency activity answer data
 *
 * @param   activityID    		Activity ID for getting activity data
 * @param   activityTypeID    	Activity Type ID for getting choices of activity by type
 * @param   getChoicesFunc    	function for getting choices
 *
 * @return the error of loading data
 */
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

/**
 * Load activity data from the database
 *
 * @param   concurrent     	Concurrent model for doing load concurrency
 * @param   activityID    	Activity ID for getting activity data
 */
func (c *checkAnswerLoader) loadActivityAsync(concurrent *general.Concurrent, activityID int) {
	defer concurrent.Wg.Done()
	var err error
	c.activityDB, err = c.learningRepo.GetActivity(activityID)
	if err != nil {
		*concurrent.Err = err
	}
}

/**
 * call getChoicesFunc for getting activity choices
 *
 * @param   concurrent     		Concurrent model for doing load concurrency
 * @param   activityID    		Activity ID for getting activity data
 * @param   activityTypeID    	Activity Type ID for getting choices of activity by type
 * @param   getChoicesFunc    	function for getting choices
 */
func (c *checkAnswerLoader) getChioceAsync(
	concurrent *general.Concurrent,
	activityID int,
	activityTypeID int,
	getChoicesFunc func(activityID int, activityTypeID int) (interface{}, error),
) {
	defer concurrent.Wg.Done()
	var err error
	c.choicesDB, err = getChoicesFunc(activityID, activityTypeID)
	if err != nil {
		*concurrent.Err = err
	}
}
