package storage

import (
	"errors"
	"ethereum-parser/internal/models"
	"sync"
)

type Cache struct {
	TransFrom *sync.Map // key address, value inbound trans
	TransTo   *sync.Map // key address, value outbound trans
}

var (
	CacheNotInitError = errors.New("cache not init")
	LoadCacheError    = errors.New("load cache error")
	requestParamError = errors.New("request param error")
)

func NewCache() *Cache {
	return &Cache{
		TransFrom: &sync.Map{},
		TransTo:   &sync.Map{},
	}
}

type item struct {
	trans []models.Transaction
}

func add(m *sync.Map, id string, trans models.Transaction) error {
	if m == nil {
		return requestParamError
	}
	d, _ := m.Load(id)
	v, yes := d.(*item)
	if !yes || v.trans == nil {
		v = &item{
			trans: make([]models.Transaction, 0),
		}
	}
	v.trans = append(v.trans, trans)
	m.Store(id, v)
	return nil
}

func get(m *sync.Map, id string) ([]models.Transaction, error) {
	if m == nil {
		return nil, CacheNotInitError
	}

	// outbound = make([]*models.Transaction, 0)
	tv, ok := m.Load(id)
	if !ok {
		return nil, nil
	}
	t, yes := tv.(*item)
	if !yes {
		return nil, nil
	}
	res := make([]models.Transaction, len(t.trans))
	for index, value := range t.trans {
		res[index] = value
	}
	return res, nil
}

func (c *Cache) Add(transaction models.Transaction) error {
	if c == nil {
		return CacheNotInitError
	}
	if err := add(c.TransFrom, transaction.From, transaction); err != nil {
		return err
	}
	if err := add(c.TransTo, transaction.To, transaction); err != nil {
		return err
	}
	return nil
}

func (c *Cache) Get(address string) ([]models.Transaction, error) {
	if c == nil {
		return nil, CacheNotInitError
	}
	res := make([]models.Transaction, 0)
	inbound, _ := get(c.TransFrom, address)
	outbound, _ := get(c.TransTo, address)
	for _, t := range inbound {
		res = append(res, t)
	}
	for _, t := range outbound {
		res = append(res, t)
	}
	return res, nil
}
