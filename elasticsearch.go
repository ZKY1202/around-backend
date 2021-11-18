// 连接elastic search
package main

import (
	"context"

	"github.com/olivere/elastic/v7"
)

const (
	ES_URL = "http://10.162.0.2:9200"
)

func readFromES(query elastic.Query, index string) (*elastic.SearchResult, error) {
	// 与 elastic search 简历链接
	client, err := elastic.NewClient(
		elastic.SetURL(ES_URL),
		elastic.SetBasicAuth("admin", "123456"))
	if err != nil {
		return nil, err
	}

	// 搜索
	searchResult, err := client.Search().
		Index(index).            // search in index "twitter"
		Query(query).            // specify the query
		Pretty(true).            // pretty print request and response JSON
		Do(context.Background()) // execute
	if err != nil {
		return nil, err
	}
	return searchResult, nil
}

func saveToES(i interface{}, index string, id string) error { // 这里 i interface{} 因为为了方便之后 存储不同数据库
	client, err := elastic.NewClient(
		elastic.SetURL(ES_URL),
		elastic.SetBasicAuth("admin", "123456"))
	if err != nil {
		return err
	}

	_, err = client.Index().
		Index(index).
		Id(id).
		BodyJson(i).
		Do(context.Background())
	return err
}
