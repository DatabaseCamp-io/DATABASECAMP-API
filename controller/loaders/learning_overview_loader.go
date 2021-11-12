package loader

import (
	"DatabaseCamp/models"
	"DatabaseCamp/repository"
	"sync"
)

type learningOverviewLoader struct {
	learningRepo          repository.ILearningRepository
	userRepo              repository.IUserRepository
	OverviewDB            []models.OverviewDB
	LearningProgressionDB []models.LearningProgressionDB
	ActivityDB            []models.ActivityDB
}

func NewLearningOverviewLoader(learningRepo repository.ILearningRepository, userRepo repository.IUserRepository) *learningOverviewLoader {
	return &learningOverviewLoader{learningRepo: learningRepo, userRepo: userRepo}
}

func (l *learningOverviewLoader) Load(userID int) error {
	var wg sync.WaitGroup
	var err error
	concurrent := models.Concurrent{Wg: &wg, Err: &err}
	wg.Add(2)
	go l.loadOverviewAsync(&concurrent, &l.OverviewDB)
	go l.loadLearningProgressionAsync(&concurrent, userID, &l.LearningProgressionDB)
	wg.Wait()
	return err
}

func (l *learningOverviewLoader) loadOverviewAsync(concurrent *models.Concurrent, overview *[]models.OverviewDB) {
	defer concurrent.Wg.Done()
	result, err := l.learningRepo.GetOverview()
	if err != nil {
		*concurrent.Err = err
	}
	*overview = append(*overview, result...)
}

func (l *learningOverviewLoader) loadLearningProgressionAsync(concurrent *models.Concurrent, id int, learningProgression *[]models.LearningProgressionDB) {
	defer concurrent.Wg.Done()
	result, err := l.userRepo.GetLearningProgression(id)
	if err != nil {
		*concurrent.Err = err
	}
	*learningProgression = append(*learningProgression, result...)
}
