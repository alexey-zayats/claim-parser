package parser

import (
	"context"
	"fmt"
	"github.com/alexey-zayats/claim-parser/internal/dict"
	"sort"
	"strings"
	"sync"
)

// Backend ...
type Backend interface {
	Parse(ctx context.Context, param *dict.Dict, out chan interface{}) error
}

// BackendInit ...
type BackendInit func() (Backend, error)

// Parser ...
type Parser interface {
	Add(name string, init BackendInit)
	Backend(name string) (Backend, error)
}

type parser struct {
	dict *dict.Dict
}

var instance *parser
var once sync.Once

// Instance ...
func Instance() Parser {
	if instance == nil {
		once.Do(func() {
			instance = new()
		})
	}
	return instance
}

func new() *parser {
	return &parser{
		dict: dict.New(),
	}
}

// Add ...
func (p *parser) Add(name string, init BackendInit) {
	p.dict.Set(name, init)
}

// Backend
func (p *parser) Backend(name string) (Backend, error) {

	if init, ok := p.dict.Get(name); ok == true {
		return init.(BackendInit)()
	}

	return nil, fmt.Errorf("'%v' isn't one of the '%s'", name, p.registered())
}

func (p *parser) registered() string {

	list := make([]string, 0)
	p.dict.Foreach(func(s string, i interface{}) {
		list = append(list, s)
	})

	sort.Strings(list)

	return strings.Join(list, ", ")
}
