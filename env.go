package envlibs

import (
	"encoding/json"
	"fmt"
	"github.wdf.sap.corp/Eureka/envlibs/util"
	"io/ioutil"
	"net/http"
	"strings"
)

func GetEnvs() (envs map[string]string, err error){
	req, err := http.NewRequest("GET", "http://dev-info/env", nil)
	if err != nil {
	    return
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
	    return
	}
	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("code not match %d", resp.StatusCode)
		return
	}
	defer resp.Body.Close()
	bs, err := ioutil.ReadAll(resp.Body)
	if err != nil {
	    return
	}
	envs = make(map[string]string)
	err = json.Unmarshal(bs, &envs)
	if err != nil {
	    return
	}
	return
}

func  GetMatch(cluster, namespace string, envs map[string]string) (stage string, found bool){
	k := fmt.Sprintf("%s/%s", cluster, namespace)
	if v, ok := envs[k]; ok {
		found = true
		stage = v
		return
	}
	k = fmt.Sprintf("%s/*", cluster)
	if v, ok := envs[k]; ok {
		found = true
		stage = v
		return
	}
	return
}

func GetOriginMatch(cluster, namespace string) (stage string, found bool, err error) {
	r, err := GetEnvs()
	if err != nil {
	    return
	}
	stage, found = GetMatch(cluster, namespace, r)
	return
}

func GetAllExtraValuesFiles(path string) (m map[string]string) {
	m = make(map[string]string)
	_ = util.IterFiles(path, false, func(apath string) interface{} {
		return nil
	}, func(level int32, path string, apath string, rootElem interface{}) (err error) {
		if level != 0 {
			return
		}
		if !strings.HasPrefix(apath, "values.") {
			return
		}
		if !strings.HasSuffix(apath, ".yaml") {
			return
		}
		fs := strings.Split(apath, ".")
		fs = fs[1:(len(fs)-1)]
		if len(fs) == 0 {
			return
		}
		m[strings.Join(fs, "-")] = apath
		return nil
	})
	return
}

