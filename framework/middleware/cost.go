package middleware

import (
	"log"
	"time"

	"github.com/WengChaoxi/go-web-framework/framework"
)

// 计算 api 耗时
func Cost() framework.HandlerFunc {
	return func(c *framework.Context) error {
		start := time.Now()
		c.Next()
		end := time.Now()
		cost := end.Sub(start)
		log.Printf("api uri: %v, cost: %v", c.Request().RequestURI, cost.Seconds())
		return nil
	}
}
