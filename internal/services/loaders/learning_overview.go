package loaders

import (
	"database-camp/internal/models/entities/content"
	"database-camp/internal/repositories"
	"sync"
)

type learningOverviewLoader struct {
	learningRepo repositories.LearningRepository
	userRepo     repositories.UserRepository

	overview            []content.Overview
	learningProgression []content.LearningProgression
	preExamID           *int
}

func NewLearningOverviewLoader(learningRepo repositories.LearningRepository, userRepo repositories.UserRepository) *learningOverviewLoader {
	return &learningOverviewLoader{learningRepo: learningRepo, userRepo: userRepo}
}

func (l *learningOverviewLoader) GetOverview() content.OverviewList {
	return l.overview
}

func (l *learningOverviewLoader) GetLearningProgression() content.LearningProgressionList {
	return l.learningProgression
}

func (l *learningOverviewLoader) GetPreExamID() *int {
	return l.preExamID
}

func (l *learningOverviewLoader) Load(userID int) error {
	var wg sync.WaitGroup
	var err error
	concurrent := Concurrent{Wg: &wg, Err: &err}
	wg.Add(3)
	go l.loadOverviewAsync(&concurrent)
	go l.loadPreExamIDAsync(&concurrent, userID)
	go l.loadLearningProgressionAsync(&concurrent, userID)
	wg.Wait()
	return err
}

func (l *learningOverviewLoader) loadOverviewAsync(concurrent *Concurrent) {
	defer concurrent.Wg.Done()
	result, err := l.learningRepo.GetOverview()
	if err != nil {
		*concurrent.Err = err
	}
	l.overview = append(l.overview, result...)
}

func (l *learningOverviewLoader) loadLearningProgressionAsync(concurrent *Concurrent, id int) {
	defer concurrent.Wg.Done()
	result, err := l.userRepo.GetLearningProgression(id)
	if err != nil {
		*concurrent.Err = err
	}
	l.learningProgression = append(l.learningProgression, result...)
}

func (l *learningOverviewLoader) loadPreExamIDAsync(concurrent *Concurrent, id int) {
	defer concurrent.Wg.Done()
	result, err := l.userRepo.GetPreExamID(id)
	if err != nil {
		*concurrent.Err = err
	}

	l.preExamID = result
}
