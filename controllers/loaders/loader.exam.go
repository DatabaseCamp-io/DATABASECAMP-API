package loaders

import (
	"DatabaseCamp/models/general"
	"DatabaseCamp/repositories"
	"sync"
)

type examLoader struct {
	examRepo         repositories.IExamRepository
	userRepo         repositories.IUserRepository
	CorrectedBadgeDB []general.CorrectedBadgeDB
	ExamActivitiesDB []general.ExamActivityDB
}

func NewExamLoader(examRepo repositories.IExamRepository, userRepo repositories.IUserRepository) *examLoader {
	return &examLoader{examRepo: examRepo, userRepo: userRepo}
}

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

func (l *examLoader) loadCorrectedBadgeAsync(concurrent *general.Concurrent, userID int) {
	defer concurrent.Wg.Done()
	result, err := l.userRepo.GetCollectedBadge(userID)
	if err != nil {
		*concurrent.Err = err
	}
	l.CorrectedBadgeDB = append(l.CorrectedBadgeDB, result...)
}

func (l *examLoader) loadExamActivityAsync(concurrent *general.Concurrent, examID int) {
	defer concurrent.Wg.Done()
	result, err := l.examRepo.GetExamActivity(examID)
	if err != nil {
		*concurrent.Err = err
	}
	l.ExamActivitiesDB = append(l.ExamActivitiesDB, result...)
}
