package office

import (
	"github.com/daida459031925/common/fmt"
	"github.com/daida459031925/common/sql/redis"
	"github.com/daida459031925/common/time"
	"testing"
)

func BaseRedis() (*redis.Redis, error) {
	r, e := redis.NewRedisConfig("..\\redis.yml")
	if e != nil {
		fmt.Println(e)
		return nil, e
	}

	myRedis := r.NewRedis()

	e = myRedis.Ping(nil)
	if e != nil {
		fmt.Println(e)
		return nil, e
	}

	return myRedis, nil
}

func TestRedis(t *testing.T) {

	myRedis, e := BaseRedis()
	if e != nil {
		return
	}
	defer myRedis.Close()

	testString(myRedis, "123", "abc", 5)
	testString(myRedis, "123", "12", 5) //只有整数才能进行加减
	testString(myRedis, "123", "12.2", 5)

	fmt.Println("String测试结束")

	a := make([]any, 2)
	a[0] = "1"
	a[1] = "2"

	testList(myRedis, "123", a, 5)
	fmt.Println("List测试结束")

	testSet(myRedis, "123", a, 5)
	fmt.Println(myRedis.Exists([]string{"123"}))
	fmt.Println(myRedis.GetType("123"))
	fmt.Println("set测试结束")

	testLinkedSet(myRedis, "123", a, 5)
	fmt.Println("LinkedSet测试结束")

	testHash(myRedis, "123", a, 5)
	fmt.Println("LinkedSet测试结束")

}

func TestJiaoChaRedis(t *testing.T) {
	myRedis, e := BaseRedis()
	if e != nil {
		return
	}
	defer myRedis.Close()
	key := "123"

	myRedis.Set(key, "123", time.GetSecond(10))
	fmt.Println(myRedis.GetType(key))
	fmt.Println(myRedis.GetSet(key))
	fmt.Println(myRedis.SetHash(key, []any{"1", 2}, time.GetSecond(10)))
	myRedis.Del([]string{key})
	fmt.Println()

	myRedis.SetListLPush(key, []any{"123"}, time.GetSecond(10))
	fmt.Println(myRedis.GetType(key))
	fmt.Println(myRedis.Get(key))
	myRedis.Del([]string{key})
	fmt.Println()

	myRedis.SetSet(key, []any{"123"}, time.GetSecond(10))
	fmt.Println(myRedis.GetType(key))
	fmt.Println(myRedis.Get(key))
	myRedis.Del([]string{key})
	fmt.Println()

	myRedis.SetLinkedSet(key, []any{"123"}, time.GetSecond(10))
	fmt.Println(myRedis.GetType(key))
	fmt.Println(myRedis.Get(key))
	myRedis.Del([]string{key})
	fmt.Println()

	myRedis.SetHash(key, []any{"123", 123}, time.GetSecond(10))
	fmt.Println(myRedis.GetType(key))
	fmt.Println(myRedis.Get(key))
	myRedis.Del([]string{key})
}

func testString(redis *redis.Redis, key, val string, t int) {
	defer redis.Del([]string{key})
	redis.Set(key, val, time.GetSecond(t))
	value, e := redis.Get(key)
	if e != nil {
		fmt.Println("error:", e)
	}
	fmt.Println(value)

	e = redis.Incr(key)
	if e != nil {
		fmt.Println("error:", e)
	}

	value, e = redis.Get(key)
	if e != nil {
		fmt.Println("error:", e)
	}
	fmt.Println(value)

	e = redis.Decr(key)
	if e != nil {
		fmt.Println("error:", e)
	}

	value, e = redis.Get(key)
	if e != nil {
		fmt.Println("error:", e)
	}
	fmt.Println(value)

	keys := make([]string, 1)
	keys[0] = key
	e = redis.Del(keys)
	if e != nil {
		fmt.Println("error:", e)
	}

	value, e = redis.Get(key)
	if e != nil {
		fmt.Println("error:", e)
	}
}

func testList(redis *redis.Redis, key string, val []any, t int) {
	defer redis.Del([]string{key})
	a := make([]any, 1)
	a[0] = "start"
	e := redis.SetListLPush(key, a, time.GetHour(t))
	if e != nil {
		fmt.Println("error:", e)
	}

	value, e := redis.GetListAll(key)
	if e != nil {
		fmt.Println("error:", e)
	}
	fmt.Println(value)

	e = redis.SetListLPush(key, val, time.GetHour(t))
	if e != nil {
		fmt.Println("error:", e)
	}
	value, e = redis.GetListAll(key)
	if e != nil {
		fmt.Println("error:", e)
	}
	fmt.Println(value)

	e = redis.SetListRPush(key, val, time.GetHour(t))
	if e != nil {
		fmt.Println("error:", e)
	}
	value, e = redis.GetListAll(key)
	if e != nil {
		fmt.Println("error:", e)
	}
	fmt.Println(value)

	r, e := redis.GetListLPop(key)
	if e != nil {
		fmt.Println("error:", e)
	}
	fmt.Println(r)

	value, e = redis.GetListAll(key)
	if e != nil {
		fmt.Println("error:", e)
	}
	fmt.Println(value)

	r, e = redis.GetListRPop(key)
	if e != nil {
		fmt.Println("error:", e)
	}
	fmt.Println(r)

	value, e = redis.GetListAll(key)
	if e != nil {
		fmt.Println("error:", e)
	}
	fmt.Println(value)
}

func testSet(redis *redis.Redis, key string, val []any, t int) {
	defer redis.Del([]string{key})
	redis.SetSet(key, val, time.GetSecond(t))
	value, e := redis.GetSet(key)
	if e != nil {
		fmt.Println("error:", e)
	}
	fmt.Println(value)

}

func testLinkedSet(redis *redis.Redis, key string, val []any, t int) {
	defer redis.Del([]string{key})
	e := redis.SetLinkedSet(key, val, time.GetSecond(t))
	if e != nil {
		fmt.Println("error:", e)
	}
	value, e := redis.GetLinkedSet(key)
	if e != nil {
		fmt.Println("error:", e)
	}
	fmt.Println(value)

}

func testHash(redis *redis.Redis, key string, val []any, t int) {
	defer redis.Del([]string{key})
	value, e := redis.GetHash(key)
	if e != nil {
		fmt.Println("error:", e)
	}
	fmt.Println(value)

	redis.SetHash(key, val, time.GetSecond(t))
	v, e := redis.GetHash(key)
	if e != nil {
		fmt.Println("error:", e)
	}
	fmt.Println(v)

	vl, e := redis.GetHashHGet(key, "1")
	if e != nil {
		fmt.Println(e)
	}

	fmt.Println(vl)

	e = redis.HDel(key, []string{"1"})
	if e != nil {
		fmt.Println(vl)
	}

	vl, e = redis.GetHashHGet(key, "1")
	if e != nil {
		fmt.Println(e)
	}

	fmt.Println(vl)
}
