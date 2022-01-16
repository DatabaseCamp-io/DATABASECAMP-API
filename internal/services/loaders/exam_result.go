package loaders

import (
	"database-camp/internal/models/entities/exam"
	"database-camp/internal/repositories"
	"sync"
)

type examResultLoader struct {
	examRepo repositories.ExamRepository

	resultActivities []exam.ResultActivity
	examResult       *exam.ExamResult
}

func NewExamResultLoader(examRepo repositories.ExamRepository) *examResultLoader {
	return &examResultLoader{examRepo: examRepo}
}

func (l examResultLoader) GetExamResult() *exam.ExamResult {
	return l.examResult
}

func (l examResultLoader) GetResultActivities() exam.ResultActivities {
	return l.resultActivities
}

func (l *examResultLoader) Load(userID int, examResultID int) error {
	var wg sync.WaitGroup
	var err error
	concurrent := Concurrent{Wg: &wg, Err: &err}
	wg.Add(2)
	go l.loadExamResult(&concurrent, userID, examResultID)
	go l.loadResultActivities(&concurrent, examResultID)
	wg.Wait()
	return err
}

func (l *examResultLoader) loadExamResult(concurrent *Concurrent, userID int, examResultID int) {
	defer concurrent.Wg.Done()

	result, err := l.examRepo.GetExamResult(userID, examResultID)
	if err != nil {
		*concurrent.Err = err
	}

	l.examResult = result
}

func (l *examResultLoader) loadResultActivities(concurrent *Concurrent, examResultID int) {
	defer concurrent.Wg.Done()

	result, err := l.examRepo.GetActivitiesResult(examResultID)
	if err != nil {
		*concurrent.Err = err
	}

	l.resultActivities = append(l.resultActivities, result...)
}
