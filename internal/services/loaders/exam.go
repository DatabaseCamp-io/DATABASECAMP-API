package loaders

import (
	"database-camp/internal/models/entities/exam"
	"database-camp/internal/models/entities/user"
	"database-camp/internal/repositories"
	"sync"
)

type examLoader struct {
	examRepo repositories.ExamRepository
	userRepo repositories.UserRepository

	correctedBadge []user.CorrectedBadge
	exam           *exam.Exam
	examActivities []exam.ExamActivity
}

func NewExamLoader(examRepo repositories.ExamRepository, userRepo repositories.UserRepository) *examLoader {
	return &examLoader{examRepo: examRepo, userRepo: userRepo}
}

func (l *examLoader) GetCorrectedBadge() user.CorrectedBadges {
	return l.correctedBadge
}

func (l *examLoader) GetExamActivities() exam.ExamActivities {
	return l.examActivities
}

func (l *examLoader) GetExam() *exam.Exam {
	return l.exam
}

func (l *examLoader) Load(userID int, examID int) error {
	var wg sync.WaitGroup
	var err error
	concurrent := Concurrent{Wg: &wg, Err: &err}
	wg.Add(3)
	go l.loadCorrectedBadgeAsync(&concurrent, userID)
	go l.loadExamActivityAsync(&concurrent, examID)
	go l.loadExam(&concurrent, examID)
	wg.Wait()
	return err
}

func (l *examLoader) loadCorrectedBadgeAsync(concurrent *Concurrent, userID int) {
	defer concurrent.Wg.Done()
	result, err := l.userRepo.GetCollectedBadge(userID)
	if err != nil {
		*concurrent.Err = err
	}
	l.correctedBadge = append(l.correctedBadge, result...)
}

func (l *examLoader) loadExamActivityAsync(concurrent *Concurrent, examID int) {
	defer concurrent.Wg.Done()
	result, err := l.examRepo.GetExamActivities(examID)
	if err != nil {
		*concurrent.Err = err
	}
	l.examActivities = append(l.examActivities, result...)
}

func (l *examLoader) loadExam(concurrent *Concurrent, examID int) {
	defer concurrent.Wg.Done()
	result, err := l.examRepo.GetExam(examID)
	if err != nil {
		*concurrent.Err = err
	}
	l.exam = result
}
