package main

import "sync"

type File struct {
	FileName []string
}

var workerInfoListInstance = make(map[string]File)
var mx sync.Mutex

func InitWorker(workerName string) {
	mx.Lock()
	workerInfoListInstance[workerName] = File{}
	mx.Unlock()

}

func SetFile(host string, name string, LastAccessTime float64) {

	result := GetFileName(host)
	result = append(result, name)
	tmpFile := File{
		FileName: result,
	}
	mx.Lock()
	workerInfoListInstance[host] = tmpFile
	mx.Unlock()
}

func GetFileName(host string) []string {
	mx.Lock()
	result := workerInfoListInstance[host].FileName
	mx.Unlock()
	return result
}
