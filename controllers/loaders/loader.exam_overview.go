package loaders

import (
	"DatabaseCamp/models/general"
	"DatabaseCamp/models/storages"
	"DatabaseCamp/repositories"
	"sync"
)

type examOverviewLoader struct {
	examRepo         repositories.IExamRepository
	userRepo         repositories.IUserRepository
	correctedBadgeDB []storages.CorrectedBadgeDB
	examDB           []storages.ExamDB
	examResultsDB    []storages.ExamResultDB
}

func NewExamOverviewLoader(examRepo repositories.IExamRepository, userRepo repositories.IUserRepository) *examOverviewLoader {
	return &examOverviewLoader{examRepo: examRepo, userRepo: userRepo}
}

func (l *examOverviewLoader) GetCorrectedBadgeDB() []storages.CorrectedBadgeDB {
	return l.correctedBadgeDB
}

func (l *examOverviewLoader) GetExamDB() []storages.ExamDB {
	return l.examDB
}

func (l *examOverviewLoader) GetExamResultsDB() []storages.ExamResultDB {
	return l.examResultsDB
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
	l.examResultsDB = append(l.examResultsDB, result...)
}

func (l *examOverviewLoader) loadExamAsync(concurrent *general.Concurrent) {
	defer concurrent.Wg.Done()
	result, err := l.examRepo.GetExamOverview()
	if err != nil {
		*concurrent.Err = err
	}
	l.examDB = append(l.examDB, result...)
}

func (l *examOverviewLoader) loadCorrectedBadgeAsync(concurrent *general.Concurrent, userID int) {
	defer concurrent.Wg.Done()
	result, err := l.userRepo.GetCollectedBadge(userID)
	if err != nil {
		*concurrent.Err = err
	}
	l.correctedBadgeDB = append(l.correctedBadgeDB, result...)
}
