package loaders

import (
	"DatabaseCamp/models/general"
	"DatabaseCamp/models/storages"
	"DatabaseCamp/repositories"
	"sync"
)

type learningOverviewLoader struct {
	learningRepo          repositories.ILearningRepository
	userRepo              repositories.IUserRepository
	overviewDB            []storages.OverviewDB
	learningProgressionDB []storages.LearningProgressionDB
}

func NewLearningOverviewLoader(learningRepo repositories.ILearningRepository, userRepo repositories.IUserRepository) *learningOverviewLoader {
	return &learningOverviewLoader{learningRepo: learningRepo, userRepo: userRepo}
}

func (l *learningOverviewLoader) GetOverviewDB() []storages.OverviewDB {
	return l.overviewDB
}

func (l *learningOverviewLoader) GetLearningProgressionDB() []storages.LearningProgressionDB {
	return l.learningProgressionDB
}

func (l *learningOverviewLoader) Load(userID int) error {
	var wg sync.WaitGroup
	var err error
	concurrent := general.Concurrent{Wg: &wg, Err: &err}
	wg.Add(2)
	go l.loadOverviewAsync(&concurrent)
	go l.loadLearningProgressionAsync(&concurrent, userID)
	wg.Wait()
	return err
}

func (l *learningOverviewLoader) loadOverviewAsync(concurrent *general.Concurrent) {
	defer concurrent.Wg.Done()
	result, err := l.learningRepo.GetOverview()
	if err != nil {
		*concurrent.Err = err
	}
	l.overviewDB = append(l.overviewDB, result...)
}

func (l *learningOverviewLoader) loadLearningProgressionAsync(concurrent *general.Concurrent, id int) {
	defer concurrent.Wg.Done()
	result, err := l.userRepo.GetLearningProgression(id)
	if err != nil {
		*concurrent.Err = err
	}
	l.learningProgressionDB = append(l.learningProgressionDB, result...)
}
