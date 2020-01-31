package datamodles

type FlowIndex struct {
	Title string
	DbName string
	Schema string
	TableName string
	ChannelId string
}

func NewFlowIndex(title, dbName, schema, tableName, channelId string) *FlowIndex {
	return &FlowIndex{
		Title: 	   title,
		DbName:    dbName,
		Schema:    schema,
		TableName: tableName,
		ChannelId: channelId,
	}
}