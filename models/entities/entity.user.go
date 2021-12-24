package entities

// entity.user.go
/**
 * 	This file is a part of models, used to collect model for entities of user
 */

import (
	"DatabaseCamp/models/request"
	"DatabaseCamp/models/storages"
	"DatabaseCamp/utils"
	"os"
	"time"
)

// Type for changes user points mode
type ChangePointMode string

// Change mode for update user point
var Mode = struct {
	Add    ChangePointMode
	Reduce ChangePointMode
}{
	"+",
	"-",
}

// Model of badge for user entity
type Badge struct {
	ID          int    `json:"badge_id"`
	ImagePath   string `json:"icon_path"`
	Name        string `json:"name"`
	IsCollected bool   `json:"is_collect"`
}

/**
 * This class manage user model
 */
type User struct {
	id               int
	name             string
	email            string
	password         string
	badges           []Badge
	createdTimestamp time.Time
	updatedTimestamp time.Time
}

/**
 * Constructor creates a new user instance by user request
 *
 * @param   request      User request
 *
 * @return 	instance of user
 */
func NewUserByRequest(request request.UserRequest) User {
	user := User{
		name:     request.Name,
		email:    request.Email,
		password: request.Password,
	}
	user.SetTimestamp()
	user.HashPassword()
	return user
}

/**
 * Constructor creates a new user instance by user database model
 *
 * @param   request      User database model
 *
 * @return 	instance of user
 */
func NewUserByUserDB(userDB storages.UserDB) User {
	return User{
		id:               userDB.ID,
		name:             userDB.Name,
		email:            userDB.Email,
		password:         userDB.Password,
		createdTimestamp: userDB.CreatedTimestamp,
		updatedTimestamp: userDB.UpdatedTimestamp,
	}
}

/**
 * Getter for getting user ID
 *
 * @return user ID
 */
func (u *User) GetID() int {
	return u.id
}

/**
 * Getter for getting user name
 *
 * @return user name
 */
func (u *User) GetName() string {
	return u.name
}

/**
 * Getter for getting user email
 *
 * @return user email
 */
func (u *User) GetEmail() string {
	return u.email
}

/**
 * Getter for getting user badges
 *
 * @return user badges
 */
func (u *User) GetBadges() []Badge {
	return u.badges
}

/**
 * Getter for getting user created timestamp
 *
 * @return user created timestamp
 */
func (u *User) GetCreatedTimestamp() time.Time {
	return u.createdTimestamp
}

/**
 * Getter for getting user updated timestamp
 *
 * @return user updated timestamp
 */
func (u *User) GetUpdatedTimstamp() time.Time {
	return u.updatedTimestamp
}

/**
 * Setter user ID
 *
 * @param id 	User ID to set
 */
func (u *User) SetID(id int) {
	u.id = id
}

/**
 * Setter user ID
 *
 * @param 	allBadgesDB 			All badges in the application
 * @param 	correctedBadgesDB 		Corrected badges of the user
 */
func (u *User) SetCorrectedBadges(allBadgesDB []storages.BadgeDB, correctedBadgesDB []storages.UserBadgeDB) {
	for _, badgeDB := range allBadgesDB {
		u.badges = append(u.badges, Badge{
			ID:          badgeDB.ID,
			ImagePath:   badgeDB.ImagePath,
			Name:        badgeDB.Name,
			IsCollected: u.isCorrectedBadge(correctedBadgesDB, badgeDB.ID),
		})
	}
}

/**
 * Setter user timestamp to now
 */
func (u *User) SetTimestamp() {
	u.createdTimestamp = time.Now().Local()
	u.updatedTimestamp = time.Now().Local()
}

/**
 * Hash user password
 */
func (u *User) HashPassword() {
	u.password = utils.NewHelper().HashAndSalt(u.password)
}

/**
 * Check corrected badge
 *
 * @param 	allBadgesDB 	All badges in the application
 * @param 	badgeID 		Badge ID to check
 *
 * @return 	true if badge id is correct, false otherwise
 */
func (u *User) isCorrectedBadge(allBadgesDB []storages.UserBadgeDB, badgeID int) bool {
	for _, correctedBadgeDB := range allBadgesDB {
		if badgeID == correctedBadgeDB.BadgeID {
			return true
		}
	}
	return false
}

/**
 * To database model
 *
 * @return user database model
 */
func (u *User) ToDB() storages.UserDB {
	tokenExpireHour := time.Hour * utils.NewType().ParseDuration(os.Getenv("TOKEN_EXPIRE_HOUR"))
	expiredTokenTimestamp := time.Now().Local().Add(tokenExpireHour)
	return storages.UserDB{
		Name:                  u.name,
		Email:                 u.email,
		Password:              u.password,
		CreatedTimestamp:      u.createdTimestamp,
		UpdatedTimestamp:      u.updatedTimestamp,
		ExpiredTokenTimestamp: expiredTokenTimestamp,
	}
}

/**
 * Check user password
 *
 * @param 	password 	Password to check
 *
 * @return 	true if input password is correct, false if not
 */
func (u *User) IsPasswordCorrect(password string) bool {
	return utils.NewHelper().ComparePasswords(u.password, password)
}
