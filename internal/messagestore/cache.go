package messagestore

type MessageCache struct {
	data       map[string]map[string][]MessageRecord
	maxRecords int
}

func NewMessageCache(maxRecords int) *MessageCache {
	return &MessageCache{
		data:       make(map[string]map[string][]MessageRecord),
		maxRecords: maxRecords,
	}
}

func (c *MessageCache) Put(serverId, channelId string, record MessageRecord) {
	if _, exists := c.data[serverId]; !exists {
		c.data[serverId] = make(map[string][]MessageRecord)
	}

	if _, exists := c.data[serverId][channelId]; !exists {
		c.data[serverId][channelId] = []MessageRecord{}
	}

	channelCache := c.data[serverId][channelId]
	if len(channelCache) >= c.maxRecords {
		// Remove the oldest record
		channelCache = channelCache[1:]
	} else {
		// Fill the cache with a DB read
		//   For the first 20 messages in any given channel, there will be about 20 reads on empty data
		//   C'est la vie
	}
	channelCache = append(channelCache, record)
	c.data[serverId][channelId] = channelCache
}

func (c *MessageCache) Delete(serverId, channelId string) {
	if channelCache, exists := c.data[serverId]; exists {
		delete(channelCache, channelId)
		if len(channelCache) == 0 {
			delete(c.data, serverId)
		}
	}
}

func (c *MessageCache) Get(serverId, channelId string) ([]MessageRecord, bool) {
	if channelCache, exists := c.data[serverId]; exists {
		if records, exists := channelCache[channelId]; exists {
			return records, true
		}
	}
	return nil, false
}
