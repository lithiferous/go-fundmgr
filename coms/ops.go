package coms

import (
	"container/list"
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"strings"

	s "github.com/deckarep/golang-set"
	m "github.com/lithiferous/succ/models"
)

func InState(filePath string) (**list.List, *s.Set) {
	l := list.New()
	set := s.NewSet()
	fl, e := os.Open(filePath)
	if e != nil {
		fmt.Println("reader politely says - I am legally blind", filePath)
		fmt.Println(e)
		return &l, &set
	}
	defer fl.Close()
	r := csv.NewReader(fl)
	r.FieldsPerRecord = -1
	all, e := r.ReadAll()
	if e != nil {
		fmt.Println(e)
		return &l, &set
	}
	//init group
	for _, rec := range all {
		v, e := strconv.ParseFloat(rec[2], 64)
		if e != nil {
			fmt.Println("err: %s", e)
			return &l, &set
		}
		set.Add(&m.Person{Name: rec[1], Value: v})
	}
	//link payers
	var t *m.Person
	for _, rec := range all {
		it := set.Iterator()
		for el := range it.C {
			if el.(*m.Person).Name == rec[1] {
				t = el.(*m.Person)
				it.Stop()
			}
		}
		l.PushBack(m.Payer{Person: t, Name: rec[0]})
	}
	return &l, &set
}

func OutState(filePath string, l **list.List) bool {
	var str strings.Builder
	for e := (*l).Front(); e != nil; e = e.Next() {
		str.WriteString(fmt.Sprintf("%s,%s,%.2f\n", e.Value.(m.Payer).Name, e.Value.(m.Payer).Person.Name, e.Value.(m.Payer).Person.Value))
	}
	f, e := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	if e != nil {
		fmt.Printf("reader err for %s: %s\n", filePath, e)
		return false
	}
	defer f.Close()
	n, e := f.WriteString(str.String())
	if e != nil {
		fmt.Printf("write err for %s: %s\n", filePath, e)
		return false
	}
	fmt.Printf("wrote %d bytes\n", n)
	return true
}

func PayerSig(p string, l *list.List) bool {
	for e := (*l).Front(); e != nil; e = e.Next() {
		if e.Value.(m.Payer).Name == p {
			fmt.Println(fmt.Sprintf("%s%s", e.Value.(m.Payer).Name, p))
			return true
		}
	}
	return false
}

func PersonSig(a string, l *list.List) bool {
	for e := (*l).Front(); e != nil; e = e.Next() {
		if e.Value.(m.Payer).Person.Name == a {
			return true
		}
	}
	return false
}

func Eval(d m.Delta, l **list.List) string {
	ass := ""
	for e := (*l).Front(); e != nil; e = e.Next() {
		if e.Value.(m.Payer).Name == d.Name {
			e.Value.(m.Payer).Person.Value += d.Value
			ass = e.Value.(m.Payer).Person.Name
			break
		}
	}
	if d.Value > 0 {
		return fmt.Sprintf("Добавил %.2f галактического кредита для %s", d.Value, ass)
	}
	return fmt.Sprintf("Вычел %.2f галактического кредита для %s", d.Value, ass)
}

func Status(set s.Set) string {
	var str strings.Builder
	it := set.Iterator()
	for el := range it.C {
		str.WriteString(fmt.Sprintf("%s: %.2f\n", el.(*m.Person).Name, el.(*m.Person).Value))
	}
	return str.String()
}

func Payer(l **list.List, ss *s.Set, a *m.Person, n string) string {
	(*l).PushBack(m.Payer{Person: a, Name: n})
	return fmt.Sprintf("Добавил %s, теперь он платит за %s\n", n, a.Name)
}

func PersonGet(l **list.List, n string) *m.Person {
	for e := (*l).Front(); e != nil; e = e.Next() {
		if e.Value.(m.Payer).Person.Name == n {
			return e.Value.(m.Payer).Person
		}
	}
	return nil
}
