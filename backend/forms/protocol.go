package forms

type Transfer struct {
	Action string `json:"action"`
	Data   string `json:"data"`
}

type StartFlow struct {
	InstanceId string `json:"instanceId"`
	ExePath    string `json:"exePath"`
	Module     string `json:"module"`
	Param      string `json:"param"`
}

type StopFlow struct {
	InstanceId string `json:"instanceId"`
}

type Result struct {
	InstanceId string `json:"instanceId"`
	Success    bool   `json:"success"`
	Message    string `json:"message"`
}
