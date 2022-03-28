package db

import (
	"sync"
)

// BasicMsgDB stores messages in local memory in a thread safe map
type BasicMsgDB struct {
	msgs sync.Map
}

func NewBasicMsgDB() *BasicMsgDB {
	return &BasicMsgDB{}
}

func (b *BasicMsgDB) GetMsg(id string) (*Msg, error) {
	msg, exists := b.msgs.Load(id)
	if !exists {
		return nil, ErrMsgNotFound{}
	}

	return msg.(*Msg), nil
}

func (b *BasicMsgDB) GetAllMsgs() ([]*Msg, error) {
	var msgs []*Msg
	b.msgs.Range(func(k, v interface{}) bool {
		msgs = append(msgs, v.(*Msg))
		return true
	})

	return msgs, nil
}

func (b *BasicMsgDB) CreateMsg(msg *Msg) error {
	_, loaded := b.msgs.LoadOrStore(msg.Id, msg)
	if loaded {
		return ErrIdUnavailable{}
	}
	return nil
}

func (b *BasicMsgDB) UpdateMsg(newMsg *Msg) error {
	msg, exists := b.msgs.Load(newMsg.Id)
	if !exists {
		return ErrMsgNotFound{}
	}

	msg.(*Msg).Content = newMsg.Content
	msg.(*Msg).IsPalindrome = newMsg.IsPalindrome
	msg.(*Msg).ModTime = newMsg.ModTime

	return nil
}

func (b *BasicMsgDB) DeleteMsg(id string) error {
	_, loaded := b.msgs.LoadAndDelete(id)
	if !loaded {
		return ErrMsgNotFound{}
	}
	return nil
}
