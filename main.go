package main

import (
	"fmt"
	"time"
)

func main() {
	InstanceMap := Getec2Instance()
	//获取所有worker的url
	workerUrl := make(map[string]string)
	InitWorker("Ec2Cluster-default-workers-0")
	workerUrl["Ec2Cluster-default-workers-0"] = "http://" + InstanceMap["Ec2Cluster-default-workers-0"] + ":39999/api/v1/paths///list-status"

	//计算处理时间
	start := time.Now()

	for k, v := range workerUrl {
		GetMultiPostWg().Add(1)
		go MultiPost(k, v)
	}
	GetMultiPostWg().Wait()
	elapsed := time.Since(start)
	fmt.Println("alluxio worker读取ufs的时间为：", elapsed)
	fmt.Println("workerList: ", workerInfoListInstance["Ec2Cluster-workers-0"])
	fmt.Println(len(workerInfoListInstance["Ec2Cluster-workers-0"].FileName))

}
