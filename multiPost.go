package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func MultiPost(workerHostname string, postUrl string) {
	//http post到指定的url，获得list-status的结果
	defer GetMultiPostWg().Done()
	client := &http.Client{}
	getReq, err := http.NewRequest("POST", postUrl, nil)
	if err != nil {
		log.Println("http.NewRequest err: ", err)
	}
	getResp, err := client.Do(getReq)
	if err != nil {
		log.Println("http.Get err: ", err)
	}

	// Check the response status code
	if getResp.StatusCode == http.StatusOK {
		// The server has started the training, exit the loop
		log.Println("Server started reading.")

	} else {
		log.Println("Server is not ready. Status code:", getResp.StatusCode)
	}
	defer getResp.Body.Close()
	var output interface{}
	if err := json.NewDecoder(getResp.Body).Decode(&output); err != nil {
		log.Println("Decode() failed: ", err)
	}
	//得到一个filesList的list，包括了UFS中所有的file的信息
	filesList, ok := output.([]interface{})
	if !ok {
		fmt.Println("not ok")
	}
	//如果循环处理filesList中的每一个file，得到一个map
	//时间效率很低，所以使用goroutine可以并行处理list中的每一个file
	for _, v := range filesList {
		//这里应该会开启300个goroutine，每个goroutine处理一个file
		//因为UFS中有300个文件
		GetFileListWg().Add(1)
		go processFile(v, workerHostname)
	}
	GetFileListWg().Wait()

}

func processFile(file interface{}, workerHostname string) {
	defer GetFileListWg().Done()
	fileMap, ok := file.(map[string]interface{})
	if !ok {
		fmt.Println("not ok")
		return
	}
	name := fileMap["name"].(string)
	lastAccessTime := fileMap["lastAccessTimeMs"].(float64)
	fileInfo := fileMap["fileInfo"].(map[string]interface{})
	fileBlockInfos := fileInfo["fileBlockInfos"].([]interface{})
	//fileBlockInfosLen := len(fileBlockInfos)
	//fmt.Println("fileBlockInfosLen:", fileBlockInfosLen)
	//fmt.Println("fileBlockInfos:", fileBlockInfos)
	for _, v := range fileBlockInfos {

		blockInfo := v.(map[string]interface{})
		//fmt.Println(blockInfo["blockInfo"])
		locations, ok := blockInfo["blockInfo"].(map[string]interface{})["locations"].([]interface{})

		//blockinfo, ok := blockInfo["blockInfo"].(map[string][]interface{})
		//fmt.Println("blockinfo:", blockinfo)
		if !ok {
			//fmt.Println("locations is not ok")
			return
		}
		for _, v := range locations {
			workerAddress := v.(map[string]interface{})["workerAddress"].(map[string]interface{})
			host := workerAddress["host"].(string)
			SetFile(host, name, lastAccessTime)
		}
		//fmt.Println("locations:", locations)

		//for _, v := range locations {
		//	_ := v.(map[string]interface{})["workerAddress"].(map[string]interface{})
		//	host, ok := workerAddress["host"].(string)
		//	if !ok {
		//		fmt.Printf("file %v, is not in alluxio", name)
		//		return
		//	}
		//	workerinfo := GetWorkerInfoList()[host]
		//	workerinfo.Files = append(workerinfo.Files, name)
		//	workerinfo.LastAccessTime = lastAccessTime
		//}
	}
	//blockInfo := fileBlockInfos[0].(map[string]interface{})
	//if !ok {
	//	//fmt.Printf("file %v, is not in alluxio", name)
	//	fmt.Println("locations is not ok")
	//	return
	//}

	//fmt.Printf("worker:%s,filename:%s,lastAccessTime:%v,\n", workerHostname, name, lastAccessTime)

	//将fileMap的内容存放的WorkerInfo中
	//选取重要的指标采集进入WorkerInfo中

}
