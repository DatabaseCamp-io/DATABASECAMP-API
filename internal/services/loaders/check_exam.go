package loaders

import (
	"database-camp/internal/models/entities/exam"
	"database-camp/internal/repositories"
	"sync"
)

type checkExamLoader struct {
	examRepo repositories.ExamRepository

	exam           *exam.Exam
	examActivities []exam.ExamActivity
}

func NewCheckExamLoader(examRepo repositories.ExamRepository) *checkExamLoader {
	return &checkExamLoader{examRepo: examRepo}
}

func (l checkExamLoader) GetExam() *exam.Exam {
	return l.exam
}

func (l checkExamLoader) GetExamActivities() []exam.ExamActivity {
	return l.examActivities
}

func (l *checkExamLoader) Load(examID int) error {
	var wg sync.WaitGroup
	var err error
	concurrent := Concurrent{Wg: &wg, Err: &err}
	wg.Add(2)
	go l.loadExamActivies(&concurrent, examID)
	go l.loadExam(&concurrent, examID)
	wg.Wait()
	return err
}

func (l *checkExamLoader) loadExamActivies(concurrent *Concurrent, examID int) {
	defer concurrent.Wg.Done()

	result, err := l.examRepo.GetExamActivities(examID)
	if err != nil {
		*concurrent.Err = err
	}

	l.examActivities = append(l.examActivities, result...)
}

func (l *checkExamLoader) loadExam(concurrent *Concurrent, examID int) {
	defer concurrent.Wg.Done()
	result, err := l.examRepo.GetExam(examID)
	if err != nil {
		*concurrent.Err = err
	}
	l.exam = result
}
