package main

import (
	"bufio"
	"fmt"
	"os"
	"time"

	"github.com/cagnosolutions/safemap"
)

var COUNT = 100000

func ms() int64 {
	return time.Now().UnixNano() % 1e6 / 1e3
}

func main() {

	m1 := safemap.SafeMapInstance(32)
	//m2 := make(map[string][]byte)

	fmt.Printf("Adding %d elements to map...\n", COUNT)

	t1 := ms()

	for i := 0; i < COUNT; i++ {
		key := fmt.Sprintf("key-%d", i)
		val := fmt.Sprintf("val-%d", i)
		m1.Set(key, []byte(val))
		//m2[key] = []byte(val)
	}

	fmt.Printf("Took %dms\n", ms()-t1)

	bufio.NewReader(os.Stdin).ReadBytes('\n')
}
