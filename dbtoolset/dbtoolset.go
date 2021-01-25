package dbtoolset

import (
	"context"
	"errors"

	"github.com/go-redis/redis/v8"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	"github.com/jiuzhou-zhao/go-fundamental/interfaces"
)

type DBToolset struct {
	ctx          context.Context
	cfg          *DBConfig
	logger       interfaces.Logger
	defaultRedis *redis.Client
	redisMap     map[string]*redis.Client
	defaultMySQL *xorm.Engine
	mySQLMap     map[string]*xorm.Engine
}

func NewDBToolset(ctx context.Context, cfg *DBConfig, logger interfaces.Logger) (*DBToolset, error) {
	if logger == nil {
		logger = &interfaces.ConsoleLogger{}
	}
	if cfg == nil {
		logger.Record(ctx, interfaces.LogLevelFatal, "no config")
		return nil, errors.New("no config")
	}
	toolset := &DBToolset{
		ctx:      ctx,
		cfg:      cfg,
		logger:   logger,
		redisMap: make(map[string]*redis.Client),
		mySQLMap: make(map[string]*xorm.Engine),
	}
	err := toolset.allRedisInit()
	if err != nil {
		logger.Recordf(ctx, interfaces.LogLevelFatal, "init redis failed: %v", err)
		return nil, err
	}
	err = toolset.allMySQLInit()
	if err != nil {
		logger.Recordf(ctx, interfaces.LogLevelFatal, "init mysql failed: %v", err)
		return nil, err
	}
	return toolset, nil
}

func (toolset *DBToolset) redisInit(cfg *RedisConfig) (redisCli *redis.Client, err error) {
	options, err := redis.ParseURL(cfg.DSN)
	if err != nil {
		toolset.logger.Recordf(toolset.ctx, interfaces.LogLevelError, "init redis failed: %v", err)
		return
	}
	redisCli = redis.NewClient(options)
	return
}

func (toolset *DBToolset) allRedisInit() error {
	for name, redisCfg := range toolset.cfg.Redis {
		redisCli, err := toolset.redisInit(&redisCfg)
		if err != nil {
			return err
		}
		toolset.redisMap[name] = redisCli
		if name == "" || name == "default" || name == "def" {
			toolset.defaultRedis = redisCli
		}
	}
	if toolset.defaultRedis == nil {
		if len(toolset.redisMap) == 1 {
			for _, client := range toolset.redisMap {
				toolset.defaultRedis = client
			}
		}
	}
	return nil
}

func (toolset *DBToolset) allMySQLInit() error {
	for name, cfg := range toolset.cfg.MySQL {
		mySQLCli, err := xorm.NewEngine("mysql", cfg.DSN)
		if err != nil {
			return err
		}
		if cfg.ShowSQL {
			mySQLCli.ShowSQL(true)
		}
		toolset.mySQLMap[name] = mySQLCli
		if name == "" || name == "default" || name == "def" {
			toolset.defaultMySQL = mySQLCli
		}
	}
	if toolset.defaultMySQL == nil {
		if len(toolset.mySQLMap) == 1 {
			for _, client := range toolset.mySQLMap {
				toolset.defaultMySQL = client
			}
		}
	}
	return nil
}

func (toolset *DBToolset) GetRedis() *redis.Client {
	if toolset.defaultRedis == nil {
		toolset.logger.Record(toolset.ctx, interfaces.LogLevelFatal, "no default redis")
	}
	return toolset.defaultRedis
}

func (toolset *DBToolset) GetRedisByName(name string) *redis.Client {
	return toolset.redisMap[name]
}

func (toolset *DBToolset) GetMySQL() *xorm.Engine {
	if toolset.defaultMySQL == nil {
		toolset.logger.Record(toolset.ctx, interfaces.LogLevelFatal, "no default mysql")
	}
	return toolset.defaultMySQL
}

func (toolset *DBToolset) GetMySQLByName(name string) *xorm.Engine {
	return toolset.mySQLMap[name]
}
