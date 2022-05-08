package loaders

import (
	"database-camp/internal/models/entities/activity"
	"database-camp/internal/models/entities/content"
	"database-camp/internal/repositories"
	"sync"
)

type contentRoadmapLoader struct {
	learningRepo repositories.LearningRepository
	userRepo     repositories.UserRepository

	content             *content.Content
	contentActivity     []activity.Activity
	learningProgression []content.LearningProgression
}

func NewContentRoadmapLoader(learningRepo repositories.LearningRepository, userRepo repositories.UserRepository) *contentRoadmapLoader {
	return &contentRoadmapLoader{learningRepo: learningRepo, userRepo: userRepo}
}

func (l *contentRoadmapLoader) GetContent() *content.Content {
	return l.content
}

func (l *contentRoadmapLoader) GetContentActivity() activity.Activities {
	return l.contentActivity
}

func (l *contentRoadmapLoader) GetLearningProgression() []content.LearningProgression {
	return l.learningProgression
}

func (l *contentRoadmapLoader) Load(userID int, contentID int) error {
	var wg sync.WaitGroup
	var err error
	concurrent := Concurrent{Wg: &wg, Err: &err}
	wg.Add(3)
	go l.loadLearningProgressionAsync(&concurrent, userID)
	go l.loadContentActivityAsync(&concurrent, contentID)
	go l.loadContentAsync(&concurrent, contentID)
	wg.Wait()
	return err
}

func (l *contentRoadmapLoader) loadContentAsync(concurrent *Concurrent, contentID int) {
	defer concurrent.Wg.Done()
	result, err := l.learningRepo.GetContent(contentID)
	if err != nil {
		*concurrent.Err = err
	}
	l.content = result
}

func (l *contentRoadmapLoader) loadContentActivityAsync(concurrent *Concurrent, contentID int) {
	defer concurrent.Wg.Done()
	result, err := l.learningRepo.GetContentActivity(contentID)
	if err != nil {
		*concurrent.Err = err
	}
	l.contentActivity = append(l.contentActivity, result...)
}

func (l *contentRoadmapLoader) loadLearningProgressionAsync(concurrent *Concurrent, userID int) {
	defer concurrent.Wg.Done()
	result, err := l.userRepo.GetLearningProgression(userID)
	if err != nil {
		*concurrent.Err = err
	}
	l.learningProgression = append(l.learningProgression, result...)
}
