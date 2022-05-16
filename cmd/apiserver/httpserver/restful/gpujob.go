package restful

import (
	"Cubernetes/cmd/apiserver/httpserver/utils"
	"Cubernetes/pkg/object"
	"Cubernetes/pkg/utils/etcdrw"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
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

	if newJob.UID != ctx.Param("uid") {
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
	delObj(ctx, object.GpuJobEtcdPrefix+ctx.Param("uid"))
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
