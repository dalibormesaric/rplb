package loadbalancing

import (
	"fmt"
	"testing"

	"github.com/dalibormesaric/rplb/internal/backend"
)

const (
	leastLoadedBpName string = RoundRobin
	leastLoadedB1     string = "http://a:1234"
	leastLoadedB2     string = "http://b:1234"
	leastLoadedB3     string = "http://c:1234"
)

func TestGet2(t *testing.T) {
	bs := func() []*backend.Backend {
		bp, _ := backend.NewBackendPool(fmt.Sprintf("%s,%s,%s,%s,%s,%s", leastLoadedBpName, leastLoadedB1, leastLoadedBpName, leastLoadedB2, leastLoadedBpName, leastLoadedB3))
		return bp[leastLoadedBpName]
	}()

	leastloaded := &leastLoaded{state: &leastLoadedState{
		loadForBackend:    make(map[string]int),
		roundRobinForLoad: make(map[int]int),
	}}
	leastloaded.state.loadForBackend[bs[0].Name] = 2
	leastloaded.state.loadForBackend[bs[1].Name] = 1
	leastloaded.state.loadForBackend[bs[2].Name] = 1

	fmt.Printf("%s %d\n", bs[0].Name, leastloaded.state.loadForBackend[bs[0].Name])
	fmt.Printf("%s %d\n", bs[1].Name, leastloaded.state.loadForBackend[bs[1].Name])
	fmt.Printf("%s %d\n", bs[2].Name, leastloaded.state.loadForBackend[bs[2].Name])

	b1, f1 := leastloaded.Get2("", bs)
	fmt.Printf("BACKEND %s\n", b1.Name)

	fmt.Printf("%s %d\n", bs[0].Name, leastloaded.state.loadForBackend[bs[0].Name])
	fmt.Printf("%s %d\n", bs[1].Name, leastloaded.state.loadForBackend[bs[1].Name])
	fmt.Printf("%s %d\n", bs[2].Name, leastloaded.state.loadForBackend[bs[2].Name])

	b2, f2 := leastloaded.Get2("", bs)
	fmt.Printf("BACKEND %s\n", b2.Name)

	fmt.Printf("%s %d\n", bs[0].Name, leastloaded.state.loadForBackend[bs[0].Name])
	fmt.Printf("%s %d\n", bs[1].Name, leastloaded.state.loadForBackend[bs[1].Name])
	fmt.Printf("%s %d\n", bs[2].Name, leastloaded.state.loadForBackend[bs[2].Name])

	f1()
	fmt.Printf("%s %d\n", bs[0].Name, leastloaded.state.loadForBackend[bs[0].Name])
	fmt.Printf("%s %d\n", bs[1].Name, leastloaded.state.loadForBackend[bs[1].Name])
	fmt.Printf("%s %d\n", bs[2].Name, leastloaded.state.loadForBackend[bs[2].Name])

	f2()
	fmt.Printf("%s %d\n", bs[0].Name, leastloaded.state.loadForBackend[bs[0].Name])
	fmt.Printf("%s %d\n", bs[1].Name, leastloaded.state.loadForBackend[bs[1].Name])
	fmt.Printf("%s %d\n", bs[2].Name, leastloaded.state.loadForBackend[bs[2].Name])

	t.Error("error")
}
