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
	email := req.GetEmail() // tainted
	req.Email = go_util.NormalizeEmail(email)

	userId := req.GetUserId() // not tainted
	queryUserReq := rpc_sdk.NewQueryUserRequest()
	queryUserReq.SetUserId(strconv.FormatInt(userId, 10))

	//get userInfo.Birthdate from rpc
	userInfo := rpc_sdk.RpcQueryUser(ctx, queryUserReq) // not tainted
	if len(userInfo.Birthdate) > 0 {
		//field email flow to rpc， 4 type case
		//case1, covered
		createEmailReq0 := &rpc_sdk.CreateEmailRequest{
			NewEmail: email, // tainted
		}
		rpc_sdk.RpcCreateEmail(ctx, createEmailReq0) // tainted, sink

		//case2, covered
		createEmailReq1 := &rpc_sdk.CreateEmailRequest{}
		createEmailReq1.NewEmail = email
		rpc_sdk.RpcCreateEmail(ctx, createEmailReq1) // tainted, sink

		//case3, covered
		emailInfo0 := rpc_sdk.EmailInfo{
			NewEmailV2: email,
			Extra:      "test_extra",
		}
		createEmailReqV20 := &rpc_sdk.CreateEmailV2Request{}
		createEmailReqV20.EmailInfo = emailInfo0
		rpc_sdk.RpcCreateEmailV2(ctx, createEmailReqV20)

		//case4, not covered
		emailInfo1 := rpc_sdk.EmailInfo{}
		emailInfo1.NewEmailV2 = email
		createEmailReqV21 := &rpc_sdk.CreateEmailV2Request{}
		createEmailReqV21.EmailInfo = emailInfo1
		rpc_sdk.RpcCreateEmailV2(ctx, createEmailReqV21) // tainted, sink

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
