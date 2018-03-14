package main

import(
	//"github.com/samuel/go-zookeeper/zk"
	//"fmt"
	//"time"
)

func main() {



	/*servers := []string{ "a", "b", "c", "d" }
	fmt.Println(servers)

	var wg sync.WaitGroup
	wg.Add(len(servers))

	var leaders []string
	leaderChl := make(chan string)
	doneChl := make(chan struct{})
	for _, s := range servers {

		//time.Sleep(500 * time.Millisecond)
		go func(server string) {
			defer wg.Done()
			leaderChl <- server
		}(s)
	}

	go func() {
		defer close(doneChl)
		wg.Wait()
	}()

loop:
	for {
		select {
		case leader := <-leaderChl:
			fmt.Println("GET One: %v", leader)
			fmt.Println("GET One: %v", &leader)
			leaders = append(leaders, leader)
		case <-doneChl:
			fmt.Println("DONE")
			break loop
		case <-time.After(time.Second * 2):
			// In case of timeout, return what we have so far.
			fmt.Println("TIMEOUT")
			break loop
		}
	}

	fmt.Println(leaders)
*/
}
