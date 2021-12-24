package loaders

// loader.exam.go
/**
 * 	This file is a part of controllers, used to load exam data
 */

import (
	"DatabaseCamp/models/general"
	"DatabaseCamp/models/storages"
	"DatabaseCamp/repositories"
	"sync"
)

/**
 * This class load concurrency exam data
 */
type examLoader struct {
	examRepo repositories.IExamRepository // repository for load activities of the exam
	userRepo repositories.IUserRepository // repository load corrected badges of the user

	correctedBadgeDB []storages.CorrectedBadgeDB // corrected badges of the user from the database
	examActivitiesDB []storages.ExamActivityDB   // activities of the exam from the database
}

/**
 * Constructor creates a new examLoader instance
 *
 * @param   examRepo    	Exam Repository for load activities of the exam
 * @param   userRepo        User Repository for load corrected badges of the user
 *
 * @return 	instance of examLoader
 */
func NewExamLoader(examRepo repositories.IExamRepository, userRepo repositories.IUserRepository) *examLoader {
	return &examLoader{examRepo: examRepo, userRepo: userRepo}
}

/**
 * Getter for getting correctedBadgeDB
 *
 * @return correctedBadgeDB
 */
func (l *examLoader) GetCorrectedBadgeDB() []storages.CorrectedBadgeDB {
	return l.correctedBadgeDB
}

/**
 * Getter for getting examActivitiesDB
 *
 * @return examActivitiesDB
 */
func (l *examLoader) GetExamActivitiesDB() []storages.ExamActivityDB {
	return l.examActivitiesDB
}

/**
 * Load concurrency exam data from the database
 *
 * @param   userID     		User ID for getting corrected badges of the user
 * @param   examID    		Exam ID for getting activities of the exam
 *
 * @return the error of loading data
 */
func (l *examLoader) Load(userID int, examID int) error {
	var wg sync.WaitGroup
	var err error
	concurrent := general.Concurrent{Wg: &wg, Err: &err}
	wg.Add(2)
	go l.loadCorrectedBadgeAsync(&concurrent, userID)
	go l.loadExamActivityAsync(&concurrent, examID)
	wg.Wait()
	return err
}

/**
 * Load corrected badges of the user from the database
 *
 * @param   userID     		User ID for getting corrected badges of the use
 */
func (l *examLoader) loadCorrectedBadgeAsync(concurrent *general.Concurrent, userID int) {
	defer concurrent.Wg.Done()
	result, err := l.userRepo.GetCollectedBadge(userID)
	if err != nil {
		*concurrent.Err = err
	}
	l.correctedBadgeDB = append(l.correctedBadgeDB, result...)
}

/**
 * Load activities of the exam from the database
 *
 * @param   concurrent     	Concurrent model for doing load concurrency
 * @param   examID    		Exam ID for getting activities of the exam
 */
func (l *examLoader) loadExamActivityAsync(concurrent *general.Concurrent, examID int) {
	defer concurrent.Wg.Done()
	result, err := l.examRepo.GetExamActivity(examID)
	if err != nil {
		*concurrent.Err = err
	}
	l.examActivitiesDB = append(l.examActivitiesDB, result...)
}
