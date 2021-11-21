package loaders

import (
	"DatabaseCamp/models/general"
	"DatabaseCamp/models/storages"
	"DatabaseCamp/repositories"
	"sync"
)

type hintLoader struct {
	learningRepo    repositories.ILearningRepository
	userRepo        repositories.IUserRepository
	activityHintsDB []storages.HintDB
	userHintsDB     []storages.UserHintDB
	userDB          *storages.UserDB
}

func NewHintLoader(learningRepo repositories.ILearningRepository, userRepo repositories.IUserRepository) *hintLoader {
	return &hintLoader{learningRepo: learningRepo, userRepo: userRepo}
}

func (l *hintLoader) GetActivityHintsDB() []storages.HintDB {
	return l.activityHintsDB
}

func (l *hintLoader) GetUserHintsDB() []storages.UserHintDB {
	return l.userHintsDB
}

func (l *hintLoader) GetUserDB() *storages.UserDB {
	return l.userDB
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
	l.userDB, err = l.userRepo.GetUserByID(userID)
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
	l.userHintsDB = append(l.userHintsDB, result...)
}

func (l *hintLoader) loadActivityHints(concurrent *general.Concurrent, activityID int) {
	defer concurrent.Wg.Done()
	result, e := l.learningRepo.GetActivityHints(activityID)
	if e != nil {
		*concurrent.Err = e
	}
	l.activityHintsDB = append(l.activityHintsDB, result...)
}
