package main

import "sync"

var multiPostWg sync.WaitGroup

var fileListWg sync.WaitGroup

func GetMultiPostWg() *sync.WaitGroup {
	return &multiPostWg
}

func GetFileListWg() *sync.WaitGroup {
	return &fileListWg
}
