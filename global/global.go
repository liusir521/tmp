package global

import (
	"github.com/patrickmn/go-cache"
	"sync"
	"time"
)

var (
	// 司机和顾客的对象缓存池
	GlobalUser   = cache.New(12*time.Hour, 10*time.Minute)
	GlobalDriver = cache.New(0, 0)

	// 订单对象缓存池
	// 订单时长不确定，让司机和用户手动删除
	GlobalOrder = cache.New(0, 0)

	// 顺风车缓存池
	// 代表司机是否发车
	GlobalTogetherOrder = cache.New(0, 0)

	// 维护乘坐顺风车的乘客
	TogetherUsers sync.Map
)
