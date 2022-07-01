package batch

import (
	"sync"
	"time"
)

type user struct {
	ID int64
}

var mu sync.Mutex

func addUser(new_user user, res *[]user) {
	// Lock the mutex before accessing res
	mu.Lock()
	// Unlock the mutex at the end of the function using a defer statement
	defer mu.Unlock()
	*res = append(*res, new_user)
}

func getOne(id int64) user {
	time.Sleep(time.Millisecond * 100)
	return user{ID: id}
}

func getBatch(n int64, pool int64) (res []user) {
	res = make([]user, 0, n)

	// channel - to know the amount of work
	amountUsers := make(chan int, n)
	for j := 0; j < int(n); j++ {
		amountUsers <- j
	}
	close(amountUsers)

	// Start pool workers
	var wg sync.WaitGroup
	for w := 0; w < int(pool); w++ {
		wg.Add(1)
		go func(amountUsers <-chan int, res *[]user) {
			defer wg.Done()
			for i := range amountUsers {
				addUser(getOne(int64(i)), res)
			}
		}(amountUsers, &res)
	}
	// to ensure that the worker goroutines have finished
	wg.Wait()

	return res
}
