package loaders

import (
	"DatabaseCamp/models/general"
	"DatabaseCamp/repositories"
	"sync"
)

type examOverviewLoader struct {
	examRepo         repositories.IExamRepository
	userRepo         repositories.IUserRepository
	CorrectedBadgeDB []general.CorrectedBadgeDB
	ExamDB           []general.ExamDB
	ExamResultsDB    []general.ExamResultDB
}

func NewExamOverviewLoader(examRepo repositories.IExamRepository, userRepo repositories.IUserRepository) *examOverviewLoader {
	return &examOverviewLoader{examRepo: examRepo, userRepo: userRepo}
}

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

func (l *examOverviewLoader) loadExamResultAsync(concurrent *general.Concurrent, userID int) {
	defer concurrent.Wg.Done()
	result, err := l.userRepo.GetExamResult(userID)
	if err != nil {
		*concurrent.Err = err
	}
	l.ExamResultsDB = append(l.ExamResultsDB, result...)
}

func (l *examOverviewLoader) loadExamAsync(concurrent *general.Concurrent) {
	defer concurrent.Wg.Done()
	result, err := l.examRepo.GetExamOverview()
	if err != nil {
		*concurrent.Err = err
	}
	l.ExamDB = append(l.ExamDB, result...)
}

func (l *examOverviewLoader) loadCorrectedBadgeAsync(concurrent *general.Concurrent, userID int) {
	defer concurrent.Wg.Done()
	result, err := l.userRepo.GetCollectedBadge(userID)
	if err != nil {
		*concurrent.Err = err
	}
	l.CorrectedBadgeDB = append(l.CorrectedBadgeDB, result...)
}
