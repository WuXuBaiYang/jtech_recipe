package common

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"server/model"
	"server/tool"
	"strings"
	"time"
)

// im_api服务地址
var baseImApi = fmt.Sprintf("http://%s:10002", ServerHost)

// im服务授权
const imSecret = "8OVL0bcFGEozHW"

// RegisterOnIM 注册到im
func RegisterOnIM(platform int, userId string) (*ImAuthRes, error) {
	url := fmt.Sprintf("%s/auth/user_register", baseImApi)
	operationID := genOperationID("register", platform, userId)
	reqBody, _ := json.Marshal(map[string]any{
		"operationID": operationID,
		"platform":    platform,
		"secret":      imSecret,
		"userID":      userId,
	})
	resp, err := http.Post(url, "application/json", strings.NewReader(string(reqBody)))
	if err != nil {
		return nil, err
	}
	resBody := &ImAuthRes{}
	if err := handleIMAuthRequest(resp, resBody); err != nil {
		return nil, err
	}
	return resBody, nil
}

// GetIMUserToken 获取IM登录授权
func GetIMUserToken(platform int, userId string) (*ImAuthRes, error) {
	url := fmt.Sprintf("%s/auth/user_token", baseImApi)
	operationID := genOperationID("get_user_token", platform, userId)
	reqBody, _ := json.Marshal(map[string]any{
		"operationID": operationID,
		"platform":    platform,
		"secret":      imSecret,
		"userID":      userId,
	})
	resp, err := http.Post(url, "application/json", strings.NewReader(string(reqBody)))
	if err != nil {
		return nil, err
	}
	resBody := &ImAuthRes{}
	if err := handleIMAuthRequest(resp, resBody); err != nil {
		return nil, err
	}
	return resBody, nil
}

// UpdateIMUserInfo 更新用户信息
func UpdateIMUserInfo(platform int, user model.User) error {
	if user.Profile == nil {
		return errors.New("用户信息不能为空")
	}
	url := fmt.Sprintf("%s/user/update_user_info", baseImApi)
	operationID := genOperationID("update_user", platform, user.IMUserId)
	reqBody, _ := json.Marshal(map[string]any{
		"operationID": operationID,
		"userID":      user.IMUserId,
		"gender":      user.Profile.Gender,
		"nickname":    user.Profile.NickName,
		"faceURL":     user.Profile.Avatar,
	})
	client := &http.Client{}
	req, err := http.NewRequest("POST", url, strings.NewReader(string(reqBody)))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("token", user.IMToken)
	resp, errResp := client.Do(req)
	if errResp != nil {
		return err
	}
	resBody := &ImAuthRes{}
	if err := handleIMAuthRequest(resp, resBody); err != nil {
		return err
	}
	return nil
}

// 生成操作id
func genOperationID(action string, platform int, userId string) string {
	platformStr := tool.Platform2Tag(platform)
	timestamp := time.Now().UnixMilli()
	return fmt.Sprintf("%s_%s_%s_%d", action, *platformStr, userId, timestamp)
}

// 处理im_api授权请求
func handleIMAuthRequest(resp *http.Response, resBody *ImAuthRes) error {
	if resp.StatusCode == 200 {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		err = json.Unmarshal(body, &resBody)
		if err != nil {
			return err
		}
		if resBody.ErrCode != 0 {
			return errors.New(resBody.ErrMsg)
		}
		return nil
	}
	return errors.New(resp.Status)
}

// ImAuthRes im_api 实体
type ImAuthRes struct {
	ErrCode int    `json:"errCode"`
	ErrMsg  string `json:"errMsg"`
	Data    struct {
		UserID      string `json:"userID"`
		Token       string `json:"token"`
		ExpiredTime int64  `json:"expiredTime"`
	} `json:"data"`
}
