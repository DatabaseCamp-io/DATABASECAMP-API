package loaders

// loader.content_roadmap.go
/**
 * 	This file is a part of controller, used to load concurrency roadmap of the learning content
 */

import (
	"DatabaseCamp/models/general"
	"DatabaseCamp/models/storages"
	"DatabaseCamp/repositories"
	"sync"
)

/**
 * This class load concurrency roadmap of the learning content
 */
type contentRoadmapLoader struct {
	learningRepo repositories.ILearningRepository // repository for load learning content data
	userRepo     repositories.IUserRepository     // repository for load user learning progression data

	contentDB             *storages.ContentDB              // learning content data from the database
	contentActivityDB     []storages.ActivityDB            // activities of the content from the database
	learningProgressionDB []storages.LearningProgressionDB // learning progression of the user from the database
}

/**
 * Constructor creates a new contentRoadmapLoader instance
 *
 * @param   learningRepo    Learning Repository for load learning data
 * @param   userRepo        User Repository for load user data
 *
 * @return 	instance of contentRoadmapLoader
 */
func NewContentRoadmapLoader(learningRepo repositories.ILearningRepository, userRepo repositories.IUserRepository) *contentRoadmapLoader {
	return &contentRoadmapLoader{learningRepo: learningRepo, userRepo: userRepo}
}

/**
 * Getter for getting contentDB
 *
 * @return contentDB
 */
func (l *contentRoadmapLoader) GetContentDB() *storages.ContentDB {
	return l.contentDB
}

/**
 * Getter for getting contentActivityDB
 *
 * @return contentActivityDB
 */
func (l *contentRoadmapLoader) GetContentActivityDB() []storages.ActivityDB {
	return l.contentActivityDB
}

/**
 * Getter for getting learningProgressionDB
 *
 * @return learningProgressionDB
 */
func (l *contentRoadmapLoader) GetLearningProgressionDB() []storages.LearningProgressionDB {
	return l.learningProgressionDB
}

/**
 * Load concurrency roadmap of the learning content
 *
 * @param   userID     		User ID for getting user hints information of the activity
 * @param   contentID   	Content ID for getting learning content data
 *
 * @return the error of loading data
 */
func (l *contentRoadmapLoader) Load(userID int, contentID int) error {
	var wg sync.WaitGroup
	var err error
	concurrent := general.Concurrent{Wg: &wg, Err: &err}
	wg.Add(3)
	go l.loadLearningProgressionAsync(&concurrent, userID)
	go l.loadContentActivityAsync(&concurrent, contentID)
	go l.loadContentAsync(&concurrent, contentID)
	wg.Wait()
	return err
}

/**
 * Load learning content data from the database
 *
 * @param   concurrent     	Concurrent model for doing load concurrency
 * @param   contentID    	Content ID for getting learning content data
 */
func (l *contentRoadmapLoader) loadContentAsync(concurrent *general.Concurrent, contentID int) {
	defer concurrent.Wg.Done()
	result, err := l.learningRepo.GetContent(contentID)
	if err != nil {
		*concurrent.Err = err
	}
	l.contentDB = result
}

/**
 * Load activities of the content from the database
 *
 * @param   concurrent     	Concurrent model for doing load concurrency
 * @param   contentID    	Content ID for getting learning content data
 */
func (l *contentRoadmapLoader) loadContentActivityAsync(concurrent *general.Concurrent, contentID int) {
	defer concurrent.Wg.Done()
	result, err := l.learningRepo.GetContentActivity(contentID)
	if err != nil {
		*concurrent.Err = err
	}
	l.contentActivityDB = append(l.contentActivityDB, result...)
}

/**
 * Load learning progression of the user from the database
 *
 * @param   concurrent     	Concurrent model for doing load concurrency
 * @param   userID    		User ID for getting learning progression of the user
 */
func (l *contentRoadmapLoader) loadLearningProgressionAsync(concurrent *general.Concurrent, userID int) {
	defer concurrent.Wg.Done()
	result, err := l.userRepo.GetLearningProgression(userID)
	if err != nil {
		*concurrent.Err = err
	}
	l.learningProgressionDB = append(l.learningProgressionDB, result...)
}
