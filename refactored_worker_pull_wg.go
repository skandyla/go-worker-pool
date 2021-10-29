package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

var actions = []string{"logged in", "logged out", "created record", "deleted record", "updated account"}

type logItem struct {
	action    string
	timestamp time.Time
}

type User struct {
	id    int
	email string
	logs  []logItem
}

func (u User) getActivityInfo() string {
	output := fmt.Sprintf("UID: %d; Email: %s;\nActivity Log:\n", u.id, u.email)
	for index, item := range u.logs {
		output += fmt.Sprintf("%d. [%s] at %s\n", index, item.action, item.timestamp.Format(time.RFC3339))
	}

	return output
}

func main() {
	rand.Seed(time.Now().Unix())
	startTime := time.Now()
	const usersCount, workerCount = 100, 100

	jobs := make(chan User, usersCount)
	results := make(chan string, usersCount)

	users := generateUsers(usersCount)

	//run our workers
	for i := 0; i < workerCount; i++ {
		go worker(i+1, jobs, results)
	}

	//send every user object to jobs channel
	for _, user := range users {
		jobs <- user
	}

	//all jobs are sended
	close(jobs)

	//collect results from results channel
	for i := 0; i < usersCount; i++ {
		fmt.Printf("result #%d: value:%v\n", i+1, <-results)
	}

	fmt.Printf("DONE! Time Elapsed: %.2f seconds\n", time.Since(startTime).Seconds())
}

func worker(id int, jobs <-chan User, results chan<- string) {
	for j := range jobs {
		result, _ := saveUserInfo(j)
		results <- result
	}
}

func saveUserInfo(user User) (string, error) {
	// -->> saving my filesystem from that IO
	//fmt.Printf("WRITING FILE FOR UID %d\n", user.id)
	//filename := fmt.Sprintf("users/uid%d.txt", user.id)
	//file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0644)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//file.WriteString(user.getActivityInfo())
	//fmt.Printf("Processing user %d\n", user.id)
	time.Sleep(time.Second)
	return fmt.Sprintf(" Processed user %d\n", user.id), nil
}

func generateUsers(count int) []User {

	users := make([]User, count)
	wg := &sync.WaitGroup{}
	mu := &sync.Mutex{}

	for i := 0; i < count; i++ {
		wg.Add(1)

		go func(i int, wg *sync.WaitGroup, mu *sync.Mutex) {
			defer wg.Done()

			mu.Lock()
			users[i] = generateUser(i)
			mu.Unlock()

		}(i, wg, mu)
	}
	wg.Wait()
	return users
}

func generateUser(i int) User {
	user := User{
		id:    i + 1,
		email: fmt.Sprintf("user%d@company.com", i+1),
		logs:  generateLogs(rand.Intn(1000)),
	}
	fmt.Printf("generated user %d\n", i+1)
	time.Sleep(time.Millisecond * 100)
	return user
}

func generateLogs(count int) []logItem {
	logs := make([]logItem, count)
	for i := 0; i < count; i++ {
		logs[i] = logItem{
			action:    actions[rand.Intn(len(actions)-1)],
			timestamp: time.Now(),
		}
	}
	return logs
}
