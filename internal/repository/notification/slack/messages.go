package slack

import (
	"context"
	"fmt"

	"github.com/alvintzz/alert-thread/internal/entity"
	sl "github.com/slack-go/slack"
)

var errorEmptyThreadID = fmt.Errorf("ThreadID is required")

// SendMessage will send new message to Slack. If thread_id is provided, the message will go to the thread. Otherwise it will create a new thread
func (s *Slack) SendMessage(ctx context.Context, param entity.Notification) (string, error) {
	attachment := sl.Attachment{
		Color: param.GetColor(),
		Blocks: sl.Blocks{
			BlockSet: []sl.Block{
				sl.NewSectionBlock(sl.NewTextBlockObject(sl.MarkdownType, param.Message, false, false), nil, nil),
			},
		},
	}
	if param.Image != "" {
		attachment.Blocks.BlockSet = append(attachment.Blocks.BlockSet, sl.NewImageBlock(param.Image, "alt text", "", nil))
	}

	options := []sl.MsgOption{
		sl.MsgOptionText(param.Title, false),
		sl.MsgOptionAttachments(attachment),
	}
	if value, ok := param.Metadata["timestamp"]; ok && value != "" {
		options = append(options, sl.MsgOptionTS(value))
	}

	_, threadID, err := s.client.PostMessageContext(ctx, param.Channel, options...)
	if err != nil {
		return threadID, err
	}

	return threadID, nil
}

// UpdateMessage will update the thread message provided
func (s *Slack) UpdateMessage(ctx context.Context, param entity.Notification) error {
	attachment := sl.Attachment{
		Color: param.GetColor(),
		Blocks: sl.Blocks{
			BlockSet: []sl.Block{
				sl.NewSectionBlock(sl.NewTextBlockObject(sl.MarkdownType, param.Message, false, false), nil, nil),
			},
		},
	}
	if param.Image != "" {
		attachment.Blocks.BlockSet = append(attachment.Blocks.BlockSet, sl.NewImageBlock(param.Image, "alt text", "", nil))
	}

	options := []sl.MsgOption{
		sl.MsgOptionText(param.Title, true),
		sl.MsgOptionAttachments(attachment),
	}

	threadID, ok := param.Metadata["timestamp"]
	if !ok || threadID == "" {
		return errorEmptyThreadID
	}

	_, _, _, err := s.client.UpdateMessageContext(ctx, param.Channel, threadID, options...)
	if err != nil {
		return err
	}

	return nil
}
