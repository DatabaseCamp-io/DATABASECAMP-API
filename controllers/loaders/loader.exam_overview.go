package loaders

// loader.exam_overview.go
/**
 * 	This file is a part of controllers, used to load concurrency overview of the exam
 */

import (
	"DatabaseCamp/models/general"
	"DatabaseCamp/models/storages"
	"DatabaseCamp/repositories"
	"sync"
)

/**
 * This class load concurrency overview of the exam
 */
type examOverviewLoader struct {
	examRepo repositories.IExamRepository // repository for load all exam data
	userRepo repositories.IUserRepository // repository load corrected badges and exam results of the user

	correctedBadgeDB []storages.CorrectedBadgeDB // corrected badges of the user from the database
	examDB           []storages.ExamDB           // exam data from the database
	examResultsDB    []storages.ExamResultDB     // exam results of the user from the database
}

/**
 * Constructor creates a new examOverviewLoader instance
 *
 * @param   examRepo    	Exam Repository for load all exam data
 * @param   userRepo        User Repository for load corrected badges and exam results of the user
 *
 * @return 	instance of examOverviewLoader
 */
func NewExamOverviewLoader(examRepo repositories.IExamRepository, userRepo repositories.IUserRepository) *examOverviewLoader {
	return &examOverviewLoader{examRepo: examRepo, userRepo: userRepo}
}

/**
 * Getter for getting correctedBadgeDB
 *
 * @return correctedBadgeDB
 */
func (l *examOverviewLoader) GetCorrectedBadgeDB() []storages.CorrectedBadgeDB {
	return l.correctedBadgeDB
}

/**
 * Getter for getting examDB
 *
 * @return examDB
 */
func (l *examOverviewLoader) GetExamDB() []storages.ExamDB {
	return l.examDB
}

/**
 * Getter for getting examResultsDB
 *
 * @return examResultsDB
 */
func (l *examOverviewLoader) GetExamResultsDB() []storages.ExamResultDB {
	return l.examResultsDB
}

/**
 * Load concurrency all activity data from the database
 *
 * @param   userID     		User ID for getting corrected badges and exam results of the user
 *
 * @return the error of loading data
 */
func (l *examOverviewLoader) Load(userID int) error {
	var wg sync.WaitGroup
	var err error
	concurrent := general.Concurrent{Wg: &wg, Err: &err}
	wg.Add(3)
	go l.loadCorrectedBadgeAsync(&concurrent, userID)
	go l.loadExamAsync(&concurrent)
	go l.loadExamResultAsync(&concurrent, userID)
	wg.Wait()
	return err
}

/**
 * Load exam results of the user from the database
 *
 * @param   concurrent     	Concurrent model for doing load concurrency
 * @param   userID     		User ID for getting exam results of the user
 *
 * @return the error of loading data
 */
func (l *examOverviewLoader) loadExamResultAsync(concurrent *general.Concurrent, userID int) {
	defer concurrent.Wg.Done()
	result, err := l.userRepo.GetExamResult(userID)
	if err != nil {
		*concurrent.Err = err
	}
	l.examResultsDB = append(l.examResultsDB, result...)
}

/**
 * Load all exam  from the database
 *
 * @param   concurrent     	Concurrent model for doing load concurrency
 */
func (l *examOverviewLoader) loadExamAsync(concurrent *general.Concurrent) {
	defer concurrent.Wg.Done()
	result, err := l.examRepo.GetExamOverview()
	if err != nil {
		*concurrent.Err = err
	}
	l.examDB = append(l.examDB, result...)
}

/**
 * Load corrected badges of the user from the database
 *
 * @param   concurrent     	Concurrent model for doing load concurrency
 * @param   userID     		User ID for getting corrected badges of the user
 */
func (l *examOverviewLoader) loadCorrectedBadgeAsync(concurrent *general.Concurrent, userID int) {
	defer concurrent.Wg.Done()
	result, err := l.userRepo.GetCollectedBadge(userID)
	if err != nil {
		*concurrent.Err = err
	}
	l.correctedBadgeDB = append(l.correctedBadgeDB, result...)
}
