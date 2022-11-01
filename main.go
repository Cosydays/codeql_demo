package main

import (
	"context"
	"fmt"
	"github.com/Cosydays/codeql_demo/client"
	"github.com/Cosydays/codeql_demo/constant"
	"github.com/Cosydays/codeql_demo/dal"
	"github.com/Cosydays/codeql_demo/go_util"
	"github.com/Cosydays/codeql_demo/model"
	"github.com/Cosydays/codeql_demo/rpc_sdk"
	"strconv"
)

func DeleteEmail(ctx context.Context, req model.DeleteEmailRequest) {
	if req.GetUserId() > 0 {
		fmt.Println(req.GetUserId())
	}
	email := req.GetEmail()
	deleteEmailReq := rpc_sdk.NewDeleteEmailRequest()
	deleteEmailReq.RpcEmail = email
	oldEmail := req.Info["old_email"]
	deleteEmailReq.OldEmail = oldEmail
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

	//field email flow to rpc
	rpc_sdk.RpcUpdateEmailInfo(ctx, updateEmailReq)

	//field email flow to redis
	redisKey := fmt.Sprintf(constant.KpA, "A")
	dal.SetValue2Redis(ctx, redisKey, email)
}

func CreateEmail(ctx context.Context, req model.CreateEmailRequest) {
	email := req.GetEmail()
	email = go_util.NormalizeEmail(email)
	req.Email = email

	userId := req.GetUserId()
	queryUserReq := rpc_sdk.NewQueryUserRequest()
	queryUserReq.SetUserId(strconv.FormatInt(userId, 10))
	userInfo := rpc_sdk.RpcQueryUser(ctx, queryUserReq)
	if len(userInfo.Birthdate) > 0 {
		createEmailReq := &rpc_sdk.CreateEmailRequest{
			NewEmail: email,
		}
		//field email flow to rpc
		rpc_sdk.RpcCreateEmail(ctx, createEmailReq)
	}
}

func ChangeEmail(ctx context.Context, req model.ChangeEmailRequest) {
	email := req.GetEmail()
	userId := req.GetUserId()

	queryUserReq := rpc_sdk.NewQueryUserRequest()
	queryUserReq.SetUserId(strconv.FormatInt(userId, 10))
	userInfo := rpc_sdk.RpcQueryUser(ctx, queryUserReq)

	//field email from RpcQueryUser
	rpcChangeEmailReq := &rpc_sdk.ChangeEmailRequest{
		OldEmail: userInfo.Email,
		NewEmail: email,
	}
	rpc_sdk.RpcChangeEmail(ctx, rpcChangeEmailReq)

}

func main() {
	ctx := context.Background()
	client.InitRedisClient(ctx)

	infoMap := map[string]string{}
	infoMap["old_email"] = "test5@email.com"
	deleteEmailReq := model.DeleteEmailRequest{
		UserId: 123,
		Email:  "test1@email.com",
		Info:   infoMap,
	}
	DeleteEmail(ctx, deleteEmailReq)

	updateEmailReq := model.UpdateEmailRequest{
		UserId: 456,
		Email:  "test2@email.com",
	}
	UpdateEmail(ctx, updateEmailReq)

	createEmailReq := model.CreateEmailRequest{
		UserId: 789,
		Email:  "test3@email.com",
		Status: "aa",
	}
	CreateEmail(ctx, createEmailReq)

	changeEmailReq := model.ChangeEmailRequest{
		UserId: 123,
		Email:  "test4@email.com",
	}
	ChangeEmail(ctx, changeEmailReq)
}
