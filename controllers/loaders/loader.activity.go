package loaders

import (
	"DatabaseCamp/models/general"
	"DatabaseCamp/repositories"
	"sync"
)

type activityLoader struct {
	learningRepo    repositories.ILearningRepository
	userRepo        repositories.IUserRepository
	ActivityDB      *general.ActivityDB
	ActivityHintsDB []general.HintDB
	UserHintsDB     []general.UserHintDB
}

func NewActivityLoader(learningRepo repositories.ILearningRepository, userRepo repositories.IUserRepository) *activityLoader {
	return &activityLoader{learningRepo: learningRepo, userRepo: userRepo}
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
	l.ActivityDB, err = l.learningRepo.GetActivity(activityID)
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
	l.UserHintsDB = append(l.UserHintsDB, result...)
}

func (l *activityLoader) loadActivityHints(concurrent *general.Concurrent, activityID int) {
	defer concurrent.Wg.Done()
	result, e := l.learningRepo.GetActivityHints(activityID)
	if e != nil {
		*concurrent.Err = e
	}
	l.ActivityHintsDB = append(l.ActivityHintsDB, result...)
}
