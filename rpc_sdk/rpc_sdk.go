package rpc_sdk

import (
	"context"
	"fmt"
)

//DeleteEmail
type DeleteEmailRequest struct {
	Id       string `json:"id"`
	RpcEmail string `json:"rpc_email"`
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

func RpcQueryUser(ctx context.Context, req *QueryUserRequest) {
	fmt.Println(req.UserId)
}

//ChangePhone
type ChangePhoneReqest struct {
	Id       string `json:"id"`
	RpcPhone string `json:"rpc_phone"`
}

func NewChangePhoneReqest() *ChangePhoneReqest {
	return &ChangePhoneReqest{}
}

func (p *ChangePhoneReqest) GetId() (v string) {
	return p.Id
}

func (p *ChangePhoneReqest) GetPhone() (v string) {
	return p.RpcPhone
}

func (p *ChangePhoneReqest) SetId(v string) {
	p.Id = v
}

func (p *ChangePhoneReqest) SetRpcPhone(v string) {
	p.RpcPhone = v
}

func ChangePhoneInfo(ctx context.Context, req *ChangePhoneReqest) {
	fmt.Println(req.RpcPhone)
}
