package gpuserver

func GetJobDockerName(jobUID string) string {
	return "Job-" + jobUID
}
