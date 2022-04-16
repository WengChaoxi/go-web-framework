package middleware

import "github.com/WengChaoxi/go-web-framework/framework"

// recovery 机制，将协程中的异常捕获
func Recovery() framework.Handler {
	return func(c *framework.Context) error {
		// 捕获 c.Next() 出现的 panic
		defer func() {
			if err := recover(); err != nil {
				c.Json(500, err)
			}
		}()

		c.Next()

		return nil
	}
}
