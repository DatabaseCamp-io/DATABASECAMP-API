package loader

import (
	"DatabaseCamp/models"
	"DatabaseCamp/repositories"
	"sync"
)

type learningOverviewLoader struct {
	learningRepo          repositories.ILearningRepository
	userRepo              repositories.IUserRepository
	OverviewDB            []models.OverviewDB
	LearningProgressionDB []models.LearningProgressionDB
}

func NewLearningOverviewLoader(learningRepo repositories.ILearningRepository, userRepo repositories.IUserRepository) *learningOverviewLoader {
	return &learningOverviewLoader{learningRepo: learningRepo, userRepo: userRepo}
}

func (l *learningOverviewLoader) Load(userID int) error {
	var wg sync.WaitGroup
	var err error
	concurrent := models.Concurrent{Wg: &wg, Err: &err}
	wg.Add(2)
	go l.loadOverviewAsync(&concurrent)
	go l.loadLearningProgressionAsync(&concurrent, userID)
	wg.Wait()
	return err
}

func (l *learningOverviewLoader) loadOverviewAsync(concurrent *models.Concurrent) {
	defer concurrent.Wg.Done()
	result, err := l.learningRepo.GetOverview()
	if err != nil {
		*concurrent.Err = err
	}
	l.OverviewDB = append(l.OverviewDB, result...)
}

func (l *learningOverviewLoader) loadLearningProgressionAsync(concurrent *models.Concurrent, id int) {
	defer concurrent.Wg.Done()
	result, err := l.userRepo.GetLearningProgression(id)
	if err != nil {
		*concurrent.Err = err
	}
	l.LearningProgressionDB = append(l.LearningProgressionDB, result...)
}
