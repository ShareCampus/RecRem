package forms

type QueryInfoForm struct {
	UserID   string `json:"userid" binding:"required"`
	Question string `json:"question" binding:"required"`
}
