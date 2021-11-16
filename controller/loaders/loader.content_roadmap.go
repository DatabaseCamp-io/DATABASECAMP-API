package loader

import (
	"DatabaseCamp/models"
	"DatabaseCamp/repository"
	"sync"
)

type contentRoadmapLoader struct {
	learningRepo          repository.ILearningRepository
	userRepo              repository.IUserRepository
	ContentDB             *models.ContentDB
	ContentActivityDB     []models.ActivityDB
	LearningProgressionDB []models.LearningProgressionDB
}

func NewContentRoadmapLoader(learningRepo repository.ILearningRepository, userRepo repository.IUserRepository) *contentRoadmapLoader {
	return &contentRoadmapLoader{learningRepo: learningRepo, userRepo: userRepo}
}

func (l *contentRoadmapLoader) Load(userID int, contentID int) error {
	var wg sync.WaitGroup
	var err error
	concurrent := models.Concurrent{Wg: &wg, Err: &err}
	wg.Add(3)
	go l.loadLearningProgressionAsync(&concurrent, userID)
	go l.loadContentActivityAsync(&concurrent, contentID)
	go l.loadContentAsync(&concurrent, contentID)
	wg.Wait()
	return err
}

func (l *contentRoadmapLoader) loadContentAsync(concurrent *models.Concurrent, contentID int) {
	defer concurrent.Wg.Done()
	result, err := l.learningRepo.GetContent(contentID)
	if err != nil {
		*concurrent.Err = err
	}
	l.ContentDB = result
}

func (l *contentRoadmapLoader) loadLearningProgressionAsync(concurrent *models.Concurrent, userID int) {
	defer concurrent.Wg.Done()
	result, err := l.userRepo.GetLearningProgression(userID)
	if err != nil {
		*concurrent.Err = err
	}
	l.LearningProgressionDB = append(l.LearningProgressionDB, result...)
}

func (l *contentRoadmapLoader) loadContentActivityAsync(concurrent *models.Concurrent, contentID int) {
	defer concurrent.Wg.Done()
	result, err := l.learningRepo.GetContentActivity(contentID)
	if err != nil {
		*concurrent.Err = err
	}
	l.ContentActivityDB = append(l.ContentActivityDB, result...)
}
