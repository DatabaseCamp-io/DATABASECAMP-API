package loaders

import (
	"database-camp/internal/models/entities/activity"
	"database-camp/internal/models/entities/content"
	"database-camp/internal/repositories"
	"sync"
)

type checkAnswerLoader struct {
	learningRepo repositories.LearningRepository

	choices     activity.Choices
	activity    *activity.Activity
	progression *content.LearningProgression
}

func NewCheckAnswerLoader(learningRepo repositories.LearningRepository) *checkAnswerLoader {
	return &checkAnswerLoader{learningRepo: learningRepo}
}

func (c *checkAnswerLoader) GetChoices() activity.Choices {
	return c.choices
}

func (c *checkAnswerLoader) GetActivity() *activity.Activity {
	return c.activity
}

func (c *checkAnswerLoader) GetProgression() *content.LearningProgression {
	return c.progression
}

func (c *checkAnswerLoader) Load(activityID int, activityTypeID int) error {
	var wg sync.WaitGroup
	var err error
	concurrent := Concurrent{Wg: &wg, Err: &err}
	wg.Add(3)
	go c.loadActivityAsync(&concurrent, activityID)
	go c.loadChioces(&concurrent, activityID, activityTypeID)
	go c.loadProgression(&concurrent, activityID)
	wg.Wait()
	return err
}

func (c *checkAnswerLoader) loadActivityAsync(concurrent *Concurrent, activityID int) {
	defer concurrent.Wg.Done()
	var err error
	c.activity, err = c.learningRepo.GetActivity(activityID)
	if err != nil {
		*concurrent.Err = err
	}
}

func (c *checkAnswerLoader) loadChioces(concurrent *Concurrent, activityID int, activityTypeID int) {
	defer concurrent.Wg.Done()
	var err error
	c.choices, err = c.learningRepo.GetActivityChoices(activityID, activityTypeID)
	if err != nil {
		*concurrent.Err = err
	}
}

func (c *checkAnswerLoader) loadProgression(concurrent *Concurrent, activityID int) {
	defer concurrent.Wg.Done()
	var err error
	c.progression, err = c.learningRepo.GetCorrectProgression(activityID)
	if err != nil {
		*concurrent.Err = err
	}
}
