package models

var TableName = struct {
	User string
	UserBadge string
	Badge string
}{
	"User",
	"UserBadge",
	"Badge",
}

var IDName = struct {
	User string
}{
	"user_id",
}

var ViewName = struct{
	Profile string
	Ranking string
}{
	"Profile",
	"Ranking",
}
