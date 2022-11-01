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

type CreateEmailRequest struct {
	UserId int64  `json:"user_id"`
	Email  string `json:"email"`
	Status string `json:"status"`
}

func (p *CreateEmailRequest) GetUserId() int64 {
	return p.UserId
}

func (p *CreateEmailRequest) GetEmail() string {
	return p.Email
}

func (p *CreateEmailRequest) GetStatus() string {
	return p.Status
}

type ChangeEmailRequest struct {
	UserId int64  `json:"user_id"`
	Email  string `json:"email"`
}

func (p *ChangeEmailRequest) GetUserId() int64 {
	return p.UserId
}

func (p *ChangeEmailRequest) GetEmail() string {
	return p.Email
}
