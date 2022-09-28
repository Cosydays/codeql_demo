package model

type DeleteEmailRequest struct {
	UserId int64  `json:"user_id"`
	Email  string `json:"email"`
}

func (p *DeleteEmailRequest) GetUserId() int64 {
	return p.UserId
}

func (p *DeleteEmailRequest) GetEmail() string {
	return p.Email
}

type UpdateEmailRequest struct {
	UserId int64  `json:"user_id"`
	Email  string `json:"email"`
}

func (p *UpdateEmailRequest) GetUserId() int64 {
	return p.UserId
}

func (p *UpdateEmailRequest) GetEmail() string {
	return p.Email
}
