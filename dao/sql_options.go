package dao

import "github.com/guzhongzhi/gmicro/logger"

type SQLDAOOptions struct {
	Logger    logger.SuperLogger
	GetSQL    string
	FindSQL   string
	DeleteSQL string
	UpdateSQL string
	InsertSQL string
}
