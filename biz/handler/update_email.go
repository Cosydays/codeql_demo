package handler

import (
	"context"
	"fmt"
	"github.com/Cosydays/codeql_demo/constant"
	"github.com/Cosydays/codeql_demo/dal"
	"github.com/Cosydays/codeql_demo/go_util"
	"github.com/Cosydays/codeql_demo/model"
	"github.com/Cosydays/codeql_demo/rpc_sdk"
)

func UpdateEmail(ctx context.Context, req model.UpdateEmailRequest) {
	if req.GetUserId() > 0 {
		fmt.Println(req.GetUserId())
	}

	email := req.GetEmail()
	field := &rpc_sdk.Field{
		FieldType:  10,
		FieldValue: email,
	}

	//get field id from Redis
	id := dal.GetRedisValue("id")
	updateEmailReq := &rpc_sdk.UpdateEmailRequest{
		Id:    id,
		Field: field,
	}
	//field email flow to rpc
	rpc_sdk.RpcUpdateEmailInfo(ctx, updateEmailReq)
	//field email flow to redis
	redisKey := fmt.Sprintf(constant.KeyPatternA, "A")
	dal.SetValue2Redis(redisKey, email)

	//get field httpEmail from http
	httpEmail := GetHttpData(ctx, nil)
	httpField := &rpc_sdk.Field{
		FieldType:  11,
		FieldValue: httpEmail,
	}
	uEReq := &rpc_sdk.UpdateEmailRequest{
		Id:    id,
		Field: httpField,
	}
	//field email flow to rpc
	rpc_sdk.RpcUpdateEmailInfo(ctx, uEReq)

}

func GetHttpData(ctx context.Context, params map[string]string) string {
	_, err := go_util.HttpGet(ctx, "http://test/get/path", params, nil)
	if err != nil {
		fmt.Println(err)
	}
	return "email"
}
