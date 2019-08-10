package client

import (
	"net/rpc"
	"HFish/utils/log"
	"HFish/utils/conf"
	"HFish/utils/ip"
)

// 上报状态结构
type Status struct {
	AgentIp                                    string
	AgentName                                  string
	Web, Deep, Ssh, Redis, Mysql, Http, Telnet string
}

// 上报结果结构
type Result struct {
	AgentIp     string
	AgentName   string
	Type        string
	ProjectName string
	SourceIp    string
	Info        string
	Id          string // 数据库ID，更新用 0 为新插入数据
}

func createClient() (*rpc.Client, bool) {
	rpcAddr := conf.Get("rpc", "addr")
	client, err := rpc.Dial("tcp", rpcAddr)

	if err != nil {
		log.Pr("RPC", "127.0.0.1", "连接 RPC Server 失败")
		return client, false
	}

	return client, true
}

func reportStatus(rpcName string, telnetStatus string, httpStatus string, mysqlStatus string, redisStatus string, sshStatus string, webStatus string, darkStatus string) {
	client, boolStatus := createClient()

	if boolStatus {
		defer client.Close()

		status := Status{
			ip.GetLocalIp(),
			rpcName,
			webStatus,
			darkStatus,
			sshStatus,
			redisStatus,
			mysqlStatus,
			httpStatus,
			telnetStatus,
		}

		var reply string
		err := client.Call("HFishRPCService.ReportStatus", status, &reply)

		if err != nil {
			log.Pr("RPC", "127.0.0.1", "上报服务状态失败", err)
		}
	}
}

func ReportResult(typex string, projectName string, sourceIp string, info string, id string) string {
	// projectName 只有 WEB 才需要传项目名 其他协议空即可
	// id 0 为 新插入数据，非 0 都是更新数据
	// id 非 0 的时候 sourceIp 为空
	client, boolStatus := createClient()

	if boolStatus {
		defer client.Close()

		rpcName := conf.Get("rpc", "name")

		result := Result{
			ip.GetLocalIp(),
			rpcName,
			typex,
			projectName,
			sourceIp,
			info,
			id,
		}

		var reply string
		err := client.Call("HFishRPCService.ReportResult", result, &reply)

		if err != nil {
			log.Pr("RPC", "127.0.0.1", "上报上钩结果失败")
		}

		return reply
	}
	return ""
}

func Start(rpcName string, telnetStatus string, httpStatus string, mysqlStatus string, redisStatus string, sshStatus string, webStatus string, darkStatus string) {
	reportStatus(rpcName, telnetStatus, httpStatus, mysqlStatus, redisStatus, sshStatus, webStatus, darkStatus)
}
