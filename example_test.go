package envx_test

import (
	"fmt"
	"time"

	"github.com/cristalhq/envx"
)

func ExampleEnvSet() {
	envs := []string{"ENVX_TIMEOUT=20s"} // or os.Environ()

	var d time.Duration
	eset := envx.NewEnvSet("ENVX")
	eset.Duration(&d, "TIMEOUT", 10*time.Second, "just a timeout")

	err := eset.Parse(envs)
	if err != nil {
		panic(err)
	}

	fmt.Println(d)

	// Output: 20s
}
