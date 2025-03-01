package v1

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/2mf8/Go-QQ-Client/dto"
	"github.com/2mf8/Go-QQ-Client/errs"
	"github.com/2mf8/Go-QQ-Client/openapi"
	"github.com/tidwall/gjson"
)

// Message 拉取单条消息
func (o *openAPI) Message(ctx context.Context, channelID string, messageID string) (*dto.Message, error) {
	resp, err := o.request(ctx).
		SetResult(dto.Message{}).
		SetPathParam("channel_id", channelID).
		SetPathParam("message_id", messageID).
		Get(o.getURL(messageURI))
	if err != nil {
		return nil, err
	}

	// 兼容处理
	result := resp.Result().(*dto.Message)
	if result.ID == "" {
		body := gjson.Get(resp.String(), "message")
		if err := json.Unmarshal([]byte(body.String()), result); err != nil {
			return nil, err
		}
	}
	return result, nil
}

// Messages 拉取消息列表
func (o *openAPI) Messages(ctx context.Context, channelID string, pager *dto.MessagesPager) ([]*dto.Message, error) {
	if pager == nil {
		return nil, errs.ErrPagerIsNil
	}
	resp, err := o.request(ctx).
		SetPathParam("channel_id", channelID).
		SetQueryParams(pager.QueryParams()).
		Get(o.getURL(messagesURI))
	if err != nil {
		return nil, err
	}

	messages := make([]*dto.Message, 0)
	if err := json.Unmarshal(resp.Body(), &messages); err != nil {
		return nil, err
	}

	return messages, nil
}

// PostMessage 发消息
func (o *openAPI) PostMessage(ctx context.Context, channelID string, msg *dto.MessageToCreate) (*dto.Message, error) {
	resp, err := o.request(ctx).
		SetResult(dto.Message{}).
		SetPathParam("channel_id", channelID).
		SetBody(msg).
		Post(o.getURL(messagesURI))
	if err != nil {
		return nil, err
	}

	return resp.Result().(*dto.Message), nil
}

func (o *openAPI) PostGroupMessage(ctx context.Context, groupID string, msg *dto.GroupMessageToCreate) (*dto.GroupMsgResp, error) {
	resp, err := o.request(ctx).
		SetResult(dto.GroupMsgResp{}).
		SetPathParam("group_openid", groupID).
		SetBody(msg).
		Post(o.getURL(groupMessageUri))
	if err != nil {
		return nil, err
	}

	return resp.Result().(*dto.GroupMsgResp), nil
}

func (o *openAPI) PostC2CMessage(ctx context.Context, userId string, msg *dto.C2CMessageToCreate) (*dto.C2CMsgResp, error) {
	resp, err := o.request(ctx).
		SetResult(dto.C2CMsgResp{}).
		SetPathParam("openid", userId).
		SetBody(msg).
		Post(o.getURL(privateMessageUri))
	if err != nil {
		return nil, err
	}

	return resp.Result().(*dto.C2CMsgResp), nil
}

func (o *openAPI) PostC2CRichMediaMessage(ctx context.Context, userId string, msg *dto.C2CRichMediaMessageToCreate) (*dto.RichMediaMsgResp, error) {
	resp, err := o.request(ctx).
		SetResult(dto.RichMediaMsgResp{}).
		SetPathParam("openid", userId).
		SetBody(msg).
		Post(o.getURL(privateRichMediaMessageUri))
	if err != nil {
		return nil, err
	}

	return resp.Result().(*dto.RichMediaMsgResp), nil
}

func (o *openAPI) PostGroupRichMediaMessage(ctx context.Context, groupID string, msg *dto.GroupRichMediaMessageToCreate) (*dto.RichMediaMsgResp, error) {
	resp, err := o.request(ctx).
		SetResult(dto.RichMediaMsgResp{}).
		SetPathParam("group_openid", groupID).
		SetBody(msg).
		Post(o.getURL(groupRichMediaMessageUri))
	if err != nil {
		return nil, err
	}

	return resp.Result().(*dto.RichMediaMsgResp), nil
}

// PatchMessage 编辑消息
func (o *openAPI) PatchMessage(ctx context.Context,
	channelID string, messageID string, msg *dto.MessageToCreate) (*dto.Message, error) {
	resp, err := o.request(ctx).
		SetResult(dto.Message{}).
		SetPathParam("channel_id", channelID).
		SetPathParam("message_id", messageID).
		SetBody(msg).
		Patch(o.getURL(messageURI))
	if err != nil {
		return nil, err
	}

	return resp.Result().(*dto.Message), nil
}

// RetractMessage 撤回消息
func (o *openAPI) RetractMessage(ctx context.Context,
	channelID, msgID string, options ...openapi.RetractMessageOption) error {
	request := o.request(ctx).
		SetPathParam("channel_id", channelID).
		SetPathParam("message_id", string(msgID))
	for _, option := range options {
		if option == openapi.RetractMessageOptionHidetip {
			request = request.SetQueryParam("hidetip", "true")
		}
	}
	_, err := request.Delete(o.getURL(messageURI))
	return err
}

// PostSettingGuide 发送设置引导消息, atUserID为要at的用户
func (o *openAPI) PostSettingGuide(ctx context.Context,
	channelID string, atUserIDs []string) (*dto.Message, error) {
	var content string
	for _, userID := range atUserIDs {
		content += fmt.Sprintf("<@%s>", userID)
	}
	msg := &dto.SettingGuideToCreate{
		Content: content,
	}
	resp, err := o.request(ctx).
		SetResult(dto.Message{}).
		SetPathParam("channel_id", channelID).
		SetBody(msg).
		Post(o.getURL(settingGuideURI))
	if err != nil {
		return nil, err
	}
	return resp.Result().(*dto.Message), nil
}
