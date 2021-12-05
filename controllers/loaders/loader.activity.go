package loaders

// loader.activity.go
/**
 * 	This file is a part of controller, used to load concurrency activity data
 */

import (
	"DatabaseCamp/models/general"
	"DatabaseCamp/models/storages"
	"DatabaseCamp/repositories"
	"sync"
)

/**
 * This class load concurrency all activity data
 */
type activityLoader struct {
	learningRepo repositories.ILearningRepository // repository for load activity data
	userRepo     repositories.IUserRepository     // repository for load user hints

	activityDB      *storages.ActivityDB  // activity data from the database
	activityHintsDB []storages.HintDB     // activity hints from the database
	userHintsDB     []storages.UserHintDB // user hints from the database
}

/**
 * Constructor creates a new activityLoader instance
 * @param   learningRepo    Learning Repository for load learning data
 * @param   userRepo        User Repository for load user hints
 * @return 	instance of activityLoader
 */
func NewActivityLoader(learningRepo repositories.ILearningRepository, userRepo repositories.IUserRepository) *activityLoader {
	return &activityLoader{learningRepo: learningRepo, userRepo: userRepo}
}

/**
 * Getter for getting activityDB
 * @return activityDB
 */
func (l *activityLoader) GetActivityDB() *storages.ActivityDB {
	return l.activityDB
}

/**
 * Getter for getting activityHintsDB
 * @return activityHintsDB
 */
func (l *activityLoader) GetActivityHintsDB() []storages.HintDB {
	return l.activityHintsDB
}

/**
 * Getter for getting userHintsDB
 * @return userHintsDB
 */
func (l *activityLoader) GetUserHintsDB() []storages.UserHintDB {
	return l.userHintsDB
}

/**
 * load concurrency all activity data from the database
 * @param   userID     		User ID for getting user hints information of the activity
 * @param   activityID    	Activity ID for getting activity data
 * @return the error of loading data
 */
func (l *activityLoader) Load(userID int, activityID int) error {
	var wg sync.WaitGroup
	var err error
	concurrent := general.Concurrent{Wg: &wg, Err: &err}
	wg.Add(3)
	go l.loadActivityAsync(&concurrent, activityID)
	go l.loadActivityHints(&concurrent, activityID)
	go l.loadUserHintsAsync(&concurrent, userID, activityID)
	wg.Wait()
	return err
}

/**
 * load activity data from the database
 * @param   concurrent     	Concurrent model for doing load concurrency
 * @param   activityID    	Activity ID for getting activity data
 */
func (l *activityLoader) loadActivityAsync(concurrent *general.Concurrent, activityID int) {
	defer concurrent.Wg.Done()
	var err error
	l.activityDB, err = l.learningRepo.GetActivity(activityID)
	if err != nil {
		*concurrent.Err = err
	}
}

/**
 * load activity data from the database
 * @param   concurrent     	Concurrent model for doing load concurrency
 * @param   userID     		User ID for getting user hints information of the activity
 * @param   activityID    	Activity ID for getting activity data
 */
func (l *activityLoader) loadUserHintsAsync(concurrent *general.Concurrent, userID int, activityID int) {
	defer concurrent.Wg.Done()
	result, e := l.userRepo.GetUserHint(userID, activityID)
	if e != nil {
		*concurrent.Err = e
	}
	l.userHintsDB = append(l.userHintsDB, result...)
}

/**
 * load activity data from the database
 * @param   concurrent     	Concurrent model for doing load concurrency
 * @param   activityID    	Activity ID for getting activity data
 */
func (l *activityLoader) loadActivityHints(concurrent *general.Concurrent, activityID int) {
	defer concurrent.Wg.Done()
	result, e := l.learningRepo.GetActivityHints(activityID)
	if e != nil {
		*concurrent.Err = e
	}
	l.activityHintsDB = append(l.activityHintsDB, result...)
}
