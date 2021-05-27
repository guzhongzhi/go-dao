package dao

import "github.com/pinguo-icc/kratos-library/logger"

type SQLDAOOptions struct {
	Logger    logger.Logger
	GetSQL    string
	FindSQL   string
	DeleteSQL string
	UpdateSQL string
	InsertSQL string
}
