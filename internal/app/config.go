package app

import (
	"github.com/core-go/core"
	"github.com/core-go/core/header/echo"
	"github.com/core-go/core/server"
	l "github.com/core-go/log/echo"
	"github.com/core-go/log/zap"
)

type Config struct {
	Server     server.ServerConfig `mapstructure:"server"`
	Mongo      MongoConfig         `mapstructure:"mongo"`
	Log        log.Config          `mapstructure:"log"`
	Response   header.Config       `mapstructure:"response"`
	MiddleWare l.LogConfig         `mapstructure:"middleware"`
	Action     *core.ActionConfig  `mapstructure:"action"`
}

type MongoConfig struct {
	Uri      string `yaml:"uri" mapstructure:"uri" json:"uri,omitempty" gorm:"column:uri" bson:"uri,omitempty" dynamodbav:"uri,omitempty" firestore:"uri,omitempty"`
	Database string `yaml:"database" mapstructure:"database" json:"database,omitempty" gorm:"column:database" bson:"database,omitempty" dynamodbav:"database,omitempty" firestore:"database,omitempty"`
}
