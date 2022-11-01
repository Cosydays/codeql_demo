package rpc_sdk

import (
	"context"
	"fmt"
)

type UserInfo struct {
	Name      string `json:"name"`
	Birthdate string `json:"birthdate"`
	Email     string `json:"email"`
}

//DeleteEmail
type DeleteEmailRequest struct {
	Id       string `json:"id"`
	RpcEmail string `json:"rpc_email"`
	OldEmail string `json:"old_email"`
}

func NewDeleteEmailRequest() *DeleteEmailRequest {
	return &DeleteEmailRequest{}
}

func (p *DeleteEmailRequest) GetId() (v string) {
	return p.Id
}

func (p *DeleteEmailRequest) GetRpcEmail() (v string) {
	return p.RpcEmail
}

func (p *DeleteEmailRequest) SetId(v string) {
	p.Id = v
}

func (p *DeleteEmailRequest) SetRpcEmail(v string) {
	p.RpcEmail = v
}

func RpcDeleteEmailInfo(ctx context.Context, req *DeleteEmailRequest) {
	fmt.Println(req.RpcEmail)
}

//UpdateEmail
type UpdateEmailRequest struct {
	Id    string `json:"id"`
	Field *Field `json:"field"`
}

type FieldType int64

type Field struct {
	FieldType  FieldType `json:"field_type"`
	FieldValue string    `json:"field_value"`
}

func NewUpdateEmailRequest() *UpdateEmailRequest {
	return &UpdateEmailRequest{}
}

func (p *UpdateEmailRequest) GetId() (v string) {
	return p.Id
}

func (p *UpdateEmailRequest) SetId(v string) {
	p.Id = v
}

func RpcUpdateEmailInfo(ctx context.Context, req *UpdateEmailRequest) {
	fmt.Println(req.Field.FieldValue)
}

//CreateEmail
type CreateEmailRequest struct {
	Id       string `json:"id"`
	NewEmail string `json:"new_email"`
}

func NewCreateEmailRequest() *CreateEmailRequest {
	return &CreateEmailRequest{}
}

func (p *CreateEmailRequest) GetId() (v string) {
	return p.Id
}

func (p *CreateEmailRequest) GetNewEmail() (v string) {
	return p.NewEmail
}

func (p *CreateEmailRequest) SetId(v string) {
	p.Id = v
}

func (p *CreateEmailRequest) SetNewEmail(v string) {
	p.NewEmail = v
}

func RpcCreateEmail(ctx context.Context, req *CreateEmailRequest) {
	fmt.Println(req.NewEmail)
}

//ChangeEmail
type ChangeEmailRequest struct {
	Id       string `json:"id"`
	OldEmail string `json:"old_email"`
	NewEmail string `json:"new_email"`
}

func NewChangeEmailRequest() *ChangeEmailRequest {
	return &ChangeEmailRequest{}
}

func (p *ChangeEmailRequest) GetId() (v string) {
	return p.Id
}

func (p *ChangeEmailRequest) GetOldEmail() (v string) {
	return p.OldEmail
}

func (p *ChangeEmailRequest) GetNewEmail() (v string) {
	return p.NewEmail
}

func (p *ChangeEmailRequest) SetId(v string) {
	p.Id = v
}

func (p *ChangeEmailRequest) SetOldEmail(v string) {
	p.OldEmail = v
}

func (p *ChangeEmailRequest) SetNewEmail(v string) {
	p.NewEmail = v
}

func RpcChangeEmail(ctx context.Context, req *ChangeEmailRequest) {
	fmt.Println(req.OldEmail)
	fmt.Println(req.NewEmail)
}

//QueryUser
type QueryUserRequest struct {
	UserId string `json:"user_id"`
}

func NewQueryUserRequest() *QueryUserRequest {
	return &QueryUserRequest{}
}

func (p *QueryUserRequest) GetUserId() (v string) {
	return p.UserId
}

func (p *QueryUserRequest) SetUserId(v string) {
	p.UserId = v
}

func RpcQueryUser(ctx context.Context, req *QueryUserRequest) *UserInfo {
	userInfo := UserInfo{
		Name:      "name_test",
		Birthdate: "99",
		Email:     "s.@a.com",
	}
	return &userInfo
}
