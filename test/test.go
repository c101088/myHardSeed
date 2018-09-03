package main
import (
	"fmt"
	"strings"
	"time"
)
func main() {
	i := 0
	for i < 100 {
		fmt.Printf("\r")
		fmt.Printf("%s", strings.	Repeat("x", i))
		time.Sleep(time.Second)
		i++
	}
}