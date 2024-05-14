package recrem

import (
	"net/http"

	"github.com/Edgenesis/IoT-Driver-Copilot/backend/pkg/database"
	"github.com/Edgenesis/IoT-Driver-Copilot/backend/pkg/utils/httpjson"
	"github.com/Edgenesis/IoT-Driver-Copilot/backend/pkg/utils/mysqlerror"
	"gorm.io/gorm"
)

type SignUpRequest struct {
	UserId string `json:"user_id" validate:"required,len=28"`
}

func (s *server) SignUp(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	request, err := httpjson.BindJson[SignUpRequest](r.Body)
	if err != nil {
		httpjson.ReturnError(w, http.StatusBadRequest, "invalid request")
		return
	}

	if err := createUser(s.dbClient, request.UserId); err != nil {
		if mysqlerror.CheckError(err, mysqlerror.ErrDuplicateEntryCode) {
			httpjson.ReturnError(w, http.StatusConflict, "user already exists")
			return
		}

		httpjson.ReturnError(w, http.StatusInternalServerError, "failed to create user")
		return
	}

	if err := httpjson.WriteJson[*struct{}](w, nil); err != nil {
		httpjson.ReturnError(w, http.StatusInternalServerError, "failed to write response")
	}
}

func createUser(txn *gorm.DB, userId string) error {
	user := database.User{
		Id: userId,
	}

	return txn.Create(&user).Error
}
