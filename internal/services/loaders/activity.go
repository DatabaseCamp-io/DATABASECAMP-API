package loaders

import (
	"database-camp/internal/models/entities/activity"
	"database-camp/internal/repositories"
	"sync"
)

type activityLoader struct {
	learningRepo repositories.LearningRepository
	userRepo     repositories.UserRepository

	activity      *activity.Activity
	activityHints []activity.Hint
	userHints     []activity.UserHint
}

func NewActivityLoader(learningRepo repositories.LearningRepository, userRepo repositories.UserRepository) *activityLoader {
	return &activityLoader{learningRepo: learningRepo, userRepo: userRepo}
}

func (l *activityLoader) GetActivity() *activity.Activity {
	return l.activity
}

func (l *activityLoader) GetActivityHints() activity.Hints {
	return l.activityHints
}

func (l *activityLoader) GetUserHints() activity.UserHints {
	return l.userHints
}

func (l *activityLoader) Load(userID int, activityID int) error {
	var wg sync.WaitGroup
	var err error
	concurrent := Concurrent{Wg: &wg, Err: &err}
	wg.Add(3)
	go l.loadActivityAsync(&concurrent, activityID)
	go l.loadActivityHints(&concurrent, activityID)
	go l.loadUserHintsAsync(&concurrent, userID, activityID)
	wg.Wait()
	return err
}

func (l *activityLoader) loadActivityAsync(concurrent *Concurrent, activityID int) {
	defer concurrent.Wg.Done()
	var err error
	l.activity, err = l.learningRepo.GetActivity(activityID)
	if err != nil {
		*concurrent.Err = err
	}
}

func (l *activityLoader) loadUserHintsAsync(concurrent *Concurrent, userID int, activityID int) {
	defer concurrent.Wg.Done()
	result, e := l.userRepo.GetUserHint(userID, activityID)
	if e != nil {
		*concurrent.Err = e
	}
	l.userHints = append(l.userHints, result...)
}

func (l *activityLoader) loadActivityHints(concurrent *Concurrent, activityID int) {
	defer concurrent.Wg.Done()
	result, e := l.learningRepo.GetActivityHints(activityID)
	if e != nil {
		*concurrent.Err = e
	}
	l.activityHints = append(l.activityHints, result...)
}
