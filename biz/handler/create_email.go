package handler

import (
	"context"
	"fmt"
	"github.com/Cosydays/codeql_demo/go_util"
	"github.com/Cosydays/codeql_demo/model"
	"github.com/Cosydays/codeql_demo/rpc_sdk"
	"strconv"
	"time"
)

func CreateEmail(ctx context.Context, req model.CreateEmailRequest) {
	email := req.GetEmail()
	email = go_util.NormalizeEmail(email)
	req.Email = email

	userId := req.GetUserId()
	queryUserReq := rpc_sdk.NewQueryUserRequest()
	queryUserReq.SetUserId(strconv.FormatInt(userId, 10))

	//get userInfo.Birthdate from rpc
	userInfo := rpc_sdk.RpcQueryUser(ctx, queryUserReq)
	if len(userInfo.Birthdate) > 0 {
		createEmailReq := &rpc_sdk.CreateEmailRequest{
			NewEmail: email,
		}

		//field email flow to rpc
		rpc_sdk.RpcCreateEmail(ctx, createEmailReq)

		//field email flow to http
		CallHttp(ctx, email)
	}
}

func CallHttp(ctx context.Context, email string) {
	_, err := go_util.HttpPost(ctx, "http://test/post/path", []byte(email), nil, time.Second*30)
	if err != nil {
		fmt.Println(err)
	}
}
