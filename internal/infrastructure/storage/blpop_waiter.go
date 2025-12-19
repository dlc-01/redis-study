package storage

type blpopWaiter struct {
	ch chan string
}
