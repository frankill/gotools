package db

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"strconv"
	"sync"

	"github.com/frankill/gotools/array"
	"github.com/olivere/elastic/v7"
)

// 定义一个类型，添加 Index 和 ReturnFields 字段
type ElasticSearchClient[T any] struct {
	Client *elastic.Client
	Index  string
	Query  elastic.Query
}

// 创建 ElasticSearchClient 实例的工厂函数
func NewElasticSearchClient[T any](client *elastic.Client, index string, query elastic.Query) *ElasticSearchClient[T] {
	return &ElasticSearchClient[T]{
		Client: client,
		Index:  index,
		Query:  query,
	}
}

// 查询分片数量并执行 Query
func (es *ElasticSearchClient[T]) QueryAnyIter() (chan T, error) {

	stringChan := make(chan T, 100)

	var wg sync.WaitGroup

	go func() {
		defer close(stringChan)

		slice, err := es.Client.IndexGetSettings(es.Index).Do(context.Background())

		if err != nil {
			log.Println(err)
			return
		}

		num, _ := strconv.Atoi(slice[es.Index].Settings["index"].(map[string]interface{})["number_of_shards"].(string))

		shardIDs := array.ArraySeq(0, num, 1)

		for _, shardID := range shardIDs {
			wg.Add(1)

			go func(shardID int) {
				defer wg.Done()

				svc := es.Client.Scroll(es.Index).KeepAlive("30m").Size(10000).
					Query(es.Query).
					Slice(elastic.NewSliceQuery().Id(shardID).Max(num))

				defer svc.Clear(context.Background())

				for {
					res, err := svc.Do(context.Background())

					if err == io.EOF {
						break
					}
					if err != nil {
						break
					}
					if res == nil {
						break
					}
					if res.Hits == nil {
						break
					}
					if res.Hits.TotalHits.Value == 0 {
						break
					}
					for _, hit := range res.Hits.Hits {
						var results T
						err := json.Unmarshal(hit.Source, &results)
						if err != nil {
							log.Println(err)
							continue
						}

						stringChan <- results
					}

				}

			}(shardID)
		}

		// 等待所有分片查询完成
		wg.Wait()
	}()

	// 返回通道和 nil 错误
	return stringChan, nil
}
