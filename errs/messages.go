package errs

// Thai and english message about server error
const (
	INTERNAL_SERVER_ERROR_TH = "เกิดข้อผิดพลาด"
	INTERNAL_SERVER_ERROR_EN = "Internal server error"

	SERVICE_UNAVAILABLE_ERROR_TH = "บริการไม่พร้อมใช้งาน"
	SERVICE_UNAVAILABLE_ERROR_EN = "Service unavailable"

	BAD_REQUEST_ERROR_TH = "คำร้องขอไม่ถูกต้อง"
	BAD_REQUEST_ERROR_EN = "Bad request"
)

// Thai and english message for insert error, load error and update error
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

	HINTS_ALREADY_USED_TH = "ได้ใช้คำใบ้ทั้งหมดของกิจกรรมแล้ว"
	HINTS_ALREADY_USED_EN = "Activity hints has been used"

	HINT_POINTS_NOT_ENOUGH_TH = "แต้มไม่เพียงพอในการขอคำใบ้"
	HINT_POINTS_NOT_ENOUGH_EN = "Not enough points to use a hint"

	ACTIVITY_TYPE_INVALID_TH = "ประเภทของกิจกรรมไม่ถูกต้อง"
	ACTIVITY_TYPE_INVALID_EN = "Activity type invalid"

	FINAL_EXAM_BAGES_NOT_ENOUGH_TH = "จำนวนเหรียญตราไม่เพียงพอในการทำข้อสอบ"
	FINAL_EXAM_BAGES_NOT_ENOUGH_EN = "Not enough badges to do final exam"
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

// Thai and english message for unexpected signing method
const (
	UNEXPECTED_SIGNING_METHOD_TH = "วิธีการลงนามที่ไม่คาดคิด"
	UNEXPECTED_SIGNING_METHOD_EN = "Unexpected signing method"
)

var (
	ErrInternalServerError     = NewInternalServerError(INTERNAL_SERVER_ERROR_TH, INTERNAL_SERVER_ERROR_EN)
	ErrServiceUnavailableError = NewServiceUnavailableError(SERVICE_UNAVAILABLE_ERROR_TH, SERVICE_UNAVAILABLE_ERROR_EN)
	ErrBadRequestError         = NewBadRequestError(BAD_REQUEST_ERROR_TH, BAD_REQUEST_ERROR_EN)
)

var (
	ErrInsertError = NewInternalServerError(INSERT_ERROR_TH, INSERT_ERROR_EN)
	ErrLoadError   = NewInternalServerError(LOAD_ERROR_TH, LOAD_ERROR_EN)
	ErrUpdateError = NewInternalServerError(UPDATE_ERROR_TH, UPDATE_ERROR_EN)
)

var (
	ErrExamNotFound              = NewNotFoundError(EXAM_NOT_FOUND_TH, EXAM_NOT_FOUND_EN)
	ErrContentNotFound           = NewNotFoundError(CONTENT_NOT_FOUND_TH, CONTENT_NOT_FOUND_EN)
	ErrActivitiesNotFound        = NewNotFoundError(ACTIVITIES_NOT_FOUND_TH, ACTIVITIES_NOT_FOUND_EN)
	ErrActivitiesNumberIncorrect = NewBadRequestError(ACTIVITIES_NUMBER_INCORRECT_TH, ACTIVITIES_NUMBER_INCORRECT_EN)
	ErrHintAlreadyUsed           = NewBadRequestError(HINTS_ALREADY_USED_TH, HINTS_ALREADY_USED_EN)
	ErrHintPointsNotEnough       = NewBadRequestError(HINT_POINTS_NOT_ENOUGH_TH, HINT_POINTS_NOT_ENOUGH_EN)
	ErrActivityTypeInvalid       = NewInternalServerError(ACTIVITY_TYPE_INVALID_TH, ACTIVITY_TYPE_INVALID_EN)
	ErrFinalExamBadgesNotEnough  = NewBadRequestError(FINAL_EXAM_BAGES_NOT_ENOUGH_TH, FINAL_EXAM_BAGES_NOT_ENOUGH_EN)
)

var (
	ErrUserNotFound              = NewNotFoundError(USER_NOT_FOUND_TH, USER_NOT_FOUND_EN)
	ErrLeaderBoardNotFound       = NewNotFoundError(LEADER_BOARD_NOT_FOUND_TH, LEADER_BOARD_NOT_FOUND_EN)
	ErrEmailAlreadyExists        = NewBadRequestError(EMAIL_ALREADY_EXISTS_TH, EMAIL_ALREADY_EXISTS_EN)
	ErrEmailOrPasswordNotCorrect = NewBadRequestError(EMAIL_OR_PASSWORD_NOT_CORRECT_TH, EMAIL_OR_PASSWORD_NOT_CORRECT_EN)
)

var (
	ErrUnExpectedsigningMethod = NewForbiddenError(UNEXPECTED_SIGNING_METHOD_TH, UNEXPECTED_SIGNING_METHOD_EN)
)
