package sql

import (
	"github.com/daida459031925/common/fmt"
	myRedis "github.com/daida459031925/common/sql/redis"
	"gorm.io/gorm"
)

type Db interface {
	Connect() (*BaseDb, error)
	ConnectCheck(checkFilePath string) (*BaseDb, error)
}

type BaseDb struct {
	Db    Db
	Gdb   gorm.DB
	Redis *myRedis.Redis
}

func (b *BaseDb) Connect() (*BaseDb, error) {
	r, e := b.Db.Connect()
	if e != nil {
		return nil, e
	}
	// 其他通用逻辑
	return r, nil
}

func (b *BaseDb) ConnectCheck(checkFilePath string) (*BaseDb, error) {
	r, e := b.Db.ConnectCheck(checkFilePath)
	if e != nil {
		return nil, e
	}
	fmt.Printlnf("2")
	// 其他通用逻辑
	return r, nil
}

func (b *BaseDb) Close() {
	db, e := b.Gdb.DB()
	if e != nil {
		fmt.Println(e)
	} else {
		e = db.Close()
		if e != nil {
			fmt.Println(e)
		}
	}

	r := b.Redis
	if r == nil {
		return
	}
	client := r.Client
	if client == nil {
		fmt.Println("Redis.Client 为空")
		return
	}
	e = client.Close()
	if e != nil {
		fmt.Println(e)
	}
}

//
//// Database 定义通用的DB接口
//type Database interface {
//	Connect() (Database, error)
//	ConnectCheck(filePath string) (DB, error)
//	Close() error
//	//Find(dest any, cons ...any) error
//	//Create(value any) error
//	//Update(value any) error
//	//Delete(value any) error
//	//Exec(sql string, values ...any) error
//}
//
//type DB struct {
//	GormDB      *gorm.DB
//	RedisClient *redis.Client
//}
//
//func (db *DB) Connect() error {
//	d, e := db.GormDB.DB()
//	if e != nil {
//		return e
//	}
//	d.Close()
//	if err := db.RedisClient.Close(); err != nil {
//		return err
//	}
//	return nil
//}
//
//// Close 关闭数据库连接和Redis连接
//func (db *DB) Close() error {
//	d, e := db.GormDB.DB()
//	if e != nil {
//		return e
//	}
//	d.Close()
//	if err := db.RedisClient.Close(); err != nil {
//		return err
//	}
//	return nil
//}

//func (db *DB) Find(dest interface{}, conds ...interface{}) error {
//	// 从Redis中查询数据
//	if err := db.redisClient.Get("cache_key").Scan(dest); err == nil {
//		return nil
//	}
//
//	// 从数据库中查询数据
//	if err := db.gormDB.Where(conds...).Find(dest).Error; err != nil {
//		return err
//	}
//
//	// 将查询结果存入Redis中
//	if err := db.redisClient.Set("cache_key", dest, time.Minute).Err(); err != nil {
//		return err
//	}
//
//	return nil
//}
//
//func (db *DB) Create(value interface{}) error {
//	// 插入数据到数据库中
//	if err := db.gormDB.Create(value).Error; err != nil {
//		return err
//	}
//
//	// 如果Redis中有缓存，删除缓存
//	if err := db.redisClient.Del("cache_key").Err(); err != nil {
//		return err
//	}
//
//	return nil
//}
//
//func (db *DB) Update(value interface{}) error {
//	// 更新数据到数据库中
//	if err := db.gormDB.Save(value).Error; err != nil {
//		return err
//	}
//
//	// 如果Redis中有缓存，删除缓存
//	if err := db.redisClient.Del("cache_key").Err(); err != nil {
//		return err
//	}
//
//	return nil
//}
//
//func (db *DB) Delete(value interface{}) error {
//	// 删除数据从数据库中
//	if err := db.gormDB.Delete(value).Error; err != nil {
//		return err
//	}
//
//	// 如果Redis中有缓存，删除缓存
//	if err := db.redisClient.Del("cache_key").Err(); err != nil {
//		return err
//	}
//
//	return nil
//}
//
//func (db *DB) Exec(sql string, values ...interface{}) error {
//	// 执行自定义的SQL语句
//	if err := db.gormDB.Exec(sql, values...).Error; err != nil {
//		return err
//	}
//
//	// 如果Redis中有缓存，删除缓存
//	if err := db.redisClient.Del("cache_key").Err(); err != nil {
//		return err
//	}
//
//	return nil
//}
