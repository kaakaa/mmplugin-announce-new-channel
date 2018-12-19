package main

import (
	"fmt"

	"github.com/mattermost/mattermost-server/model"
	"github.com/mattermost/mattermost-server/plugin"
)

const MessageFormat = `Channel ~%s has been created by %s.`

func (p *Plugin) ChannelHasBeenCreated(c *plugin.Context, channel *model.Channel) {
	if channel.Type != model.CHANNEL_OPEN {
		return
	}

	u, appErr := p.API.GetUser(channel.CreatorId)
	if appErr != nil {
		p.API.LogError("Failed to get user", "details", appErr)
		return
	}
	townSquare, appErr := p.API.GetChannelByName(channel.TeamId, model.DEFAULT_CHANNEL, false)
	if appErr != nil {
		p.API.LogError("Failed to get channel", "details", appErr)
		return
	}

	post := &model.Post{
		Type:      model.POST_DEFAULT,
		ChannelId: townSquare.Id,
		UserId:    u.Id,
		Message:   fmt.Sprintf(MessageFormat, channel.Name, u.GetDisplayName(model.SHOW_USERNAME)),
	}

	overrideIconUrl := p.getConfiguration().IconURL
	overrideUserNmae := p.getConfiguration().UserName
	props := map[string]interface{}{}
	if overrideIconUrl != "" || overrideUserNmae != "" {
		if overrideIconUrl != "" {
			props["override_icon_url"] = overrideIconUrl
		}
		if overrideUserNmae != "" {
			props["override_username"] = overrideUserNmae
		}
	}
	props["from_webhook"] = "true"
	post.Props = props

	if _, appErr := p.API.CreatePost(post); appErr != nil {
		p.API.LogError("Failed to create post", "details", appErr)
	}
}
