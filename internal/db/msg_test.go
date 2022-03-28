package db

import (
	"github.com/stretchr/testify/assert"
	"strconv"
	"testing"
	"time"
)

func TestNewMsg(t *testing.T) {
	t0 := time.Now()
	msg := NewMsg("unicorn", "a message")

	assert.Equal(t, "unicorn", msg.Id)
	assert.Equal(t, "a message", msg.Content)
	assert.False(t, msg.IsPalindrome)
	assert.True(t, msg.ModTime.After(t0))
}

func TestIsPalindrome(t *testing.T) {
	testDetails := []struct {
		msg          string
		isPalindrome bool
	}{
		{
			msg:          "",
			isPalindrome: true,
		},
		{
			msg:          "a",
			isPalindrome: true,
		},
		{
			msg:          "kayak",
			isPalindrome: true,
		},
		{
			msg:          "Step on no pets",
			isPalindrome: true,
		},
		{
			msg:          "Live on time emit no evil",
			isPalindrome: true,
		},
		{
			msg:          "613A316",
			isPalindrome: true,
		},
		{
			msg:          " 6 7 8   8 7 6 ",
			isPalindrome: true,
		},
		{
			msg:          "12345321",
			isPalindrome: false,
		},
		{
			msg:          "potato",
			isPalindrome: false,
		},
		{
			msg:          "kjbhwro348hf9pni io  phof sdfhj   //  adsfa",
			isPalindrome: false,
		},
		{
			msg:          "===s===a==",
			isPalindrome: false,
		},
	}

	for i, test := range testDetails {
		t.Run("test#"+strconv.Itoa(i), func(t *testing.T) {
			assert.Equal(t, test.isPalindrome, isPalindrome(test.msg))
		})
	}
}
