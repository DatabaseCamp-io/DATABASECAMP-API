package loader

import (
	"DatabaseCamp/models"
	"DatabaseCamp/repository"
	"sync"
)

type activityLoader struct {
	learningRepo    repository.ILearningRepository
	userRepo        repository.IUserRepository
	ActivityDB      *models.ActivityDB
	ActivityHintsDB []models.HintDB
	UserHintsDB     []models.UserHintDB
}

func NewActivityLoader(learningRepo repository.ILearningRepository, userRepo repository.IUserRepository) *activityLoader {
	return &activityLoader{learningRepo: learningRepo, userRepo: userRepo}
}

func (l *activityLoader) Load(userID int, activityID int) error {
	var wg sync.WaitGroup
	var err error
	concurrent := models.Concurrent{Wg: &wg, Err: &err}
	wg.Add(3)
	go l.loadActivityAsync(&concurrent, activityID)
	go l.loadActivityHints(&concurrent, activityID)
	go l.loadUserHintsAsync(&concurrent, userID, activityID)
	wg.Wait()
	return err
}

func (l *activityLoader) loadActivityAsync(concurrent *models.Concurrent, activityID int) {
	defer concurrent.Wg.Done()
	var err error
	l.ActivityDB, err = l.learningRepo.GetActivity(activityID)
	if err != nil {
		*concurrent.Err = err
	}
}

func (l *activityLoader) loadUserHintsAsync(concurrent *models.Concurrent, userID int, activityID int) {
	defer concurrent.Wg.Done()
	result, e := l.userRepo.GetUserHint(userID, activityID)
	if e != nil {
		*concurrent.Err = e
	}
	l.UserHintsDB = append(l.UserHintsDB, result...)
}

func (l *activityLoader) loadActivityHints(concurrent *models.Concurrent, activityID int) {
	defer concurrent.Wg.Done()
	result, e := l.learningRepo.GetActivityHints(activityID)
	if e != nil {
		*concurrent.Err = e
	}
	l.ActivityHintsDB = append(l.ActivityHintsDB, result...)
}
