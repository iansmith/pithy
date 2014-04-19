package pithy

import (
	"fmt"
	"github.com/gopherjs/gopherjs/js"
	"github.com/gopherjs/jquery"
	"reflect"
	"strings"
)

type Pithy struct {
	Tag      string
	Id       string
	Classes  []string
	Attrs    map[string]string
	Children []*Pithy
}

const (
	seq_none     = 0
	seq_id       = 1
	seq_classes  = 2
	seq_attr     = 3
	seq_children = 4
)

type seq int

func (t seq) String() string {
	switch t {
	case seq_none:
		return "NONE"
	case seq_id:
		return "ID"
	case seq_classes:
		return "CLASSES"
	case seq_attr:
		return "ATTRIBUTES"
	case seq_children:
		return "CHILDREN"
	}
	panic("unknown sequence element")
}

func (p *Pithy) Root() jquery.JQuery {
	id := ""
	classes := ""

	if p.Id != "" {
		id = fmt.Sprintf("id='%s'", p.Id)
	}

	if p.Classes != nil {
		classes = fmt.Sprintf("class='%s'", strings.Join(p.Classes, " "))
	}

	t := fmt.Sprintf("<%s %s %s/>", p.Tag, id, classes)
	parsed := jquery.ParseHTML(t)
	first := parsed[0].(js.Object)
	j := jquery.NewJQuery(first)
	for k, v := range p.Attrs {
		j.SetAttr(k, v)
	}

	for _, child := range p.Children {
		j.Append(child.Root())
	}

	return j
}

func IMG(obj ...interface{}) *Pithy {
	return tag("img", obj...)
}
func DIV(obj ...interface{}) *Pithy {
	return tag("div", obj...)
}

func INPUT(obj ...interface{}) *Pithy {
	return tag("input", obj...)
}

func TEXTAREA(obj ...interface{}) *Pithy {
	return tag("textarea", obj...)
}

func LABEL(obj ...interface{}) *Pithy {
	return tag("label", obj...)
}

func A(obj ...interface{}) *Pithy {
	return tag("a", obj...)
}

func SPAN(obj ...interface{}) *Pithy {
	return tag("span", obj...)
}

func H1(obj ...interface{}) *Pithy {
	return tag("h1", obj...)
}

func H2(obj ...interface{}) *Pithy {
	return tag("h2", obj...)
}

func H3(obj ...interface{}) *Pithy {
	return tag("h3", obj...)
}

func H4(obj ...interface{}) *Pithy {
	return tag("h4", obj...)
}

func H5(obj ...interface{}) *Pithy {
	return tag("h5", obj...)
}

func H6(obj ...interface{}) *Pithy {
	return tag("h6", obj...)
}

func HR(obj ...interface{}) *Pithy {
	return tag("hr", obj...)
}

func tag(tagName string, obj ...interface{}) *Pithy {
	state := seq_none
	p := &Pithy{Tag: tagName}
	for i := 0; i < len(obj); i++ {
		if obj[i] == nil {
			panic("pithy does not allow nil values, try omitting the value you don't think you need")
		}
		t := reflect.TypeOf(obj[i])

		if t.Kind() == reflect.String {
			if state >= seq_id {
				panic(fmt.Sprintf("found an id but was expecting at least %s", (state + 1)))
			}
			state = seq_id
			p.Id = obj[i].(string)
			continue
		}
		if t.Kind() == reflect.Slice {
			if state >= seq_classes {
				panic(fmt.Sprintf("found a slice of classes but was expecting at least %s", (state + 1)))
			}
			if t.Elem().Kind() != reflect.String {
				panic(fmt.Sprintf("found []%s but expected []string for classes", t.Elem().Kind()))
			}
			p.Classes = obj[i].([]string)
			state = seq_classes
			continue
		}
		if t.Kind() == reflect.Map {
			if state >= seq_attr {
				panic(fmt.Sprintf("found a map of attributes but was expecting at least %s", (state + 1)))
			}
			if t.Key().Kind() != reflect.String {
				panic(fmt.Sprintf("found map but expected key to be string for attributes (got %s)", t.Key().Kind()))
			}
			if t.Elem().Kind() != reflect.String {
				panic(fmt.Sprintf("found map but expected key to be string for attributes (got %s)", t.Key().Kind()))
			}
			p.Attrs = obj[i].(map[string]string)
			state = seq_classes
			continue
		}
		if t.Kind() == reflect.Ptr {
			if t.Elem().Name() != "Pithy" {
				panic(fmt.Sprintf("expected pointer to Pithy but got pointer to %s", t.Elem().Name()))
			}
			state = seq_children
			p.Children = append(p.Children, obj[i].(*Pithy))
			continue
		}

		panic(fmt.Sprintf("unable to understand type of parameter: %v", obj[i]))
	}
	return p
}
