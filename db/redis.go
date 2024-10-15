package db

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/frankill/gotools"
	"github.com/go-redis/redis/v7"
)

type Redis[T any] struct {
	client *redis.Client
	ctx    context.Context
	NUM    int64
}

func NewRedisClient[T any](host, pwd string, db int) *Redis[T] {
	rdb := redis.NewClient(&redis.Options{
		Addr:     host,
		Password: pwd,
		DB:       db, // 默认数据库
	})

	ctx := gotools.SysStop()

	return &Redis[T]{
		client: rdb,
		ctx:    ctx,
		NUM:    100,
	}
}

func (r *Redis[T]) PopList(key string) (chan T, chan error) {

	ch := make(chan T)
	errs := make(chan error, 1)

	go func() {
		defer close(ch)
		defer close(errs)
		defer r.client.Close()

		pipe := r.client.Pipeline()
		for i := 0; i < int(r.NUM); i++ {
			pipe.LPop(key)
		}

		for {

			if _, err := r.client.Ping().Result(); err != nil {
				errs <- err
				return
			}

			select {

			case <-r.ctx.Done():
				return

			default:
				num, _ := r.client.LLen(key).Result()
				if num == 0 {
					time.Sleep(100 * time.Second)
					continue
				}
				if num < r.NUM {
					pipe = r.client.Pipeline()
					for i := 0; i < int(num); i++ {
						pipe.LPop(key)
					}

				}

				res, err := pipe.Exec()
				if err != nil {
					errs <- err
					return
				}

				for _, v := range res {
					var tmp T
					err = json.Unmarshal([]byte(v.String()[len(key)+7:]), tmp)

					if err != nil {
						errs <- errors.New("json unmarshal error:" + v.String())
						continue
					}
					ch <- tmp
				}

			}

		}
	}()

	return ch, errs
}
