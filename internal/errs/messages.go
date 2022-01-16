package errs

// messages.go
/**
 * 	This file used to defind all error in the application
 */

// Thai and english message about server error
const (
	INTERNAL_SERVER_ERROR_TH = "เกิดข้อผิดพลาด"
	INTERNAL_SERVER_ERROR_EN = "Internal server error"

	SERVICE_UNAVAILABLE_ERROR_TH = "บริการไม่พร้อมใช้งาน"
	SERVICE_UNAVAILABLE_ERROR_EN = "Service unavailable"

	BAD_REQUEST_ERROR_TH = "คำร้องขอไม่ถูกต้อง"
	BAD_REQUEST_ERROR_EN = "Bad request"
)

// Thai and english message about manipulation error
const (
	INSERT_ERROR_TH = "เกิดข้อผิดพลาดในการบันทึกข้อมูล"
	INSERT_ERROR_EN = "Saving data error"

	LOAD_ERROR_TH = "เกิดข้อผิดพลาดในการโหลดข้อมูล"
	LOAD_ERROR_EN = "Load data error"

	UPDATE_ERROR_TH = "เกิดข้อผิดพลาดในการอัพเดตข้อมูล"
	UPDATE_ERROR_EN = "Update data error"
)

// Thai and english message about exam and activity error
const (
	EXAM_NOT_FOUND_TH = "ไม่พบข้อสอบ"
	EXAM_NOT_FOUND_EN = "Exam not found"

	CONTENT_NOT_FOUND_TH = "ไม่พบเนื้อหา"
	CONTENT_NOT_FOUND_EN = "Content not found"

	ACTIVITIES_NOT_FOUND_TH = "ไม่พบกิจกรรม"
	ACTIVITIES_NOT_FOUND_EN = "Activities not found"

	ACTIVITIES_NUMBER_INCORRECT_TH = "จำนวนของกิจกรรมไม่ถูกต้อง"
	ACTIVITIES_NUMBER_INCORRECT_EN = "Number of activities incorrect"

	ACTIVITYID_NOTFOUND_TH = "ไม่พบรหัสของกิจกรรมในคำร้องขอ"
	ACTIVITYID_NOTFOUND_EN = "Activity ID not found"

	ACTIVITY_TYPEID_NOTFOUND_TH = "ไม่พบรหัสของประเภทกิจกรรมในคำร้องขอ"
	ACTIVITY_TYPEID_NOTFOUND_EN = "Activity Type ID not found"

	HINTS_ALREADY_USED_TH = "ได้ใช้คำใบ้ทั้งหมดของกิจกรรมแล้ว"
	HINTS_ALREADY_USED_EN = "Activity hints has been used"

	HINT_POINTS_NOT_ENOUGH_TH = "แต้มไม่เพียงพอในการขอคำใบ้"
	HINT_POINTS_NOT_ENOUGH_EN = "Not enough points to use a hint"

	ACTIVITY_TYPE_INVALID_TH = "ประเภทของกิจกรรมไม่ถูกต้อง"
	ACTIVITY_TYPE_INVALID_EN = "Activity type invalid"

	FINAL_EXAM_BAGES_NOT_ENOUGH_TH = "จำนวนเหรียญตราไม่เพียงพอในการทำข้อสอบ"
	FINAL_EXAM_BAGES_NOT_ENOUGH_EN = "Not enough badges to do final exam"

	ANSWER_INVALID_TH = "รูปแบบของคำตอบไม่ถูกต้อง"
	ANSWER_INVALID_EN = "Answer invalid"

	ANSWER_NOT_FOUND_TH = "ไม่พบคำตอบในคำร้องขอ"
	ANSWER_NOT_FOUND_EN = "Answer not found"
)

// Thai and english message about user error
const (
	USER_NOT_FOUND_TH = "ไม่พบผู้ใช้"
	USER_NOT_FOUND_EN = "User not found"

	LEADER_BOARD_NOT_FOUND_TH = "ไม่พบตารางคะแนน"
	LEADER_BOARD_NOT_FOUND_EN = "Leader board not found"

	EMAIL_ALREADY_EXISTS_TH = "อีเมลมีการใช้งานแล้ว"
	EMAIL_ALREADY_EXISTS_EN = "Email is already exists"

	EMAIL_OR_PASSWORD_NOT_CORRECT_TH = "อีเมลหรือรหัสผ่านไม่ถูกต้อง"
	EMAIL_OR_PASSWORD_NOT_CORRECT_EN = "Email or password not correct"
)

// Thai and english message about verification
const (
	UNEXPECTED_SIGNING_METHOD_TH = "วิธีการลงนามที่ไม่คาดคิด"
	UNEXPECTED_SIGNING_METHOD_EN = "Unexpected signing method"
)

// Server error
var (
	ErrInternalServerError     = NewInternalServerError(INTERNAL_SERVER_ERROR_TH, INTERNAL_SERVER_ERROR_EN)
	ErrServiceUnavailableError = NewServiceUnavailableError(SERVICE_UNAVAILABLE_ERROR_TH, SERVICE_UNAVAILABLE_ERROR_EN)
	ErrBadRequestError         = NewBadRequestError(BAD_REQUEST_ERROR_TH, BAD_REQUEST_ERROR_EN)
)

// Manipulation error
var (
	ErrInsertError = NewInternalServerError(INSERT_ERROR_TH, INSERT_ERROR_EN)
	ErrLoadError   = NewInternalServerError(LOAD_ERROR_TH, LOAD_ERROR_EN)
	ErrUpdateError = NewInternalServerError(UPDATE_ERROR_TH, UPDATE_ERROR_EN)
)

// Exam and Activity error
var (
	ErrExamNotFound              = NewNotFoundError(EXAM_NOT_FOUND_TH, EXAM_NOT_FOUND_EN)
	ErrContentNotFound           = NewNotFoundError(CONTENT_NOT_FOUND_TH, CONTENT_NOT_FOUND_EN)
	ErrActivitiesNotFound        = NewNotFoundError(ACTIVITIES_NOT_FOUND_TH, ACTIVITIES_NOT_FOUND_EN)
	ErrActivittyIDNotFound       = NewNotFoundError(ACTIVITYID_NOTFOUND_TH, ACTIVITYID_NOTFOUND_EN)
	ErrActivittyTypeIDNotFound   = NewNotFoundError(ACTIVITY_TYPEID_NOTFOUND_TH, ACTIVITY_TYPEID_NOTFOUND_EN)
	ErrActivitiesNumberIncorrect = NewBadRequestError(ACTIVITIES_NUMBER_INCORRECT_TH, ACTIVITIES_NUMBER_INCORRECT_EN)
	ErrHintAlreadyUsed           = NewBadRequestError(HINTS_ALREADY_USED_TH, HINTS_ALREADY_USED_EN)
	ErrHintPointsNotEnough       = NewBadRequestError(HINT_POINTS_NOT_ENOUGH_TH, HINT_POINTS_NOT_ENOUGH_EN)
	ErrActivityTypeInvalid       = NewInternalServerError(ACTIVITY_TYPE_INVALID_TH, ACTIVITY_TYPE_INVALID_EN)
	ErrAnswerInvalid             = NewBadRequestError(ANSWER_INVALID_TH, ANSWER_INVALID_EN)
	ErrAnswerNotFound            = NewBadRequestError(ANSWER_NOT_FOUND_TH, ANSWER_NOT_FOUND_EN)
	ErrFinalExamBadgesNotEnough  = NewBadRequestError(FINAL_EXAM_BAGES_NOT_ENOUGH_TH, FINAL_EXAM_BAGES_NOT_ENOUGH_EN)
)

// User error
var (
	ErrUserNotFound              = NewNotFoundError(USER_NOT_FOUND_TH, USER_NOT_FOUND_EN)
	ErrLeaderBoardNotFound       = NewNotFoundError(LEADER_BOARD_NOT_FOUND_TH, LEADER_BOARD_NOT_FOUND_EN)
	ErrEmailAlreadyExists        = NewBadRequestError(EMAIL_ALREADY_EXISTS_TH, EMAIL_ALREADY_EXISTS_EN)
	ErrEmailOrPasswordNotCorrect = NewBadRequestError(EMAIL_OR_PASSWORD_NOT_CORRECT_TH, EMAIL_OR_PASSWORD_NOT_CORRECT_EN)
)

// Verification error
var (
	ErrUnExpectedsigningMethod = NewForbiddenError(UNEXPECTED_SIGNING_METHOD_TH, UNEXPECTED_SIGNING_METHOD_EN)
)
