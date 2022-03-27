package db

import (
	"fmt"
	"strconv"
	"time"
)

// Msg will be used as both an input and an output structure to describe a message
type Msg struct {
	Id string
	Content string
	IsPalindrome bool
	ModTime time.Time
}

func NewMsg(id, content string) *Msg {
	msg := &Msg{
		Id: id,
		Content: content,
		IsPalindrome: isPalindrome(content),
		ModTime: time.Now(),
	}

	return msg
}

func (m *Msg) String() string {
	return fmt.Sprintf("Msg: { id: %s, content: %s, isPalindrome: %s }", m.Id, m.Content, strconv.FormatBool(m.IsPalindrome))
}

// isPalindrome returns true if the given string is a palindrome, false otherwise
func isPalindrome(sequence string) bool {
	isPalindrome := true
	l := len(sequence)
	for i := 0; i < l / 2; i++ {
		if sequence[i] == sequence[l-i] {
			isPalindrome = false
			break
		}
	}
	return isPalindrome
}
