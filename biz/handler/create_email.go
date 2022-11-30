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
	email := req.GetEmail()               //COMPLIANT, The Field `email` to be taint, As SourceNode
	email = go_util.NormalizeEmail(email) //NON_COMPLIANT, encrypt email
	req.Email = email                     //NON_COMPLIANT, Reassign the value to req.Email

	userId := req.GetUserId() //NON_COMPLIANT, get the userId
	queryUserReq := rpc_sdk.NewQueryUserRequest()
	queryUserReq.SetUserId(strconv.FormatInt(userId, 10))

	//get userInfo.Birthdate from rpc
	userInfo := rpc_sdk.RpcQueryUser(ctx, queryUserReq) // NON_COMPLIANT, get the userInfo by userId
	if len(userInfo.Birthdate) > 0 {
		createEmailReq := &rpc_sdk.CreateEmailRequest{
			NewEmail: email, //COMPLIANT
		}

		//field email flow to rpc
		rpc_sdk.RpcCreateEmail(ctx, createEmailReq) //COMPLIANT, As SinkNode

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
