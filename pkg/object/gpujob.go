package object

const GpuJobEtcdPrefix = "/apis/gpuJob/"

type GpuJob struct {
	TypeMeta   `json:",inline" yaml:",inline"`
	ObjectMeta `json:"metadata" yaml:"metadata"`
	Spec       GpuJobSpec   `json:"spec" yaml:"spec"`
	Status     GpuJobStatus `json:"status" yaml:"status"`
}

type GpuJobSpec struct {
	Filename string `json:"filename" yaml:"filename"`
}

type GpuJobStatus string

const (
	JobCreating  GpuJobStatus = "Creating"
	JobCreated   GpuJobStatus = "Created"
	JobWaiting   GpuJobStatus = "Waiting"
	JobRunning   GpuJobStatus = "Running"
	JobSucceeded GpuJobStatus = "Succeeded"
	JobFailed    GpuJobStatus = "Failed"
)
