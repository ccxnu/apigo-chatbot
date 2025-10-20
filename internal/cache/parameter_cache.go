package cache

import (
	"sync"

	d "api-chatbot/domain"
)

type parameterCache struct {
	mu     sync.RWMutex
	params map[string]*d.Parameter
}

func NewParameterCache() d.ParameterCache {
	return &parameterCache{
		params: make(map[string]*d.Parameter),
	}
}

func (c *parameterCache) Get(code string) (*d.Parameter, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	param, exists := c.params[code]
	return param, exists
}

func (c *parameterCache) GetValue(code string) (d.Data, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	param, exists := c.params[code]
	if !exists {
		return nil, false
	}

	data, err := param.GetDataAsMap()
	if err != nil {
		return nil, false
	}

	return data, true
}

func (c *parameterCache) Set(code string, param *d.Parameter) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.params[code] = param
}

func (c *parameterCache) Delete(code string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.params, code)
}

func (c *parameterCache) LoadAll(params []d.Parameter) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.params = make(map[string]*d.Parameter)

	for i := range params {
		c.params[params[i].Code] = &params[i]
	}
}

func (c *parameterCache) GetAll() []d.Parameter {
	c.mu.RLock()
	defer c.mu.RUnlock()

	result := make([]d.Parameter, 0, len(c.params))
	for _, param := range c.params {
		result = append(result, *param)
	}

	return result
}

func (c *parameterCache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.params = make(map[string]*d.Parameter)
}
