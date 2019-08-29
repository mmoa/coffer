package records

// Coffer
// Records storage
// Copyright © 2019 Eduard Sesigin. All rights reserved. Contacts: <claygod@yandex.ru>

import (
	//"fmt"

	"github.com/claygod/coffer/domain"
)

/*
storage - easy data storage (not parallel mode).
*/
type storage struct {
	data map[string][]byte
}

func newStorage() *storage {
	return &storage{
		data: make(map[string][]byte),
	}
}

func (r *storage) readList(keys []string) (map[string][]byte, []string) {
	//var errOut error
	notFound := make([]string, 0, len(keys))
	list := make(map[string][]byte)
	for _, key := range keys {
		if value, ok := r.data[key]; ok {
			list[key] = value
			//out = append(out, &domain.Record{Key: key, Value: value})
		} else {
			notFound = append(notFound, key)
			//errOut = fmt.Errorf("%v %v", errOut, fmt.Errorf("Key `%s` not found", key))
		}
	}
	return list, notFound
}

// func (r *storage) get(keys []string) ([]*domain.Record, error) {
// 	var errOut error
// 	out := make([]*domain.Record, 0, len(keys))
// 	for _, key := range keys {
// 		if value, ok := r.data[key]; ok {
// 			out = append(out, &domain.Record{Key: key, Value: value})
// 		} else {
// 			errOut = fmt.Errorf("%v %v", errOut, fmt.Errorf("Key `%s` not found", key))
// 		}
// 	}
// 	return out, nil
// }

func (r *storage) writeList(list map[string][]byte) {
	for key, value := range list {
		r.data[key] = value
	}
}

func (r *storage) writeOne(key string, value []byte) {
	r.data[key] = value
}

// func (r *storage) set(in []*domain.Record) {
// 	for _, rec := range in {
// 		r.data[rec.Key] = rec.Value
// 	}
// }

func (r *storage) setOne(rec *domain.Record) {
	r.data[rec.Key] = rec.Value
}

func (r *storage) removeWhatIsPossible(keys []string) ([]string, []string) {
	removedList := make([]string, 0, len(keys))
	notFound := make([]string, 0, len(keys))
	for _, key := range keys {
		if _, ok := r.data[key]; ok {
			removedList = append(removedList, key)
			delete(r.data, key)
		} else {
			notFound = append(notFound, key)
		}
	}
	return removedList, notFound
}

func (r *storage) delAllOrNothing(keys []string) []string {
	//var errOut error
	notFound := make([]string, 0, len(keys))
	for _, key := range keys { // сначала проверяем, есть ли все эти ключи
		if _, ok := r.data[key]; !ok {
			notFound = append(notFound, key)
			//errOut = fmt.Errorf("%v %v", errOut, fmt.Errorf("Key `%s` not found", key))
		}
	}
	if len(notFound) != 0 {
		return notFound
	}
	for _, key := range keys { // теперь удаляем
		delete(r.data, key)
	}
	return notFound
}

// func (r *storage) keys() []string { // Resource-intensive method
// 	out := make([]string, 0, len(r.data))
// 	for key, _ := range r.data {
// 		out = append(out, key)
// 	}
// 	return out
// }

// func (r *storage) len() int {
// 	return len(r.data)
// }

func (r *storage) iterator(chRecord chan *domain.Record, chFinish chan struct{}) {
	//fmt.Println(r.data)
	for key, value := range r.data {
		//fmt.Println("++++++ ", key, value)
		chRecord <- &domain.Record{
			Key:   key,
			Value: value,
		}
	}
	close(chFinish)
}

func (r *storage) countRecords() int {
	return len(r.data)
}
