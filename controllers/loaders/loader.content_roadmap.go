package loaders

import (
	"DatabaseCamp/models/general"
	"DatabaseCamp/repositories"
	"sync"
)

type contentRoadmapLoader struct {
	learningRepo          repositories.ILearningRepository
	userRepo              repositories.IUserRepository
	ContentDB             *general.ContentDB
	ContentActivityDB     []general.ActivityDB
	LearningProgressionDB []general.LearningProgressionDB
}

func NewContentRoadmapLoader(learningRepo repositories.ILearningRepository, userRepo repositories.IUserRepository) *contentRoadmapLoader {
	return &contentRoadmapLoader{learningRepo: learningRepo, userRepo: userRepo}
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
	l.ContentDB = result
}

func (l *contentRoadmapLoader) loadLearningProgressionAsync(concurrent *general.Concurrent, userID int) {
	defer concurrent.Wg.Done()
	result, err := l.userRepo.GetLearningProgression(userID)
	if err != nil {
		*concurrent.Err = err
	}
	l.LearningProgressionDB = append(l.LearningProgressionDB, result...)
}

func (l *contentRoadmapLoader) loadContentActivityAsync(concurrent *general.Concurrent, contentID int) {
	defer concurrent.Wg.Done()
	result, err := l.learningRepo.GetContentActivity(contentID)
	if err != nil {
		*concurrent.Err = err
	}
	l.ContentActivityDB = append(l.ContentActivityDB, result...)
}
