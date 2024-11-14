package db

import (
	"context"
	"encoding/json"
	"time"

	"github.com/go-redis/redis/v7"
)

type Redis[T any] struct {
	client *redis.Client
	clear  bool
	ctx    context.Context
}

func NewRedisClient[T any](host, pwd string, db int) *Redis[T] {
	rdb := redis.NewClient(&redis.Options{
		Addr:     host,
		Password: pwd,
		DB:       db, // 默认数据库
	})

	return &Redis[T]{
		client: rdb,
	}
}

func (r *Redis[T]) Context(ctx context.Context) *Redis[T] {
	r.ctx = ctx
	return r
}

func (r *Redis[T]) Clear(t bool) *Redis[T] {
	r.clear = t
	return r
}

func (r *Redis[T]) Close() error {
	return r.client.Close()
}

func (r *Redis[T]) LLen(key string) (int64, error) {
	return r.client.LLen(key).Result()
}

func (r *Redis[T]) PushList(key string, data chan T) error {

	if r.ctx == nil {
		r.ctx = context.Background()
	}
	defer func() {
		if r.clear {
			r.Close()
		}

	}()

	for {
		select {

		case <-r.ctx.Done():
			return nil

		case item, ok := <-data:

			if !ok {
				return nil
			}

			itemBytes, err := json.Marshal(item)
			if err != nil {
				return err
			}

			_, err = r.client.RPush(key, itemBytes).Result()
			if err != nil {
				return err
			}
		}
	}

}

func (r *Redis[T]) PushListStr(key string, data chan string) error {

	if r.ctx == nil {
		r.ctx = context.Background()
	}
	defer func() {
		if r.clear {
			r.Close()
		}

	}()

	for {
		select {

		case <-r.ctx.Done():
			return nil

		case item, ok := <-data:

			if !ok {
				return nil
			}

			itemBytes := []byte(item)

			_, err := r.client.RPush(key, itemBytes).Result()
			if err != nil {
				return err
			}
		}
	}

}

func (r *Redis[T]) PopList(key string, num int) (chan T, chan error) {

	if r.ctx == nil {
		r.ctx = context.Background()
	}

	ch := make(chan T, 100)
	errs := make(chan error, 3)

	go func() {

		defer close(ch)
		defer close(errs)

		defer func() {
			if r.clear {
				r.Close()
			}
		}()

		if num == 0 {
			return
		}

		n := 0

		for {

			if num > 0 && n >= num {
				return
			}

			select {

			case <-r.ctx.Done():
				return

			default:
				var item T
				value, err := r.client.LPop(key).Result()

				if err == redis.Nil {
					time.Sleep(time.Second * 10)
					continue
				}

				if err != nil {
					errs <- err
					return
				}

				if err := json.Unmarshal([]byte(value), &item); err != nil {
					errs <- err
					return
				}

				ch <- item

				n++
			}
		}

	}()

	return ch, errs
}

func (r *Redis[T]) PopListStr(key string, num int) (chan string, chan error) {

	if r.ctx == nil {
		r.ctx = context.Background()
	}

	ch := make(chan string, 100)
	errs := make(chan error, 3)

	go func() {

		defer close(ch)
		defer close(errs)

		defer func() {
			if r.clear {
				r.Close()
			}
		}()

		if num == 0 {
			return
		}

		n := 0

		for {

			if num > 0 && n >= num {
				return
			}

			select {

			case <-r.ctx.Done():
				return

			default:

				value, err := r.client.LPop(key).Result()

				if err == redis.Nil {
					time.Sleep(time.Second * 10)
					continue
				}

				if err != nil {
					errs <- err
					return
				}

				ch <- value

				n++
			}
		}

	}()

	return ch, errs
}
