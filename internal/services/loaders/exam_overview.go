package loaders

import (
	"database-camp/internal/models/entities/exam"
	"database-camp/internal/models/entities/user"
	"database-camp/internal/repositories"
	"sync"
)

type examOverviewLoader struct {
	examRepo repositories.ExamRepository
	userRepo repositories.UserRepository

	correctedBadges []user.CorrectedBadge
	exams           []exam.Exam
	examResults     []exam.ExamResult
}

func NewExamOverviewLoader(examRepo repositories.ExamRepository, userRepo repositories.UserRepository) *examOverviewLoader {
	return &examOverviewLoader{examRepo: examRepo, userRepo: userRepo}
}

func (l *examOverviewLoader) GetCorrectedBadges() user.CorrectedBadges {
	return l.correctedBadges
}

func (l *examOverviewLoader) GetExams() exam.Exams {
	return l.exams
}

func (l *examOverviewLoader) GetExamResults() exam.ExamResults {
	return l.examResults
}

func (l *examOverviewLoader) Load(userID int) error {
	var wg sync.WaitGroup
	var err error
	concurrent := Concurrent{Wg: &wg, Err: &err}
	wg.Add(3)
	go l.loadCorrectedBadgeAsync(&concurrent, userID)
	go l.loadExamAsync(&concurrent)
	go l.loadExamResultAsync(&concurrent, userID)
	wg.Wait()
	return err
}

func (l *examOverviewLoader) loadExamResultAsync(concurrent *Concurrent, userID int) {
	defer concurrent.Wg.Done()
	result, err := l.examRepo.GetExamResults(userID)
	if err != nil {
		*concurrent.Err = err
	}
	l.examResults = append(l.examResults, result...)
}

func (l *examOverviewLoader) loadExamAsync(concurrent *Concurrent) {
	defer concurrent.Wg.Done()
	result, err := l.examRepo.GetExams()
	if err != nil {
		*concurrent.Err = err
	}
	l.exams = append(l.exams, result...)
}

func (l *examOverviewLoader) loadCorrectedBadgeAsync(concurrent *Concurrent, userID int) {
	defer concurrent.Wg.Done()
	result, err := l.userRepo.GetCollectedBadge(userID)
	if err != nil {
		*concurrent.Err = err
	}
	l.correctedBadges = append(l.correctedBadges, result...)
}
