package db

import (
	"fmt"
	"strconv"
	"strings"
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
// it will ignore the case, but not the whitespaces or punctuations
func isPalindrome(sequence string) bool {
	seq := strings.ToLower(sequence)
	isPalindrome := true
	l := len(seq)
	for i := 0; i < l/2; i++ {
		if seq[i] != seq[l-1-i] {
			isPalindrome = false
			break
		}
	}
	return isPalindrome
}
