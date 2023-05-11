package redis

import (
	"context"
	cont "github.com/daida459031925/common/context"
	err "github.com/daida459031925/common/error"
	"github.com/daida459031925/common/file"
	"github.com/daida459031925/common/fmt"
	reqUtil "github.com/daida459031925/common/util/requestUtil"
	"github.com/go-redis/redis/v8"
	"net/http"
	sysTime "time"
)

type Redis struct {
	*redis.Client
}

type config struct {
	Redis struct {
		Host     string `yaml:"host"`     //地址127.0.0.1:3306
		Port     string `yaml:"port"`     //地址127.0.0.1:3306
		Password string `yaml:"password"` //密码********
		Db       int    `yaml:"db"`       //redis数据库编号
		PoolSize int    `yaml:"poolSize"` //连接池大小
	} `yaml:"redis"`
}

const (
	STRING     Type = "string"
	LIST       Type = "list"
	SET        Type = "set"
	ZSET       Type = "zset"
	SORTED_SET Type = "sorted set"
	HASH       Type = "hash"
)

type Type string

func NewRedisConfig(filePath string) (*config, error) {
	return file.NewConfig[config](filePath)
}

func (c *config) NewRedis() *Redis {
	redisDb := redis.Options{
		Addr:     fmt.Sprintf("%s:%s", c.Redis.Host, c.Redis.Port),
		Password: c.Redis.Password,
		DB:       c.Redis.Db,
	}

	var poolSize = c.Redis.PoolSize
	if poolSize != 0 {
		redisDb.PoolSize = poolSize
	}

	redisClient := redis.NewClient(&redisDb)

	return &Redis{
		redisClient,
	}
}

// WithContext 设置前置方法
func (r *Redis) WithContext(key string, expectedType Type, f func(context.Context) error) error {
	return r.WithRequestContext(key, expectedType, f)
}

// WithRequestContext 配置一个前置方法使得所有的redis都可以拿到context
func (r *Redis) WithRequestContext(key string, expectedType Type, f func(context.Context) error, reqs ...*http.Request) error {
	ctx, cancel := cont.GetBaseContext().GetRequestContext(reqUtil.GetRequest(reqs...))
	defer cancel()
	if key != "" && expectedType != "" {
		ok, e := r.CheckKeyExistAndTypeMatch(key, expectedType, reqs...)
		if ok {
			return e
		}
	}
	e := f(ctx)

	if e != nil {
		return e
	}

	return nil
}

// Del 删除指定key的数据保
func (r *Redis) Del(keys []string, reqs ...*http.Request) error {
	return r.WithRequestContext("", "", func(ctx context.Context) error {
		return r.Client.Del(ctx, keys...).Err()
	}, reqs...)
}

// Close 对redis整个链接进行关闭，在项目关闭时候进行使用
func (r *Redis) Close() error {
	return r.Client.Close()
}

// Ping 测试连接是否可以正常使用
func (r *Redis) Ping(reqs ...*http.Request) error {
	return r.WithRequestContext("", "", func(ctx context.Context) error {
		return r.Client.Ping(ctx).Err()
	}, reqs...)
}

// Exists 用于检测当前redis中是否存在对应key的数据
func (r *Redis) Exists(keys []string, reqs ...*http.Request) (int64, error) {
	var value int64
	e := r.WithRequestContext("", "", func(ctx context.Context) error {
		var e error
		value, e = r.Client.Exists(ctx, keys...).Result()
		return e
	}, reqs...)
	if e != nil {
		return 0, e
	}
	return value, nil
}

// GetType 用于获取当前key是什么类型
func (r *Redis) GetType(key string, reqs ...*http.Request) (string, error) {
	var value string
	e := r.WithRequestContext("", "", func(ctx context.Context) error {
		var e error
		value, e = r.Client.Type(ctx, key).Result()
		return e
	}, reqs...)
	if e != nil {
		return "", e
	}
	return value, nil
}

// CheckKeyExistAndTypeMatch 检查指定键是否存在并且类型匹配
func (r *Redis) CheckKeyExistAndTypeMatch(key string, expectedType Type, reqs ...*http.Request) (bool, error) {
	value := false
	e := r.WithRequestContext("", "", func(ctx context.Context) error {
		// 检查键是否存在
		res, e := r.Client.Exists(ctx, key).Result()
		if e != nil {
			return err.NewSprintf("redis无法获取到KEY: %s的值", key)
		}
		if res == 0 {
			return err.NewSprintf("redis中不存在KEY: %s当前值", key)
		}

		// 检查键值类型是否匹配
		actualType, e := r.Client.Type(ctx, key).Result()
		if e != nil {
			return err.NewSprintf("redis获取KEY: %s的类型失败", key)
		}

		value = actualType != string(expectedType)
		if value {
			return err.NewSprintf("redis获取KEY: %s的类型%s与需要匹配的类型%s不一致", key, actualType, expectedType)
		}
		return nil
	}, reqs...)

	return value, e
}

/* String 相关方法*/

// Get 方法，获取指定键的值
func (r *Redis) Get(key string, reqs ...*http.Request) (string, error) {
	var value string
	e := r.WithRequestContext(key, STRING, func(ctx context.Context) error {
		var e error
		value, e = r.Client.Get(ctx, key).Result()
		return e
	}, reqs...)
	if e != nil {
		return "", e
	}
	return value, nil
}

// Set 方法，设置指定键的值，并指定过期时间
func (r *Redis) Set(key, value string, expiration sysTime.Duration, reqs ...*http.Request) error {
	return r.WithRequestContext(key, STRING, func(ctx context.Context) error {
		return r.Client.Set(ctx, key, value, expiration).Err()
	}, reqs...)
}

// Incr 方法，将指定键的值增加 1，并返回自增后的结果存储的是整数类型的就可以进行加减
func (r *Redis) Incr(key string, reqs ...*http.Request) error {
	return r.WithRequestContext("", "", func(ctx context.Context) error {
		return r.Client.Incr(ctx, key).Err()
	}, reqs...)
}

// Decr 方法，将指定键的值减少 1，并返回自减后的结果 的是整数类型的就可以进行加减
func (r *Redis) Decr(key string, reqs ...*http.Request) error {
	return r.WithRequestContext("", "", func(ctx context.Context) error {
		return r.Client.Decr(ctx, key).Err()
	}, reqs...)
}

/* List  相关方法*/

// GetListAll 获取指定key的数据返还数据为哈希：map[string]string
func (r *Redis) GetListAll(key string, reqs ...*http.Request) ([]string, error) {
	var value []string
	e := r.WithRequestContext(key, LIST, func(ctx context.Context) error {
		var e error
		value, e = r.Client.LRange(ctx, key, 0, -1).Result()
		return e
	}, reqs...)
	if e != nil {
		return nil, e
	}
	return value, nil
}

// SetListLPush 将一组元素从列表左侧插入
func (r *Redis) SetListLPush(key string, value []any, expiration sysTime.Duration, reqs ...*http.Request) error {
	return r.WithRequestContext(key, LIST, func(ctx context.Context) error {
		e := r.Client.LPush(ctx, key, value...).Err()
		if e != nil {
			return e
		}
		return r.Client.Expire(ctx, key, expiration).Err()
	}, reqs...)
}

// SetListRPush 将一组元素从列表右侧插入
func (r *Redis) SetListRPush(key string, value []any, expiration sysTime.Duration, reqs ...*http.Request) error {
	return r.WithRequestContext(key, LIST, func(ctx context.Context) error {
		e := r.Client.RPush(ctx, key, value...).Err()
		if e != nil {
			return e
		}
		return r.Client.Expire(ctx, key, expiration).Err()
	}, reqs...)
}

// GetListLPop 方法，从列表左侧弹出一个元素并删除它
func (r *Redis) GetListLPop(key string, reqs ...*http.Request) (string, error) {
	var value string
	e := r.WithRequestContext(key, LIST, func(ctx context.Context) error {
		var e error
		value, e = r.Client.LPop(ctx, key).Result()
		return e
	}, reqs...)
	if e != nil {
		return "", e
	}
	return value, nil
}

// GetListRPop 方法，从列表右侧弹出一个元素并删除它
func (r *Redis) GetListRPop(key string, reqs ...*http.Request) (string, error) {
	var value string
	e := r.WithRequestContext(key, LIST, func(ctx context.Context) error {
		var e error
		value, e = r.Client.RPop(ctx, key).Result()
		return e
	}, reqs...)
	if e != nil {
		return "", e
	}
	return value, nil
}

/* Set   相关方法*/

// GetSet 向集合中添加一组成员
func (r *Redis) GetSet(key string, reqs ...*http.Request) ([]string, error) {
	var value []string
	e := r.WithRequestContext(key, SET, func(ctx context.Context) error {
		var e error
		value, e = r.Client.SMembers(ctx, key).Result()
		return e
	}, reqs...)
	if e != nil {
		return nil, e
	}
	return value, nil
}

// SetSet 获取集合中所有成员
func (r *Redis) SetSet(key string, value []any, expiration sysTime.Duration, reqs ...*http.Request) error {
	return r.WithRequestContext(key, SET, func(ctx context.Context) error {
		e := r.Client.SAdd(ctx, key, value...).Err()
		if e != nil {
			return e
		}
		return r.Client.Expire(ctx, key, expiration).Err()
	}, reqs...)
}

/* Hash  相关方法*/

// GetHash 获取指定key的数据返还数据为哈希：struct
func (r *Redis) GetHash(key string, reqs ...*http.Request) (map[string]string, error) {
	var value map[string]string
	e := r.WithRequestContext(key, HASH, func(ctx context.Context) error {
		var e error
		value, e = r.Client.HGetAll(ctx, key).Result()
		return e
	}, reqs...)
	if e != nil {
		return nil, e
	}
	return value, nil
}

// GetHashHGet 方法，获取哈希表中指定字段的值
func (r *Redis) GetHashHGet(key, field string, reqs ...*http.Request) (string, error) {
	var value string
	e := r.WithRequestContext(key, HASH, func(ctx context.Context) error {
		var e error
		value, e = r.Client.HGet(ctx, key, field).Result()
		return e
	}, reqs...)
	if e != nil {
		return "", e
	}
	return value, nil
}

// SetHash 保存指定key的数据保存数据为哈希：
func (r *Redis) SetHash(key string, value []any, expiration sysTime.Duration, reqs ...*http.Request) error {
	return r.WithRequestContext(key, HASH, func(ctx context.Context) error {
		e := r.Client.HSet(ctx, key, value...).Err()
		if e != nil {
			return e
		}
		return r.Client.Expire(ctx, key, expiration).Err()
	}, reqs...)
}

// HDel 方法，删除哈希表中指定字段
func (r *Redis) HDel(key string, fields []string, reqs ...*http.Request) error {
	return r.WithRequestContext(key, HASH, func(ctx context.Context) error {
		return r.Client.HDel(ctx, key, fields...).Err()
	}, reqs...)
}

/* Sorted set  相关方法*/

// GetLinkedSet 获取指定key的数据返还数据为哈希：map[string]string
func (r *Redis) GetLinkedSet(key string, reqs ...*http.Request) ([]string, error) {
	var value []string
	e := r.WithRequestContext(key, ZSET, func(ctx context.Context) error {
		var e error
		value, e = r.Client.ZRange(ctx, key, 0, -1).Result()
		return e
	}, reqs...)
	if e != nil {
		return nil, e
	}
	return value, nil
}

// SetLinkedSet 保存指定key的数据保存数据为哈希：
func (r *Redis) SetLinkedSet(key string, value []any, expiration sysTime.Duration, reqs ...*http.Request) error {
	//SORTED_SET
	return r.WithRequestContext(key, ZSET, func(ctx context.Context) error {
		var e error
		for i := 0; i < len(value); i++ {
			rz := &redis.Z{Score: float64(i + 1.0), Member: value[i]}
			e = r.Client.ZAdd(ctx, key, rz).Err()
			if e != nil {
				var str = make([]string, 1)
				str[0] = key
				r.Del(str, reqs...)
				break
			}
		}
		if e != nil {
			return e
		}

		return r.Client.Expire(ctx, key, expiration).Err()
	}, reqs...)
}
