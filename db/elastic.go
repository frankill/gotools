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
	"github.com/frankill/gotools/query"
	"github.com/olivere/elastic/v7"
)

func EsSimple(host ...string) func(User, Pwd string) (*elastic.Client, error) {

	return func(User, Pwd string) (*elastic.Client, error) {
		return elastic.NewClient(
			elastic.SetSniff(false),
			elastic.SetURL(host...),
			elastic.SetHealthcheck(false),
			elastic.SetBasicAuth(User, Pwd),
		)
	}
}

// EsScriptID returns a function that creates an Elasticsearch stored script
// with a specified ID and parameters.
func EsScriptID(id string) func(data map[string]any) *elastic.Script {
	return func(doc map[string]any) *elastic.Script {
		return elastic.NewScriptStored(id).Params(doc)
	}
}

// EsScript returns a function that creates an Elasticsearch inline script
// with a specified language and script content, and parameters.
func EsScript(lang string, script string) func(data map[string]any) *elastic.Script {
	return func(data map[string]any) *elastic.Script {
		if lang == "" {
			lang = "painless" // Default script language if not provided
		}
		if script == "" {
			log.Panic("script cannot be empty")
		}
		return elastic.NewScript(script).Lang(lang).Params(data)
	}
}

// ElasticBluk represents a bulk request for Elasticsearch
type ElasticBluk[U any] struct {
	Index          string
	OpType         string // Operation type: "index", "create", "update", "delete", "script"
	Id             string
	Routing        string
	Source         U
	DocAsUpsert    bool            // For update operations
	Script         *elastic.Script // For script-based operations
	ScriptAsUpsert bool
	ScriptUpsert   any
}

// 定义一个类型，添加 Index 和 ReturnFields 字段
type ElasticSearchClient[U any] struct {
	Client *elastic.Client
	Query  elastic.Query
}

// 创建 ElasticSearchClient 实例的工厂函数
func NewElasticSearchClient[U any](client *elastic.Client) *ElasticSearchClient[U] {
	return &ElasticSearchClient[U]{
		Client: client,
	}
}

func (es *ElasticSearchClient[U]) Indexs() chan string {

	ch := make(chan string, 100)

	go func() {
		defer close(ch)
		indexs, err := es.Client.IndexNames()
		if err != nil {
			log.Println(err)
			return
		}
		for _, index := range indexs {
			ch <- index
		}
	}()

	return ch
}

func (es *ElasticSearchClient[U]) BulkInsert() func(ch chan ElasticBluk[U]) error {
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

	const maxRetries = 10
	retries := 0

	for {
		request := es.Client.Bulk()
		_, err := request.Add(data...).Refresh("false").Do(context.TODO())
		if err != nil {
			if strings.Contains(err.Error(), "429") {
				retries++
				if retries >= maxRetries {
					return err
				}
				time.Sleep(30 * time.Second)
				continue
			}
			return err
		}
		return nil
	}

}

func createdoc[U any](doc ElasticBluk[U]) elastic.BulkableRequest {

	switch doc.OpType {
	case "index", "create":
		return elastic.NewBulkIndexRequest().Index(doc.Index).OpType(doc.OpType).
			Routing(doc.Routing).Id(doc.Id).UseEasyJSON(true).
			Doc(doc.Source)
	case "update":
		update := elastic.NewBulkUpdateRequest().Index(doc.Index).Id(doc.Id).
			Routing(doc.Routing).Doc(doc.Source).UseEasyJSON(true)
		if doc.DocAsUpsert {
			update.DocAsUpsert(true)
		}
		return update
	case "delete":
		return elastic.NewBulkDeleteRequest().Index(doc.Index).Id(doc.Id).Routing(doc.Routing).UseEasyJSON(true)
	case "script":
		update := elastic.NewBulkUpdateRequest().
			Index(doc.Index).
			Id(doc.Id).
			Routing(doc.Routing).
			Script(doc.Script).
			ScriptedUpsert(doc.ScriptAsUpsert)

		if doc.ScriptAsUpsert {
			if doc.ScriptUpsert != nil {
				update.Upsert(doc.ScriptUpsert)
			} else {
				// Default upsert document if not provided
				update.Upsert(map[string]interface{}{})
			}
		}
		return update
	default:
		log.Panicln("unsupported operation type:", doc.OpType)
		return nil
	}
}

// 查询分片数量并执行 Query
// 用于查询分片数量
// 参数:
//
//	index: 索引名称
//	q: 查询条件 , 支持 string 和 *query.EsQuery, elastic.Query
//
// 返回:
//
//	chan ElasticBluk[U]: 查询结果通道
//	chan error: 错误通道
func (es *ElasticSearchClient[U]) QueryAnyIter(index string, q any) (chan ElasticBluk[U], chan error) {

	stringChan := make(chan ElasticBluk[U], 100)
	errors := make(chan error, 1)

	var wg sync.WaitGroup

	var q_ elastic.Query

	switch v := q.(type) {
	case elastic.Query:
		q_ = v
	case *query.EsQuery:
		q_ = v.Build()
	case string:
		q_ = elastic.NewRawStringQuery(v)
	default:
		log.Panicln("unsupported query type")
	}

	go func() {
		defer close(stringChan)
		defer close(errors)

		slice, err := es.Client.IndexGetSettings(index).Do(context.Background())

		if err != nil {
			errors <- err
			return
		}

		num, _ := strconv.Atoi(slice[index].Settings["index"].(map[string]interface{})["number_of_shards"].(string))

		shardIDs := array.Seq(0, num, 1)

		for _, shardID := range shardIDs {
			wg.Add(1)

			go func(shardID int) {
				defer wg.Done()

				svc := es.Client.Scroll(index).KeepAlive("30m").Size(10000).
					Query(q_).
					Slice(elastic.NewSliceQuery().Id(shardID).Max(num))

				defer svc.Clear(context.Background())

				for {
					res, err := svc.Do(context.Background())

					if err == io.EOF {
						break
					}
					if err != nil {
						errors <- err
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
							errors <- err
							continue
						}

						stringChan <- ElasticBluk[U]{
							Index:   hit.Index,
							Id:      hit.Id,
							Routing: hit.Routing,
							Source:  results,
						}
					}

				}

			}(shardID)
		}

		wg.Wait()

	}()

	// 返回通道和   错误通道
	return stringChan, errors
}
