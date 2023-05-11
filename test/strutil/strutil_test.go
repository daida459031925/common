package strutil

import (
	"github.com/daida459031925/common/fmt"
	"github.com/daida459031925/common/util/strutil"
	"testing"
)

func TestString(t *testing.T) {
	fmt.Println(strUtil.TrimSpaceAndEmpty("  "))                  // Output: true
	fmt.Println(strUtil.TrimSpaceAndEmpty("  hello, world!  "))   // Output: false
	fmt.Println(strUtil.FirstNonEmpty("", "hello", "world"))      // Output: "hello"
	fmt.Println(strUtil.FirstNonEmpty("", "", ""))                // Output: ""
	fmt.Println(strUtil.Concat([]string{"hello", "world"}, ", ")) // Output: "hello, world"
	fmt.Println(strUtil.CountOccurrences("hello, world!", "l"))   // Output: 3
	fmt.Println(strUtil.Reverse("hello, world!"))                 // Output: "!dlrow ,olleh"
}
