package req

import corev1 "k8s.io/api/core/v1"

type Pod struct {
	Base           Base                `json:"base"` // base definition info
	Volume         []Volume            `json:"volume"`
	Network        Network             `json:"network"`
	InitContainers []Container         `json:"initContainers"`
	Containers     []Container         `json:"containers"`
	Tolerations    []corev1.Toleration `json:"tolerations"` // pod toleration params
	NodeScheduling NodeScheduling      `json:"nodeScheduling"`
}

type Base struct {
	Name          string `json:"name"`
	Labels        []Item `json:"labels"`
	Namespace     string `json:"namespace"`
	RestartPolicy string `json:"restartPolicy"` // reboot strategy: Always | Never | On-Failure
}

type Item struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type Volume struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

type Network struct {
	HostNetwork bool      `json:"hostNetwork"`
	HostName    string    `json:"hostName"`
	DnsPolicy   string    `json:"dnsPolicy"`
	DnsConfig   DnsConfig `json:"dnsConfig"`
	HostAliases []Item    `json:"hostAliases"`
}

type DnsConfig struct {
	Nameservers []string `json:"nameservers"`
}

type Container struct {
	Name            string               `json:"name"`
	Image           string               `json:"image"`
	ImagePullPolicy string               `json:"imagePullPolicy"` // Always | IfNotPresent | Never
	Tty             bool                 `json:"tty"`
	Ports           []ContainerPort      `json:"ports"`
	WorkingDir      string               `json:"workingDir"`
	Command         []string             `json:"command"`
	Args            []string             `json:"args"`
	Env             []EnvVar             `json:"env"`
	EnvsFrom        []EnvVarFromResource `json:"envsFrom"`
	Privileged      bool                 `json:"privileged"`    // Whether to enable privileged mode (e.g. root)
	Resources       Resource             `json:"resources"`     // Container application quota
	VolumeMounts    []VolumeMount        `json:"volumeMounts"`  // Mounted volumes
	StartUpProbe    ContainerProbe       `json:"startUpProbe"`  // Start the probe
	LivenessProbe   ContainerProbe       `json:"livenessProbe"` // Survival probes
	ReadinessProbe  ContainerProbe       `json:"readyProbe"`    // Readiness Probe
}

type ContainerPort struct {
	Name          string `json:"name"`
	ContainerPort int32  `json:"containerPort"`
	HostPort      int32  `json:"hostPort"`
}

type Resource struct {
	Enable      bool  `json:"enable"` // Whether to configure container quotas
	MemoryReq   int32 `json:"memory"` // memory mi
	MemoryLimit int32 `json:"memoryLimit"`
	CPUReq      int32 `json:"CPUReq"` // cpu m
	CPULimit    int32 `json:"CPULimit"`
}

type VolumeMount struct {
	MountName string `json:"mountName"`
	MountPath string `json:"mountPath"` // Mounted volume -> path in the corresponding container
	ReadOnly  bool   `json:"readOnly"`  // Read-only or not
}

type ContainerProbe struct {
	Enable    bool           `json:"enable"` // Whether to turn on the probe
	Type      string         `json:"type"`   // Probe type: tcp | http | command
	HttpGet   ProbeHTTPGet   `json:"httpGet"`
	Command   ProbeCommand   `json:"command"`
	TcpSocket ProbeTcpSocket `json:"tcpSocket"`
	ProbeTime
}

type ProbeHTTPGet struct {
	Scheme  string `json:"scheme"` // http | https
	Host    string `json:"host"`   // If "", it's an in-pod request.
	Path    string `json:"path"`
	Port    int32  `json:"port"`
	Headers []Item `json:"headers"` // http headers
}

type ProbeCommand struct {
	// cat /test/test.txt
	Command []string `json:"command"`
}

type ProbeTcpSocket struct {
	Host string `json:"host"`
	Port int32  `json:"port"`
}

type ProbeTime struct {
	InitialDelaySeconds int32 `json:"initialDelaySeconds"` // Initialize for a number of seconds before starting the probe
	PeriodSeconds       int32 `json:"periodSeconds"`       // Probe after every few seconds
	TimeoutSeconds      int32 `json:"timeoutSeconds"`      // Probe timeout time
	SuccessThreshold    int32 `json:"successThreshold"`    // Threshold number of probes to consider a probe successful
	FailureThreshold    int32 `json:"failureThreshold"`    // Threshold number of probes to consider a probe failure
}

type NodeScheduling struct {
	Type         string                        `json:"type"` // nodeName | nodeSelector | nodeAffinity
	NodeName     string                        `json:"nodeName"`
	NodeSelector []Item                        `json:"nodeSelector"`
	NodeAffinity []NodeAffinityTermExpressions `json:"nodeAffinity"`
}

type NodeAffinityTermExpressions struct {
	Key      string                      `json:"key"`
	Value    string                      `json:"value"`
	Operator corev1.NodeSelectorOperator `json:"operator"`
}

type EnvVar struct {
	Name    string `json:"name"`
	RefName string `json:"refName"`
	Value   string `json:"value"`
	Type    string `json:"type"` // configMap | secret | default(k/v)
}

type EnvVarFromResource struct {
	Name    string `json:"name"`
	RefType string `json:"refType"` // configMap | secret
	Prefix  string `json:"prefix"`
}
