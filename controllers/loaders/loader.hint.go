package loaders

// loader.hint.go
/**
 * 	This file is a part of controllers, used to load concurrency hints data
 */

import (
	"DatabaseCamp/models/general"
	"DatabaseCamp/models/storages"
	"DatabaseCamp/repositories"
	"sync"
)

/**
 * This class load concurrency hints data
 */
type hintLoader struct {
	learningRepo repositories.ILearningRepository // repository for hints data
	userRepo     repositories.IUserRepository     // repository for load user data

	activityHintsDB []storages.HintDB     // hints of the activity from the database
	userHintsDB     []storages.UserHintDB // user hints from the database
	userDB          *storages.UserDB      // user data from the database
}

/**
 * Constructor creates a new hintLoader instance
 *
 * @param   learningRepo    Learning Repository for load hints data
 * @param   userRepo        User Repository for load user data
 *
 * @return 	instance of hintLoader
 */
func NewHintLoader(learningRepo repositories.ILearningRepository, userRepo repositories.IUserRepository) *hintLoader {
	return &hintLoader{learningRepo: learningRepo, userRepo: userRepo}
}

/**
 * Getter for getting activityHintsDB
 *
 * @return activityHintsDB
 */
func (l *hintLoader) GetActivityHintsDB() []storages.HintDB {
	return l.activityHintsDB
}

/**
 * Getter for getting userHintsDB
 *
 * @return userHintsDB
 */
func (l *hintLoader) GetUserHintsDB() []storages.UserHintDB {
	return l.userHintsDB
}

/**
 * Getter for getting userDB
 *
 * @return userDB
 */
func (l *hintLoader) GetUserDB() *storages.UserDB {
	return l.userDB
}

/**
 * Load concurrency all activity data from the database

 * @param   userID     		User ID for getting user data
 * @param   activityID    	Activity ID for getting hints data
 *
 * @return the error of loading data
 */
func (l *hintLoader) Load(userID int, activityID int) error {
	var wg sync.WaitGroup
	var err error
	concurrent := general.Concurrent{Wg: &wg, Err: &err}
	wg.Add(3)
	go l.loadActivityHints(&concurrent, activityID)
	go l.loadUserHintsAsync(&concurrent, userID, activityID)
	go l.loadUser(&concurrent, userID)
	wg.Wait()
	return err
}

/**
 * Load user data from the database
 *
 * @param   concurrent     	Concurrent model for doing load concurrency
 * @param   userID    		User ID for getting user data
 */
func (l *hintLoader) loadUser(concurrent *general.Concurrent, userID int) {
	defer concurrent.Wg.Done()
	var err error
	l.userDB, err = l.userRepo.GetUserByID(userID)
	if err != nil {
		*concurrent.Err = err
	}
}

/**
 * Load user hints from the database
 *
 * @param   concurrent     	Concurrent model for doing load concurrency
 * @param   userID    		User ID for getting user hints of the activity data
 * @param   activityID    	Activity ID for indicate activity
 */
func (l *hintLoader) loadUserHintsAsync(concurrent *general.Concurrent, userID int, activityID int) {
	defer concurrent.Wg.Done()
	result, e := l.userRepo.GetUserHint(userID, activityID)
	if e != nil {
		*concurrent.Err = e
	}
	l.userHintsDB = append(l.userHintsDB, result...)
}

/**
 * Load hints of the activity from the database
 *
 * @param   concurrent     	Concurrent model for doing load concurrency
 * @param   activityID    	Activity ID for getting hints of the activity
 */
func (l *hintLoader) loadActivityHints(concurrent *general.Concurrent, activityID int) {
	defer concurrent.Wg.Done()
	result, e := l.learningRepo.GetActivityHints(activityID)
	if e != nil {
		*concurrent.Err = e
	}
	l.activityHintsDB = append(l.activityHintsDB, result...)
}
