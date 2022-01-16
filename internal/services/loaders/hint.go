package loaders

import (
	"database-camp/internal/models/entities/activity"
	"database-camp/internal/models/entities/user"
	"database-camp/internal/repositories"
	"sync"
)

type hintLoader struct {
	learningRepo repositories.LearningRepository
	userRepo     repositories.UserRepository

	activityHints []activity.Hint
	userHints     []activity.UserHint
	user          *user.User
}

func NewHintLoader(learningRepo repositories.LearningRepository, userRepo repositories.UserRepository) *hintLoader {
	return &hintLoader{learningRepo: learningRepo, userRepo: userRepo}
}

func (l *hintLoader) GetActivityHintsDB() activity.Hints {
	return l.activityHints
}

func (l *hintLoader) GetUserHintsDB() activity.UserHints {
	return l.userHints
}

func (l *hintLoader) GetUser() *user.User {
	return l.user
}

func (l *hintLoader) Load(userID int, activityID int) error {
	var wg sync.WaitGroup
	var err error
	concurrent := Concurrent{Wg: &wg, Err: &err}
	wg.Add(3)
	go l.loadActivityHints(&concurrent, activityID)
	go l.loadUserHintsAsync(&concurrent, userID, activityID)
	go l.loadUser(&concurrent, userID)
	wg.Wait()
	return err
}

func (l *hintLoader) loadUser(concurrent *Concurrent, userID int) {
	defer concurrent.Wg.Done()
	var err error
	l.user, err = l.userRepo.GetUserByID(userID)
	if err != nil {
		*concurrent.Err = err
	}
}

func (l *hintLoader) loadUserHintsAsync(concurrent *Concurrent, userID int, activityID int) {
	defer concurrent.Wg.Done()
	result, e := l.userRepo.GetUserHint(userID, activityID)
	if e != nil {
		*concurrent.Err = e
	}
	l.userHints = append(l.userHints, result...)
}

func (l *hintLoader) loadActivityHints(concurrent *Concurrent, activityID int) {
	defer concurrent.Wg.Done()
	result, e := l.learningRepo.GetActivityHints(activityID)
	if e != nil {
		*concurrent.Err = e
	}
	l.activityHints = append(l.activityHints, result...)
}
