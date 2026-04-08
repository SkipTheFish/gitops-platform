package store

import (
	"context"

	"gorm.io/gorm"

	"github.com/SkipTheFish/gitops-platform/backend/internal/model"
)

type AppStore struct {
	db *gorm.DB
}

func NewAppStore(db *gorm.DB) *AppStore {
	return &AppStore{db: db}
}

// Create 插入一条 app 记录
func (s *AppStore) Create(ctx context.Context, app *model.App) error {
	return s.db.WithContext(ctx).Create(app).Error
}

// List 查询所有 app
func (s *AppStore) List(ctx context.Context) ([]model.App, error) {
	var apps []model.App
	err := s.db.WithContext(ctx).Order("id desc").Find(&apps).Error
	return apps, err
}

// GetByID 按 id 查询 app
func (s *AppStore) GetByID(ctx context.Context, id int64) (*model.App, error) {
	// 声明承载容器，一个app结构体（model.app格式）
	var app model.App

	// 1. WithContext(ctx)：链路追踪与超时控制
	//    如果前端取消了请求，或者请求处理超时了，ctx 会发出信号。
	// 	  WithContext(ctx) 确保数据库发现连接断开后，立即停止这个耗时查询，释放数据库资源。

	// 2. .First(&app, id)
	// 	翻译成的 SQL：SELECT * FROM apps WHERE id = 123 LIMIT 1;（假设传进来的 id 是 123）。
	// &app (传址)：这是一个“填空”动作。GORM 查到数据后，会通过反射把数据库的字段值一一填入 app 结构体的内存地址里。
	// First 的特性：它是 GORM 的内置方法。如果没有找到记录，它会返回一个特定的错误：gorm.ErrRecordNotFound。
	err := s.db.WithContext(ctx).First(&app, id).Error

	if err != nil {
		return nil, err // 查不到或者数据库挂了
	}
	return &app, nil // 查到了，返回填充好数据的对象地址
}

// GetByName 按名称查询，主要用来做唯一性校验
func (s *AppStore) GetByName(ctx context.Context, name string) (*model.App, error) {
	var app model.App
	err := s.db.WithContext(ctx).Where("name = ?", name).First(&app).Error
	if err != nil {
		return nil, err
	}
	return &app, nil
}

// Update 保存更新后的对象
func (s *AppStore) Update(ctx context.Context, app *model.App) error {
	return s.db.WithContext(ctx).Save(app).Error
}
