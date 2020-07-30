package envlibs

import (
	"context"
	"encoding/json"
	"fmt"
	"github.wdf.sap.corp/Eureka/envlibs/util"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
)

func GetEnvs() (envs map[string]string, err error) {
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
	log.Println("[envlibs]env  return", string(bs))
	if err != nil {
		return
	}
	envs = make(map[string]string)
	err = json.Unmarshal(bs, &envs)
	if err != nil {
		return
	}
	log.Printf("[envlibs]env  return unmarshal %+v\n", envs)
	return
}

func GetMatch(cluster, namespace string, envs map[string]string) (stage string, found bool) {
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
		log.Println("[envlibs] get env failed", err)
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
		fs = fs[1:(len(fs) - 1)]
		if len(fs) == 0 {
			return
		}
		m[strings.Join(fs, "-")] = apath
		return nil
	})
	return
}

var cm sync.Map

func GetAllExtraValuesFilesWithCache(owner, repo, version, root, harborUser, harborPwd string) (m map[string]string, err error) {
	k := fmt.Sprintf("%s:%s:%s", owner, repo, version)
	v, ok := cm.Load(k)
	if ok {
		m = v.(map[string]string)
		return
	}
	path, tgzPath, err := FetchPkg(owner, repo, version, root, harborUser, harborPwd)
	if err != nil {
		return
	}
	m = GetAllExtraValuesFiles(tgzPath)
	util.Exec(context.Background(), "rm -rf "+path)
	cm.Store(k, m)
	return
}

func GetAllExtraValuesFilesWithCallback(owner, repo, version, root, harborUser, harborPwd string, ab bool) (m map[string]string, fn func(), err error) {
	path, tgzPath, err := FetchPkg(owner, repo, version, root, harborUser, harborPwd)
	if err != nil {
		return
	}
	m = GetAllExtraValuesFiles(tgzPath)
	if ab {
		for k := range m {
			m[k] = tgzPath + "/" + m[k]
		}
	}
	fn = func() {
		util.Exec(context.Background(), "rm -rf "+path)
	}
	return
}

func FetchPkg(owner, repo, version, root, harborUser, harborPwd string) (path string, tgzPath string, err error) {
	for {
		path = root + fmt.Sprintf("/%s%s%s", owner, repo, version) + util.RandomString(5)
		_, ferr := os.Stat(path)
		if ferr != nil && os.IsNotExist(ferr) {
			rs := util.Exec(context.Background(), "mkdir -p "+path)
			if rs.Error != nil {
				err = rs.Error
				return
			}
			break
		}
	}
	uc := util.NewExecCmd(path)
	ho := strings.ToLower(owner)
	ctx := context.Background()
	sr := uc.Exec(ctx, fmt.Sprintf("curl https://harbor.eurekacloud.io/chartrepo/%s/charts/%s-%s.tgz -u %s:%s  --output %s-%s.tgz", ho, repo, version, harborUser, harborPwd, repo, version))
	if sr.Error != nil {
		log.Println("error: repo fetch", sr.Stderr)
		err = sr.Error
		return
	}
	rs := uc.Exec(ctx, fmt.Sprintf("tar -zxvf %s-%s.tgz", repo, version))
	if rs.Error != nil {
		log.Println("error: tgz extract", rs.Stderr)
		err = rs.Error
		return
	}
	tgzPath = path + "/" + repo
	return
}

func GetValues(owner, repo, version, root, cluster, namespace, harborUser, harborPwd, stage string) (rs []string, err error) {
	m, err := GetAllExtraValuesFilesWithCache(owner, repo, version, root, harborUser, harborPwd)
	if err != nil {
		return
	}
	if r, ok := m[stage]; ok {
		rs = append(rs, r)
	}
	if r, ok := m[namespace+"-"+cluster]; ok {
		rs = append(rs, r)
	}
	return
}
