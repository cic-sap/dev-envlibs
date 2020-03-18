package envlibs

import "testing"

func TestFetchPkg(t *testing.T) {
	p, tg, err := FetchPkg("eureka-jenkins-test", "gopm12", "0.0.1-ead8e9248eef394bdeaa2d6ca97ea2873c70bff7", "./output")
	t.Logf("%+v, %+v, %+v", p, tg, err)

}
