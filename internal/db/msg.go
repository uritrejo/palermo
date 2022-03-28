package db

import (
	"fmt"
	"strconv"
	"time"
)

// Msg will be used as both an input and an output structure to describe a message
type Msg struct {
	Id           string    `json:"id"           bson:"id"`
	Content      string    `json:"content"      bson:"content"`
	IsPalindrome bool      `json:"isPalindrome" bson:"isPalindrome"`
	ModTime      time.Time `json:"modTime"      bson:"modTime"`
}

func NewMsg(id, content string) *Msg {
	msg := &Msg{
		Id:           id,
		Content:      content,
		IsPalindrome: isPalindrome(content),
		ModTime:      time.Now(),
	}

	return msg
}

func (m *Msg) String() string {
	return fmt.Sprintf("Msg: { id: %s, content: %s, isPalindrome: %s, modTime: %s }",
		m.Id, m.Content, strconv.FormatBool(m.IsPalindrome), m.ModTime.Format(time.RFC822Z))
}

// isPalindrome returns true if the given string is a palindrome, false otherwise
func isPalindrome(sequence string) bool {
	isPalindrome := true
	l := len(sequence)
	for i := 0; i < l/2; i++ {
		if sequence[i] != sequence[l-1-i] {
			isPalindrome = false
			break
		}
	}
	return isPalindrome
}
