package restful

import (
	"Cubernetes/cmd/apiserver/httpserver/utils"
	cubeconfig "Cubernetes/config"
	"Cubernetes/pkg/object"
	"Cubernetes/pkg/utils/etcdrw"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"os"
	"path"
)

func GetGpuJob(ctx *gin.Context) {
	getObj(ctx, object.GpuJobEtcdPrefix+ctx.Param("uid"))
}

func GetGpuJobs(ctx *gin.Context) {
	getObjs(ctx, object.GpuJobEtcdPrefix)
}

func PostGpuJob(ctx *gin.Context) {
	job := object.GpuJob{}
	err := ctx.BindJSON(&job)
	if err != nil {
		utils.ParseFail(ctx)
		return
	}
	if job.Name == "" {
		utils.BadRequest(ctx)
		return
	}

	job.UID = uuid.New().String()
	job.Status.Phase = object.JobCreating

	buf, _ := json.Marshal(job)
	err = etcdrw.PutObj(object.GpuJobEtcdPrefix+job.UID, string(buf))
	if err != nil {
		utils.ServerError(ctx)
		return
	}
	ctx.JSON(http.StatusOK, job)
}

func PutGpuJob(ctx *gin.Context) {
	newJob := object.GpuJob{}
	err := ctx.BindJSON(&newJob)
	if err != nil {
		utils.ParseFail(ctx)
		return
	}

	if newJob.UID != ctx.Param("uid") || newJob.Status.Phase == object.JobCreating {
		utils.BadRequest(ctx)
		return
	}

	oldBuf, err := etcdrw.GetObj(object.GpuJobEtcdPrefix + newJob.UID)
	if err != nil {
		utils.ServerError(ctx)
		return
	}
	if oldBuf == nil {
		utils.NotFound(ctx)
		return
	}

	newBuf, _ := json.Marshal(newJob)
	err = etcdrw.PutObj(object.GpuJobEtcdPrefix+newJob.UID, string(newBuf))
	if err != nil {
		utils.ServerError(ctx)
		return
	}

	ctx.Header("Content-Type", "application/json")
	ctx.String(http.StatusOK, string(newBuf))
}

func DelGpuJob(ctx *gin.Context) {
	key := object.GpuJobEtcdPrefix + ctx.Param("uid")
	oldBuf, err := etcdrw.GetObj(key)
	if err != nil {
		utils.ServerError(ctx)
		return
	}
	if oldBuf == nil {
		utils.NotFound(ctx)
		return
	}

	err = etcdrw.DelObj(key)
	if err != nil {
		utils.ServerError(ctx)
		return
	}

	filename := path.Join(cubeconfig.JobFileDir, ctx.Param("uid"))
	_ = os.RemoveAll(filename + ".tar.gz")
	_ = os.RemoveAll(filename + ".out")

	ctx.String(http.StatusOK, "deleted")
}

func SelectGpuJobs(ctx *gin.Context) {
	var selectors map[string]string
	err := ctx.BindJSON(&selectors)
	if err != nil {
		utils.ParseFail(ctx)
		return
	}

	if len(selectors) == 0 {
		getObjs(ctx, object.GpuJobEtcdPrefix)
		return
	}

	selectObjs(ctx, object.GpuJobEtcdPrefix, func(str []byte) bool {
		var job object.GpuJob
		err = json.Unmarshal(str, &job)
		if err != nil {
			return false
		}
		for key, val := range selectors {
			v := job.Labels[key]
			if v != val {
				return false
			}
		}
		return true
	})
}
