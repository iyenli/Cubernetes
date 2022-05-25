package main

import (
	"Cubernetes/pkg/apiserver/crudobj"
	"Cubernetes/pkg/apiserver/objfile"
	"Cubernetes/pkg/cubenetwork/nodenetwork"
	"Cubernetes/pkg/object"
	"Cubernetes/pkg/utils/sshutils"
	"fmt"
	"golang.org/x/crypto/ssh"
	"log"
	"net"
	"os"
	"path"
	"strings"
	"time"
)

var gpuJobUID string
var job object.GpuJob

const JobDir = "./jobs/"

const sshUser = ""
const sshPwd = ""

// gpujobserver [gpuJobUID] [MasterIP]
func main() {
	if len(os.Args) < 3 {
		log.Fatal("[FATAL] Lack arguments")
	}

	gpuJobUID = os.Args[1]
	nodenetwork.SetMasterIP(os.Args[2])

	var err error
	job, err = crudobj.GetGpuJob(gpuJobUID)
	if err != nil {
		jobFail("[FATAL] Fail to get GpuJob, err: ", err)
	}

	log.Println("[INFO]: Get GPU Job meta success")
	_ = os.MkdirAll(JobDir, 0777)
	err = objfile.GetJobFile(gpuJobUID, JobDir+gpuJobUID+".tar.gz")
	if err != nil {
		jobFail("[FATAL] Fail to get job file, err: ", err)
	}

	log.Println("[INFO]: Get GPU Job file success")
	auth := make([]ssh.AuthMethod, 0)
	auth = append(auth, ssh.Password(sshPwd))

	config := ssh.ClientConfig{
		User: sshUser,
		Auth: auth,
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
	}

	job.Status.Phase = object.JobSubmitting
	_, err = crudobj.UpdateGpuJob(job)
	if err != nil {
		jobFail("[Fatal] Fail to update GpuJob phase, err: ", err)
	}

	// upload job file to sjtu data server
	fileClient, err := ssh.Dial("tcp", "data.hpc.sjtu.edu.cn:22", &config)
	if err != nil {
		jobFail("[FATAL] Fail to establish ssh connection with data server, err: ", err)
	}

	log.Println("[INFO]: ssh connection created success")
	filename := gpuJobUID + ".tar.gz"
	filePath := path.Join(JobDir, filename)
	err = sshutils.UploadFile(fileClient, filePath, filePath, JobDir)
	_ = fileClient.Close()
	if err != nil {
		jobFail("[FATAL] Fail to upload file to data server, err: ", err)
	}

	log.Println("[INFO]: Job file uploaded to remote hpc")
	client, err := ssh.Dial("tcp", "login.hpc.sjtu.edu.cn:22", &config)
	if err != nil {
		jobFail("[FATAL] Fail to establish ssh connection with login server, err: ", err)
	}
	defer func() { _ = client.Close() }()

	cmd := fmt.Sprintf("cd %s && rm -rf ./%s && mkdir %s && tar zxvf %s -C ./%s --strip-components 1 && cd ./%s && sbatch *.slurm", JobDir, gpuJobUID, gpuJobUID, filename, gpuJobUID, gpuJobUID)
	res, err := sshutils.Exec(client, cmd)
	if err != nil {
		jobFail("[FATAL] Fail to extract and commit job, err: ", err)
	}

	log.Println("[INFO]: Job committed successfully")
	log.Println("[INFO]: Waiting with your best patience...")

	results := strings.Split(res, "\n")
	var s, slurmJobId string
	_, err = fmt.Sscanf(results[len(results)-2], "%s %s %s %s", &s, &s, &s, &slurmJobId)
	if err != nil {
		jobFail("[FATAL] Fail to parse Slurm Job Id, err: ", err)
	}

	job.Status.SlurmJobId = slurmJobId
	job.Status.Phase = object.JobWaiting
	_, err = crudobj.UpdateGpuJob(job)
	if err != nil {
		jobFail("[Fatal] Fail to update GpuJob phase, err: ", err)
	}

	for {
		time.Sleep(5 * time.Second)
		res, err := sshutils.Exec(client, "squeue -j "+slurmJobId)
		if err != nil {
			jobFail("[FATAL] Fail to parse Slurm output, err: ", err)
		}
		results := strings.Split(res, "\n")

		var s, phase string
		_, err = fmt.Sscanf(results[len(results)-2], "%s %s %s %s %s %s %s %s", &s, &s, &s, &s, &phase, &s, &s, &s)
		if job.Status.Phase != object.JobRunning && phase == "R" {
			job.Status.Phase = object.JobRunning
			_, err = crudobj.UpdateGpuJob(job)
			if err != nil {
				log.Println("[WARNING] Fail to update job status to running")
			}
		}
		if len(results) == 2 {
			break
		}
	}

	bufOut, err := sshutils.ReadFile(client, path.Join(JobDir, gpuJobUID, slurmJobId+".out"))
	if err != nil {
		jobFail("[FATAL] Fail to read output, err: ", err)
	}

	bufErr, err := sshutils.ReadFile(client, path.Join(JobDir, gpuJobUID, slurmJobId+".err"))
	if err != nil {
		jobFail("[FATAL] Fail to read error output, err: ", err)
	}

	output := "Errors:\n" + string(bufErr) + "\n\nOutputs:\n" + string(bufOut)
	jobSuccess(output)
}

func jobSuccess(output string) {
	err := objfile.PostJobOutput(gpuJobUID, output)
	if err != nil {
		jobFail("[FATAL] Fail to upload output, err: ", err)
	}

	job.Status.Phase = object.JobSucceeded
	_, err = crudobj.UpdateGpuJob(job)
	if err != nil {
		jobFail("[FATAL] Fail to update GpuJob, err: ", err)
	}

	log.Printf("Job UID=%s completed successfully\n", gpuJobUID)
	os.Exit(0)
}

func jobFail(msg string, err error) {
	output := msg + err.Error()
	lerr := objfile.PostJobOutput(gpuJobUID, output)
	if lerr != nil {
		log.Println("[FATAL] Fail to upload output, err: ", lerr)
	}

	job.Status.Phase = object.JobFailed
	_, lerr = crudobj.UpdateGpuJob(job)
	if lerr != nil {
		log.Println("[FATAL] Fail to update GpuJob, err: ", lerr)
	}

	log.Fatal(msg, err)
}
