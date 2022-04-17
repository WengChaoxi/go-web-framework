package framework

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"strconv"
	"sync"
	"text/template"
	"time"
)

type Context struct {
	request *http.Request
	writer  http.ResponseWriter
	ctx     context.Context

	handlers []HandlerFunc // 调用链：中间件a -> 中间件b -> ... -> 业务逻辑
	index    int           // 当前请求调用到调用链的位置, 默认 -1

	mu   *sync.RWMutex
	keys map[string]interface{}

	hasTimeout bool // 是否超时
}

func NewContext(req *http.Request, rw http.ResponseWriter) *Context {
	return &Context{
		request: req,
		writer:  rw,
		ctx:     req.Context(),
		index:   -1,
		mu:      &sync.RWMutex{},
	}
}

//
// 基本函数功能
//

func (c *Context) WriterMux() *sync.RWMutex {
	return c.mu
}

func (c *Context) Request() *http.Request {
	return c.request
}

func (c *Context) ResponseWriter() http.ResponseWriter {
	return c.writer
}

func (c *Context) SetHasTimeout() {
	c.hasTimeout = true
}

func (c *Context) HasTimeout() bool {
	return c.hasTimeout
}

func (c *Context) Set(key string, value interface{}) {
	c.mu.Lock()
	if c.keys == nil {
		c.keys = make(map[string]interface{})
	}
	c.keys[key] = value
	c.mu.Unlock()
}

func (c *Context) Get(key string) (value interface{}, exists bool) {
	c.mu.RLock()
	value, exists = c.keys[key]
	c.mu.Unlock()
	return
}

func (c *Context) GetString(key string) (s string) {
	if value, ok := c.Get(key); ok && value != nil {
		s, _ = value.(string)
	}
	return
}

// 设置 handlers 调用链
func (c *Context) SetHandlers(handlers []HandlerFunc) {
	c.handlers = handlers
}

// 调用调用链的下一个函数
func (c *Context) Next() error {
	c.index++
	if c.index < len(c.handlers) {
		err := c.handlers[c.index](c)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *Context) BaseContext() context.Context {
	return c.request.Context()
}

//
// 实现 context.Context 的接口
//

func (c *Context) Deadline() (deadline time.Time, ok bool) {
	return c.BaseContext().Deadline()
}

func (c *Context) Done() <-chan struct{} {
	return c.BaseContext().Done()
}

func (c *Context) Err() error {
	return c.BaseContext().Err()
}

func (c *Context) Value(key interface{}) interface{} {
	return c.BaseContext().Value(key)
}

//
// 获取 URL 参数相关方法
//

func (c *Context) QueryInt(key string, default_ int) int {
	data := c.Query()
	if tmp, ok := data[key]; ok {
		len := len(tmp)
		if len > 0 {
			val, err := strconv.Atoi(tmp[len-1])
			if err != nil {
				return default_
			}
			return val
		}
	}
	return default_
}

func (c *Context) QueryString(key string, default_ string) string {
	data := c.Query()
	if tmp, ok := data[key]; ok {
		len := len(tmp)
		if len > 0 {
			return tmp[len-1]
		}
	}
	return default_
}

func (c *Context) QueryArray(key string, default_ []string) []string {
	data := c.Query()
	if tmp, ok := data[key]; ok {
		return tmp
	}
	return default_
}

func (c *Context) Query() map[string][]string {
	if c.request != nil {
		return map[string][]string(c.request.URL.Query())
	}
	return map[string][]string{}
}

//
// 获取 Form 数据相关方法 (post)
//

func (c *Context) FormInt(key string, default_ int) int {
	data := c.PostForm()
	if tmp, ok := data[key]; ok {
		len := len(tmp)
		if len > 0 {
			val, err := strconv.Atoi(tmp[len-1])
			if err != nil {
				return default_
			}
			return val
		}
	}
	return default_
}

func (c *Context) FormString(key string, default_ string) string {
	data := c.PostForm()
	if tmp, ok := data[key]; ok {
		len := len(tmp)
		if len > 0 {
			return tmp[len-1]
		}
	}
	return default_
}

func (c *Context) FormArray(key string, default_ []string) []string {
	data := c.PostForm()
	if tmp, ok := data[key]; ok {
		return tmp
	}
	return default_
}

func (c *Context) PostForm() map[string][]string {
	if c.request != nil {
		return map[string][]string(c.request.PostForm)
	}
	return map[string][]string{}
}

//
// 获取 application/json 数据相关方法 (post)
//

func (c *Context) BindJson(object interface{}) error {
	if c.request != nil {
		body, err := ioutil.ReadAll(c.request.Body)
		if err != nil {
			return err
		}
		c.request.Body = ioutil.NopCloser(bytes.NewBuffer(body))

		err = json.Unmarshal(body, object)
		if err != nil {
			return err
		}
	} else {
		return errors.New("request empty")
	}
	return nil
}

//
// 响应数据相关方法
//

func (c *Context) Json(status int, object interface{}) error {
	if c.HasTimeout() {
		return nil
	}
	c.writer.Header().Set("Content-Type", "application/json")
	c.writer.WriteHeader(status)
	bytes_, err := json.Marshal(object)
	if err != nil {
		c.writer.WriteHeader(500)
		return err
	}
	c.writer.Write(bytes_)
	return nil
}

func (c *Context) HTML(template_ string, object interface{}) error {
	t, err := template.New("output").Parse(template_)
	if err != nil {
		fmt.Printf("%s\n", err.Error())
		return err
	}
	if err := t.Execute(c.writer, object); err != nil {
		fmt.Printf("%s\n", err.Error())
		return err
	}
	c.writer.Header().Set("Content-Type", "application/html")
	return nil
}

func (c *Context) HTMLFromFile(filename string, object interface{}) error {
	_, fn := filepath.Split(filename)
	t, err := template.New(fn).ParseFiles(filename)
	if err != nil {
		fmt.Printf("%s\n", err.Error())
		return err
	}
	if err := t.Execute(c.writer, object); err != nil {
		fmt.Printf("%s\n", err.Error())
		return err
	}
	c.writer.Header().Set("Content-Type", "application/html")
	return nil
}

func (c *Context) Text(text string) error {
	c.writer.Header().Set("Content-Type", "text/plain")
	c.writer.WriteHeader(200)
	c.writer.Write([]byte(text))
	return nil
}
