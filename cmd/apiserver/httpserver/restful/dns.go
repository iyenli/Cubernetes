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

func GetDns(ctx *gin.Context) {
	getObj(ctx, object.DnsEtcdPrefix+ctx.Param("uid"))
}

func GetDnses(ctx *gin.Context) {
	getObjs(ctx, object.DnsEtcdPrefix)
}

func PostDns(ctx *gin.Context) {
	dns := object.Dns{}
	err := ctx.BindJSON(&dns)
	if err != nil {
		utils.ParseFail(ctx)
		return
	}
	if dns.Name == "" {
		utils.BadRequest(ctx)
		return
	}
	dns.UID = uuid.New().String()
	buf, _ := json.Marshal(dns)
	err = etcdrw.PutObj(object.DnsEtcdPrefix+dns.UID, string(buf))
	if err != nil {
		utils.ServerError(ctx)
		return
	}
	ctx.JSON(http.StatusOK, dns)
}

func PutDns(ctx *gin.Context) {
	newDns := object.Dns{}
	err := ctx.BindJSON(&newDns)
	if err != nil {
		utils.ParseFail(ctx)
		return
	}

	if newDns.UID != ctx.Param("uid") {
		utils.BadRequest(ctx)
		return
	}

	oldBuf, err := etcdrw.GetObj(object.DnsEtcdPrefix + newDns.UID)
	if err != nil {
		utils.ServerError(ctx)
		return
	}
	if oldBuf == nil {
		utils.NotFound(ctx)
		return
	}

	newBuf, _ := json.Marshal(newDns)
	err = etcdrw.PutObj(object.DnsEtcdPrefix+newDns.UID, string(newBuf))
	if err != nil {
		utils.ServerError(ctx)
		return
	}

	ctx.Header("Content-Type", "application/json")
	ctx.String(http.StatusOK, string(newBuf))
}

func DelDns(ctx *gin.Context) {
	delObj(ctx, object.DnsEtcdPrefix+ctx.Param("uid"))
}

func SelectDnses(ctx *gin.Context) {
	var selectors map[string]string
	err := ctx.BindJSON(&selectors)
	if err != nil {
		utils.ParseFail(ctx)
		return
	}

	if len(selectors) == 0 {
		getObjs(ctx, object.DnsEtcdPrefix)
		return
	}

	selectObjs(ctx, object.DnsEtcdPrefix, func(str []byte) bool {
		var dns object.Dns
		err = json.Unmarshal(str, &dns)
		if err != nil {
			return false
		}

		for key, val := range selectors {
			v := dns.Labels[key]
			if v != val {
				return false
			}
		}
		return true
	})
}
