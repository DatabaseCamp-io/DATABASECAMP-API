package loaders

import (
	"DatabaseCamp/models/general"
	"DatabaseCamp/models/storages"
	"DatabaseCamp/repositories"
	"sync"
)

type contentRoadmapLoader struct {
	learningRepo          repositories.ILearningRepository
	userRepo              repositories.IUserRepository
	contentDB             *storages.ContentDB
	contentActivityDB     []storages.ActivityDB
	learningProgressionDB []storages.LearningProgressionDB
}

func NewContentRoadmapLoader(learningRepo repositories.ILearningRepository, userRepo repositories.IUserRepository) *contentRoadmapLoader {
	return &contentRoadmapLoader{learningRepo: learningRepo, userRepo: userRepo}
}

func (l *contentRoadmapLoader) GetContentDB() *storages.ContentDB {
	return l.contentDB
}

func (l *contentRoadmapLoader) GetContentActivityDB() []storages.ActivityDB {
	return l.contentActivityDB
}

func (l *contentRoadmapLoader) GetLearningProgressionDB() []storages.LearningProgressionDB {
	return l.learningProgressionDB
}

func (l *contentRoadmapLoader) Load(userID int, contentID int) error {
	var wg sync.WaitGroup
	var err error
	concurrent := general.Concurrent{Wg: &wg, Err: &err}
	wg.Add(3)
	go l.loadLearningProgressionAsync(&concurrent, userID)
	go l.loadContentActivityAsync(&concurrent, contentID)
	go l.loadContentAsync(&concurrent, contentID)
	wg.Wait()
	return err
}

func (l *contentRoadmapLoader) loadContentAsync(concurrent *general.Concurrent, contentID int) {
	defer concurrent.Wg.Done()
	result, err := l.learningRepo.GetContent(contentID)
	if err != nil {
		*concurrent.Err = err
	}
	l.contentDB = result
}

func (l *contentRoadmapLoader) loadLearningProgressionAsync(concurrent *general.Concurrent, userID int) {
	defer concurrent.Wg.Done()
	result, err := l.userRepo.GetLearningProgression(userID)
	if err != nil {
		*concurrent.Err = err
	}
	l.learningProgressionDB = append(l.learningProgressionDB, result...)
}

func (l *contentRoadmapLoader) loadContentActivityAsync(concurrent *general.Concurrent, contentID int) {
	defer concurrent.Wg.Done()
	result, err := l.learningRepo.GetContentActivity(contentID)
	if err != nil {
		*concurrent.Err = err
	}
	l.contentActivityDB = append(l.contentActivityDB, result...)
}
