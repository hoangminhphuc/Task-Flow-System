package rpccaller

import (
	"flag"

	"github.com/200Lab-Education/go-sdk/logger"
)

type apiItemCaller struct {
	// We do not need client here
	name       string
	serviceURL string
	logger     logger.Logger
}

func NewApiItemCaller(name string) *apiItemCaller {
	return &apiItemCaller{
		name:       name,
	}
}

func (c *apiItemCaller) GetPrefix() string {
	return c.name
}

func (c *apiItemCaller) Get() interface{} {
	return c
}

func (c *apiItemCaller) Name() string {
	return c.name
}

func (c *apiItemCaller) InitFlags() {
	flag.StringVar(&c.serviceURL, "item-service-url", "http://localhost:8080", "URL of item service")
}

// Before running
func (c *apiItemCaller) Configure() error {
	c.logger = logger.GetCurrent().GetLogger("api.item")
	return nil
}

func (c *apiItemCaller) Run() error {
	return nil
}

func (c *apiItemCaller) Stop() <-chan bool {
	ch := make(chan bool)
	go func() {
		ch <- true
	}()
	return ch
}

func (c *apiItemCaller) GetServiceURL() string {
	return c.serviceURL
}