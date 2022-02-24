package framework

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strconv"
	"sync"
	"time"
)

type Context struct {
	req     *http.Request
	rw      http.ResponseWriter
	ctx     context.Context
	handler Handler

	hasTimeout bool        // 是否超时
	writerMux  *sync.Mutex // 写锁
}

func NewContext(req *http.Request, rw http.ResponseWriter) *Context {
	return &Context{
		req:       req,
		rw:        rw,
		ctx:       req.Context(),
		writerMux: &sync.Mutex{},
	}
}

//
// 基本函数功能
//

func (c *Context) WriterMux() *sync.Mutex {
	return c.writerMux
}

func (c *Context) GetRequest() *http.Request {
	return c.req
}

func (c *Context) GetResponse() http.ResponseWriter {
	return c.rw
}

func (c *Context) SetHasTimeout() {
	c.hasTimeout = true
}

func (c *Context) HasTimeout() bool {
	return c.hasTimeout
}

func (c *Context) BaseContext() context.Context {
	return c.req.Context()
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

func (c *Context) Query() map[string][]string {
	if c.req != nil {
		return map[string][]string(c.req.URL.Query())
	}
	return map[string][]string{}
}

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

//
// 获取 Form 数据相关方法 (post)
//

func (c *Context) PostForm() map[string][]string {
	if c.req != nil {
		return map[string][]string(c.req.PostForm)
	}
	return map[string][]string{}
}

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

//
// 获取 application/json 数据相关方法 (post)
//

func (c *Context) BindJson(object interface{}) error {
	if c.req != nil {
		body, err := ioutil.ReadAll(c.req.Body)
		if err != nil {
			return err
		}
		c.req.Body = ioutil.NopCloser(bytes.NewBuffer(body))

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
	c.rw.Header().Set("Content-Type", "application/json")
	c.rw.WriteHeader(status)
	bytes_, err := json.Marshal(object)
	if err != nil {
		c.rw.WriteHeader(500)
		return err
	}
	c.rw.Write(bytes_)
	return nil
}

func (c *Context) HTML(status int, object interface{}, template string) error {
	// TODO
	return nil
}

func (c *Context) Text(status int, object string) error {
	// TODO
	return nil
}
