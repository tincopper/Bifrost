/*
Copyright [2018] [jc3wish]

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
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