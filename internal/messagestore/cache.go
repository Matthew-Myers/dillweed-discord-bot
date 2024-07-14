package messagestore

import "internal/botconfig"

type MessageCache struct {
	data       map[string]map[string]*ChannelStore
	maxRecords int
	config     *botconfig.AppConfig
}

type ChannelStore struct {
	messages         []MessageRecord
	personalityIndex int
}

func NewMessageCache(maxRecords int, config *botconfig.AppConfig) *MessageCache {
	return &MessageCache{
		data:       make(map[string]map[string]*ChannelStore),
		maxRecords: maxRecords,
		config:     config,
	}
}

func (c *MessageCache) Put(serverId, channelId string, record MessageRecord) {
	if _, exists := c.data[serverId]; !exists {
		c.data[serverId] = make(map[string]*ChannelStore)
	}

	if _, exists := c.data[serverId][channelId]; !exists {
		c.data[serverId][channelId] = &ChannelStore{
			messages:         []MessageRecord{},
			personalityIndex: c.config.Prompt.Settings.InitialIndex,
		}
		// Fill the cache with a DB read
	}

	channelCache := c.data[serverId][channelId]
	if len(channelCache.messages) >= c.maxRecords {
		// Remove the oldest record
		channelCache.messages = channelCache.messages[1:]
	}

	channelCache.messages = append(channelCache.messages, record)
}

func (c *MessageCache) Delete(serverId, channelId string) {
	if channelCache, exists := c.data[serverId]; exists {
		delete(channelCache, channelId)
		if len(channelCache) == 0 {
			delete(c.data, serverId)
		}
	}
}

func (c *MessageCache) ChannelStore(serverId, channelId string) (*ChannelStore, bool) {
	if channelCache, exists := c.data[serverId]; exists {
		if channelStore, exists := channelCache[channelId]; exists {
			return channelStore, true
		}
	}
	return nil, false
}
