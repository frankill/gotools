package db

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/frankill/gotools/array"
	"github.com/olivere/elastic/v7"
)

type ElasticBluk[U any] struct {
	Index   string
	Type    string
	Id      string
	Routing string
	Source  U
}

// 定义一个类型，添加 Index 和 ReturnFields 字段
type ElasticSearchClient[U any] struct {
	Client *elastic.Client
	Index  string
	Query  elastic.Query
}

// 创建 ElasticSearchClient 实例的工厂函数
func NewElasticSearchClient[U any](client *elastic.Client) *ElasticSearchClient[U] {
	return &ElasticSearchClient[U]{
		Client: client,
	}
}

func (es *ElasticSearchClient[U]) BulkInsert(index string, ctype string) func(ch chan ElasticBluk[U]) error {
	return func(ch chan ElasticBluk[U]) error {
		bulkSize := 3000
		bulkData := make([]elastic.BulkableRequest, 0, bulkSize)

		for data := range ch {
			bulkData = append(bulkData, createdoc(data))
			if len(bulkData) >= bulkSize {
				if err := es.sendBulk(bulkData); err != nil {
					return err
				}
				bulkData = bulkData[:0]
			}
		}

		// Send remaining data
		if len(bulkData) > 0 {
			if err := es.sendBulk(bulkData); err != nil {
				return err
			}
		}

		return nil
	}
}

// 假设sendBulk是一个发送批次数据到Elasticsearch的函数
func (es *ElasticSearchClient[U]) sendBulk(data []elastic.BulkableRequest) error {

	request := es.Client.Bulk()
	_, err := request.Add(data...).Refresh("false").Do(context.TODO())

	if err != nil && strings.Index(err.Error(), "429") > 0 {
		time.Sleep(30 * time.Second)
		es.sendBulk(data)
		return nil
	}

	return err

}

func createdoc[U any](doc ElasticBluk[U]) elastic.BulkableRequest {

	return elastic.NewBulkIndexRequest().Index(doc.Index).Type(doc.Type).
		Routing(doc.Routing).Id(doc.Id).UseEasyJSON(true).
		Doc(doc.Source)

}

// 查询分片数量并执行 Query
func (es *ElasticSearchClient[U]) QueryAnyIter(index string, query elastic.Query) (chan ElasticBluk[U], error) {

	stringChan := make(chan ElasticBluk[U], 100)

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
						var results U
						err := json.Unmarshal(hit.Source, &results)
						if err != nil {
							log.Println(err)
							continue
						}

						stringChan <- ElasticBluk[U]{
							Index:   hit.Index,
							Type:    hit.Type,
							Id:      hit.Id,
							Routing: hit.Routing,
							Source:  results,
						}
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
