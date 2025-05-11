package resp

type DaemonSet struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
	Desired   int32  `json:"desired"`   //表示 DaemonSet 所需的守护进程副本数，即期望在集群中运行的守护进程副本数量。
	Current   int32  `json:"current"`   //表示当前正在运行的守护进程副本数，即当前在集群中实际运行的守护进程副本数量。
	Ready     int32  `json:"ready"`     //表示已经就绪的守护进程副本数，即当前运行且已准备好接收工作负载的守护进程副本数量。
	UpToDate  int32  `json:"upToDate"`  //表示已经更新到最新版本的守护进程副本数，即已经与 DaemonSet 的定义保持一致的守护进程副本数量
	Available int32  `json:"available"` //表示可用的守护进程副本数，即已经就绪并且可以提供服务的守护进程副本数量
	Age       int64  `json:"age"`
}
