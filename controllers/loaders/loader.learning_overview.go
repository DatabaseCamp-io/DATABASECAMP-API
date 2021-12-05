package loaders

// loader.learning_overview.go
/**
 * 	This file is a part of controllers, used to load concurrency learning overview data
 */

import (
	"DatabaseCamp/models/general"
	"DatabaseCamp/models/storages"
	"DatabaseCamp/repositories"
	"sync"
)

/**
 * This class load concurrency learning overview data
 */
type learningOverviewLoader struct {
	learningRepo repositories.ILearningRepository // repository for load learning overview data
	userRepo     repositories.IUserRepository     // repository for load learning progression of the user data

	overviewDB            []storages.OverviewDB            // learning overview data from the database
	learningProgressionDB []storages.LearningProgressionDB // learning progression of the user from the
}

/**
 * Constructor creates a new activityLoader instance
 *
 * @param   learningRepo    Learning Repository for load learning overview data
 * @param   userRepo        User Repository for load learning progression of the user data
 *
 * @return 	instance of activityLoader
 */
func NewLearningOverviewLoader(learningRepo repositories.ILearningRepository, userRepo repositories.IUserRepository) *learningOverviewLoader {
	return &learningOverviewLoader{learningRepo: learningRepo, userRepo: userRepo}
}

/**
 * Getter for getting overviewDB
 *
 * @return overviewDB
 */
func (l *learningOverviewLoader) GetOverviewDB() []storages.OverviewDB {
	return l.overviewDB
}

/**
 * Getter for getting learningProgressionDB
 *
 * @return learningProgressionDB
 */
func (l *learningOverviewLoader) GetLearningProgressionDB() []storages.LearningProgressionDB {
	return l.learningProgressionDB
}

/**
 * Load concurrency learning overview data from the database
 *
 * @param   userID     		User ID for getting learning progression of the user
 *
 * @return the error of loading data
 */
func (l *learningOverviewLoader) Load(userID int) error {
	var wg sync.WaitGroup
	var err error
	concurrent := general.Concurrent{Wg: &wg, Err: &err}
	wg.Add(2)
	go l.loadOverviewAsync(&concurrent)
	go l.loadLearningProgressionAsync(&concurrent, userID)
	wg.Wait()
	return err
}

/**
 * Load learning overview data from the database
 *
 * @param   concurrent     	Concurrent model for doing load concurrency
 */
func (l *learningOverviewLoader) loadOverviewAsync(concurrent *general.Concurrent) {
	defer concurrent.Wg.Done()
	result, err := l.learningRepo.GetOverview()
	if err != nil {
		*concurrent.Err = err
	}
	l.overviewDB = append(l.overviewDB, result...)
}

/**
 * Load activity data from the database
 *
 * @param   concurrent     	Concurrent model for doing load concurrency
 * @param   userID     		User ID for getting learning progression of the user
 */
func (l *learningOverviewLoader) loadLearningProgressionAsync(concurrent *general.Concurrent, id int) {
	defer concurrent.Wg.Done()
	result, err := l.userRepo.GetLearningProgression(id)
	if err != nil {
		*concurrent.Err = err
	}
	l.learningProgressionDB = append(l.learningProgressionDB, result...)
}
