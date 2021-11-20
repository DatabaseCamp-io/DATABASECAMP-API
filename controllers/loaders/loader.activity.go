package loaders

import (
	"DatabaseCamp/models/general"
	"DatabaseCamp/models/storages"
	"DatabaseCamp/repositories"
	"sync"
)

type activityLoader struct {
	learningRepo    repositories.ILearningRepository
	userRepo        repositories.IUserRepository
	activityDB      *storages.ActivityDB
	activityHintsDB []storages.HintDB
	userHintsDB     []storages.UserHintDB
}

func NewActivityLoader(learningRepo repositories.ILearningRepository, userRepo repositories.IUserRepository) *activityLoader {
	return &activityLoader{learningRepo: learningRepo, userRepo: userRepo}
}

func (l *activityLoader) GetActivityDB() *storages.ActivityDB {
	return l.activityDB
}

func (l *activityLoader) GetActivityHintsDB() []storages.HintDB {
	return l.activityHintsDB
}

func (l *activityLoader) GetUserHintsDB() []storages.UserHintDB {
	return l.userHintsDB
}

func (l *activityLoader) Load(userID int, activityID int) error {
	var wg sync.WaitGroup
	var err error
	concurrent := general.Concurrent{Wg: &wg, Err: &err}
	wg.Add(3)
	go l.loadActivityAsync(&concurrent, activityID)
	go l.loadActivityHints(&concurrent, activityID)
	go l.loadUserHintsAsync(&concurrent, userID, activityID)
	wg.Wait()
	return err
}

func (l *activityLoader) loadActivityAsync(concurrent *general.Concurrent, activityID int) {
	defer concurrent.Wg.Done()
	var err error
	l.activityDB, err = l.learningRepo.GetActivity(activityID)
	if err != nil {
		*concurrent.Err = err
	}
}

func (l *activityLoader) loadUserHintsAsync(concurrent *general.Concurrent, userID int, activityID int) {
	defer concurrent.Wg.Done()
	result, e := l.userRepo.GetUserHint(userID, activityID)
	if e != nil {
		*concurrent.Err = e
	}
	l.userHintsDB = append(l.userHintsDB, result...)
}

func (l *activityLoader) loadActivityHints(concurrent *general.Concurrent, activityID int) {
	defer concurrent.Wg.Done()
	result, e := l.learningRepo.GetActivityHints(activityID)
	if e != nil {
		*concurrent.Err = e
	}
	l.activityHintsDB = append(l.activityHintsDB, result...)
}
