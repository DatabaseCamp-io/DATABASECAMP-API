package loaders

import (
	"DatabaseCamp/models/general"
	"DatabaseCamp/repositories"
	"sync"
)

type hintLoader struct {
	learningRepo    repositories.ILearningRepository
	userRepo        repositories.IUserRepository
	ActivityHintsDB []general.HintDB
	UserHintsDB     []general.UserHintDB
	UserDB          *general.UserDB
}

func NewHintLoader(learningRepo repositories.ILearningRepository, userRepo repositories.IUserRepository) *hintLoader {
	return &hintLoader{learningRepo: learningRepo, userRepo: userRepo}
}

func (l *hintLoader) Load(userID int, activityID int) error {
	var wg sync.WaitGroup
	var err error
	concurrent := general.Concurrent{Wg: &wg, Err: &err}
	wg.Add(3)
	go l.loadActivityHints(&concurrent, activityID)
	go l.loadUserHintsAsync(&concurrent, userID, activityID)
	go l.loadUser(&concurrent, userID)
	wg.Wait()
	return err
}

func (l *hintLoader) loadUser(concurrent *general.Concurrent, userID int) {
	defer concurrent.Wg.Done()
	var err error
	l.UserDB, err = l.userRepo.GetUserByID(userID)
	if err != nil {
		*concurrent.Err = err
	}
}

func (l *hintLoader) loadUserHintsAsync(concurrent *general.Concurrent, userID int, activityID int) {
	defer concurrent.Wg.Done()
	result, e := l.userRepo.GetUserHint(userID, activityID)
	if e != nil {
		*concurrent.Err = e
	}
	l.UserHintsDB = append(l.UserHintsDB, result...)
}

func (l *hintLoader) loadActivityHints(concurrent *general.Concurrent, activityID int) {
	defer concurrent.Wg.Done()
	result, e := l.learningRepo.GetActivityHints(activityID)
	if e != nil {
		*concurrent.Err = e
	}
	l.ActivityHintsDB = append(l.ActivityHintsDB, result...)
}
