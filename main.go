package main

import (
	"context"
	"fmt"
	"main/model"
	"main/rpc_sdk"
)

func DeleteEmail(ctx context.Context, req model.DeleteEmailRequest) {
	if req.GetUserId() > 0 {
		fmt.Println(req.GetUserId())
	}
	email := req.GetEmail()
	deleteEmailReq := rpc_sdk.NewDeleteEmailRequest()
	deleteEmailReq.RpcEmail = email
	rpc_sdk.RpcDeleteEmailInfo(ctx, deleteEmailReq)
}

func UpdateEmail(ctx context.Context, req model.UpdateEmailRequest) {
	if req.GetUserId() > 0 {
		fmt.Println(req.GetUserId())
	}

	email := req.GetEmail()
	field := &rpc_sdk.Field{
		FieldType:  10,
		FieldValue: email,
	}

	updateEmailReq := &rpc_sdk.UpdateEmailRequest{
		Id:    "123",
		Field: field,
	}
	rpc_sdk.RpcUpdateEmailInfo(ctx, updateEmailReq)
}

func main() {
	ctx := context.Background()
	deleteEmailReq := model.DeleteEmailRequest{
		UserId: 123,
		Email:  "test1@email.com",
	}
	DeleteEmail(ctx, deleteEmailReq)
	updateEmailReq := model.UpdateEmailRequest{
		UserId: 456,
		Email:  "test2@email.com",
	}
	UpdateEmail(ctx, updateEmailReq)
}
