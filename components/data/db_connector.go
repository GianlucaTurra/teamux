package data

import (
	"context"

	"gorm.io/gorm"
)

type Connector struct {
	DB  *gorm.DB
	Ctx context.Context
}
