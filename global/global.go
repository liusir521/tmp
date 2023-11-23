package global

import (
	"github.com/patrickmn/go-cache"
	"time"
)

var (
	// 司机和顾客的对象缓存池
	GlobalUser   = cache.New(12*time.Hour, 10*time.Minute)
	GlobalDriver = cache.New(0, 0)

	// 订单对象缓存池
	// 订单时长不确定，让司机和用户手动删除
	GlobalOrder = cache.New(0, 0)
)
