package coms

import (
	"container/list"
	"fmt"
	"strconv"
	s "strings"

	ss "github.com/deckarep/golang-set"
	m "github.com/lithiferous/succ/models"
)

func DropCmd(q string, sep string) []string {
	return s.SplitAfter(q, sep+" ")
}

func Delta(args []string, sep string, l *list.List) (string, m.Delta) {
	str := s.Split(args[1], sep)
	if len(args) != 2 {
		return fmt.Sprintf("incorrect number of args"), m.Delta{}
	}

	n := s.TrimSpace(s.Join(str[:len(str)-1], " "))
	if !PayerSig(n, l) {
		return fmt.Sprintf("name - %s, does not exist :c", n), m.Delta{}
	}

	v, e := strconv.ParseFloat(str[len(str)-1], 64)

	if e != nil {
		return fmt.Sprintf("number error: %s", e), m.Delta{}
	}

	return "", m.Delta{Name: n, Value: v}
}

func Person(args []string, sep string, l *list.List) (string, string, *m.Person) {
	str := s.Split(args[1], sep)
	if len(args) != 2 {
		return "incorrect number of args", "", nil
	}
	//each human described by name surname mdlname
	n := s.TrimSpace(s.Join(str[len(str)-3:], " "))
	if !PersonSig(n, l) {
		return fmt.Sprintf("name %s - does not exist :c\n", n), "", nil
	}

	return "", str[0], PersonGet(&l, n)
}

func Pay(p []string, sep string, l **ss.Set) string {
	pp := s.Split(p[1], sep)
	v, e := strconv.ParseFloat(pp[len(pp)-1], 64)

	if e != nil {
		return fmt.Sprintf("number error: %s", e)
	}

	pay := -v / float64((*(*l)).Cardinality())

	it := (*(*l)).Iterator()
	for el := range it.C {
		el.(*m.Person).Value += pay
	}
	return Status((*(*l)))
}
