package loaders

import (
	"database-camp/internal/models/entities/activity"
	"database-camp/internal/models/entities/content"
	"database-camp/internal/repositories"
	"sync"
)

type checkPeerReviewLoader struct {
	learningRepo repositories.LearningRepository

	erAnswer     *activity.ERAnswer
	erChoice     *activity.ERChoice
	progresstion *content.LearningProgression
}

func NewCheckPeerReviewLoader(learningRepo repositories.LearningRepository) *checkPeerReviewLoader {
	return &checkPeerReviewLoader{learningRepo: learningRepo}
}

func (l checkPeerReviewLoader) GetERAnswer() *activity.ERAnswer {
	return l.erAnswer
}

func (l checkPeerReviewLoader) GetERChoice() *activity.ERChoice {
	return l.erChoice
}

func (l checkPeerReviewLoader) GetProgression() *content.LearningProgression {
	return l.progresstion
}

func (l *checkPeerReviewLoader) Load(erAnswerID int, activityID int) error {
	var wg sync.WaitGroup
	var err error
	concurrent := Concurrent{Wg: &wg, Err: &err}
	wg.Add(3)
	go l.loadERAnswer(&concurrent, erAnswerID)
	go l.loadERChoice(&concurrent, activityID)
	go l.loadProgression(&concurrent, activityID)
	wg.Wait()
	return err
}

func (l *checkPeerReviewLoader) loadERAnswer(concurrent *Concurrent, erAnswerID int) {
	defer concurrent.Wg.Done()
	result, err := l.learningRepo.GetPeerChoice(&erAnswerID)
	if err != nil {
		*concurrent.Err = err
	}
	l.erAnswer = &result
}

func (l *checkPeerReviewLoader) loadERChoice(concurrent *Concurrent, activityID int) {
	defer concurrent.Wg.Done()
	result, err := l.learningRepo.GetERChoice(activityID)
	if err != nil {
		*concurrent.Err = err
	}
	l.erChoice = &result
}

func (l *checkPeerReviewLoader) loadProgression(concurrent *Concurrent, activityID int) {
	defer concurrent.Wg.Done()
	var err error
	l.progresstion, err = l.learningRepo.GetCorrectProgression(activityID)
	if err != nil {
		*concurrent.Err = err
	}
}
