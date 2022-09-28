package model

type DeleteEmailRequest struct {
	UserId int64  `json:"user_id"`
	Email  string `json:"email"`
}

func (deleteEmailRequest *DeleteEmailRequest) GetUserId() int64 {
	return deleteEmailRequest.UserId
}

func (deleteEmailRequest *DeleteEmailRequest) GetEmail() string {

}
