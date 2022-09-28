package rpc_sdk

import (
	"context"
	"fmt"
)

//DeleteEmail
type DeleteEmailRequest struct {
	Id string `json:"id"`
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
	Id string `json:"id"`
	Field *Field `json:"field"`
}

type FieldType int64

type Field struct {
	FieldType      FieldType `json:"field_type"`
	FieldValue     string    `json:"field_value"`
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



//ChangePhone
type ChangePhoneReqest struct {
	Id string `json:"id"`
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
