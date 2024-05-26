package parser

import (
	"errors"
	"ethereum-parser/internal/infrastructure/remotely"
	"ethereum-parser/internal/infrastructure/storage"
	"ethereum-parser/internal/models"
	"log"
	"strings"
	"sync"
	"time"
)

var (
	p Parser
)

func SetParser(par Parser) {
	p = par
}

func GetParser() Parser {
	return p
}

type Parser interface {
	// GetCurrentBlock last parsed block
	GetCurrentBlock() int

	// Subscribe add address to observer
	Subscribe(address string) bool

	// GetTransactions list of inbound or outbound transactions for an address
	GetTransactions(address string) []models.Transaction
}

func NewParser() (Parser, error) {
	return NewEthParser()
}

type EthParser struct {
	subscribes     map[string]bool // Subscribe
	subscribesLock *sync.RWMutex   // Subscribe

	latestBlock int            //
	cache       *storage.Cache // trans cache

	_localBlock      int
	_pool            chan struct{}
	_retryBlocks     []int
	_retryBlocksLock *sync.Mutex
	_retryInterval   int
	_updateInterval  int
}

func NewEthParser() (*EthParser, error) {
	e := &EthParser{}
	b, err := remotely.GetCurrentBlock()
	if err != nil {
		return nil, errors.New("init parser fail")
	}
	e.latestBlock = int(b)
	e.subscribes = make(map[string]bool)
	e.subscribesLock = &sync.RWMutex{}
	e.cache = storage.NewCache()

	e._localBlock = 0 // init
	e._retryBlocksLock = &sync.Mutex{}
	e._retryInterval = 10             // retry seconds, can be configured
	e._updateInterval = 10            // update interval
	e._pool = make(chan struct{}, 10) // work pool size, can be configured
	//
	go e.doUpdateLatestBlock()
	go e.doUpdate()
	go e.doRetry()
	go e.showLogs()
	return e, nil
}
func (e *EthParser) doUpdateLatestBlock() {
	for {
		b, err := remotely.GetCurrentBlock()
		if err == nil {
			e.latestBlock = int(b)
		}
		time.Sleep(time.Second * time.Duration(e._updateInterval))
	}
}

func (e *EthParser) doUpdate() {
	for {
		if e._localBlock < e.latestBlock {
			e._pool <- struct{}{}
			go func(num int) {
				_ = e.updateByNumber(num, false)
			}(e._localBlock)
			e._localBlock++
		} else {
			time.Sleep(time.Second * time.Duration(e._updateInterval))
		}
	}
}

func (e *EthParser) doRetry() {
	for {
		if len(e._retryBlocks) == 0 {
			// no retry
			time.Sleep(time.Second * time.Duration(e._retryInterval))
			continue
		}

		//
		e._retryBlocksLock.Lock()
		// has retry
		newRetryBlocks := make([]int, 0)
		for _, b := range e._retryBlocks {
			if e.updateByNumber(b, true) {
				newRetryBlocks = append(newRetryBlocks, b) // keep retry
			}
		}
		e._retryBlocks = newRetryBlocks
		e._retryBlocksLock.Unlock()

		time.Sleep(time.Second * time.Duration(e._retryInterval)) // wait time interval
	}
}

func (e *EthParser) updateByNumber(num int, inRetry bool) (retry bool) {
	trans, err := remotely.GetTransactionByNumber(int64(num))
	defer func() {
		<-e._pool
	}()
	if err != nil {
		if !inRetry {
			e._retryBlocksLock.Lock()
			e._retryBlocks = append(e._retryBlocks, num)
			e._retryBlocksLock.Unlock()
			log.Printf("eth block:[%s] needs retry", num)
		}
		retry = true
		return
	}
	for _, t := range trans {
		// if in subscribes
		fromSub, toSub := e.ifSubscribe(&t)
		if fromSub {
			log.Printf("[subscribe] attention, subscribed user:[%s] has make a trans: [%v] in block :[%d]", t.From, t, num)
		}
		if toSub {
			log.Printf("[subscribe] attention, subscribed user:[%s] has make a trans: [%v] in block:[%d]", t.To, t, num)
		}
		_ = e.cache.Add(t)
	}
	return
}

func (e *EthParser) showLogs() {
	for {
		log.Printf("[info] sync block successfully to :[%d], retry block length:[%d]", e._localBlock, len(e._retryBlocks))
		log.Printf("[info] current subscribes:[%s]", e.listSubscribes())
		time.Sleep(time.Second * 5)
	}
}

func (e *EthParser) listSubscribes() string {
	if e == nil {
		return ""
	}
	res := make([]string, len(e.subscribes))
	index := 0
	for s, _ := range e.subscribes {
		res[index] = s
		index++
	}
	return strings.Join(res, ",")
}

func (e *EthParser) ifSubscribe(tran *models.Transaction) (fromSub, toSub bool) {
	// if Subscribe, then print
	if tran == nil {
		return false, false
	}
	e.subscribesLock.RLock()
	_, ok := e.subscribes[tran.From]
	if ok {
		fromSub = true
	}
	_, ok = e.subscribes[tran.To]
	if ok {
		toSub = true
	}
	e.subscribesLock.RUnlock()
	return
}

func (e *EthParser) GetCurrentBlock() int {
	if e == nil {
		return 0
	}
	return e.latestBlock
}

func (e *EthParser) Subscribe(address string) bool {
	if e == nil {
		return false
	}

	e.subscribesLock.Lock()
	_, ok := e.subscribes[address]
	if !ok {
		e.subscribes[address] = true
	}
	e.subscribesLock.Unlock()

	return true
}

func (e *EthParser) GetTransactions(address string) []models.Transaction {
	res, _ := e.cache.Get(address)
	return res
}
