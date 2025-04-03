module github.com/barani129/osp-director-operator

go 1.23.0

toolchain go1.23.5

exclude k8s.io/cluster-bootstrap v0.0.0

require (
	github.com/Masterminds/sprig v2.22.0+incompatible
	github.com/blang/semver v3.5.1+incompatible
	github.com/ghodss/yaml v1.0.1-0.20190212211648-25d852aebe32
	github.com/go-git/go-git/v5 v5.13.2
	github.com/go-logr/logr v1.4.2
	github.com/golang/glog v1.2.4
	github.com/google/uuid v1.6.0
	github.com/gophercloud/gophercloud v1.14.1
	github.com/k8snetworkplumbingwg/network-attachment-definition-client v0.0.0-20200626054723-37f83d1996bc
	github.com/metal3-io/baremetal-operator/apis v0.0.0-20220105105621-0ee9ce37c7bc
	github.com/nmstate/kubernetes-nmstate/api v0.0.0-20230518145043-eefd3b0fcd65
	github.com/onsi/ginkgo/v2 v2.21.0
	github.com/onsi/gomega v1.35.1
	github.com/openshift/cluster-api v0.0.0-20191129101638-b09907ac6668
	github.com/openshift/sriov-network-operator v0.0.0-20201204053545-49045c36efb9
	github.com/openstack-k8s-operators/osp-director-operator v1.2.3
	github.com/pkg/errors v0.9.1
	github.com/spf13/cobra v1.9.1
	github.com/tidwall/gjson v1.18.0
	golang.org/x/crypto v0.33.0
	k8s.io/api v0.32.2
	k8s.io/apimachinery v0.32.2
	k8s.io/client-go v12.0.0+incompatible
	k8s.io/utils v0.0.0-20241210054802-24370beab758
	kubevirt.io/api v1.5.0
	kubevirt.io/client-go v1.5.0
	kubevirt.io/containerized-data-importer-api v1.60.3-0.20241105012228-50fbed985de9
	sigs.k8s.io/controller-runtime v0.16.1
	sigs.k8s.io/yaml v1.4.0
)

replace github.com/openstack-k8s-operators/osp-director-operator => ./

replace k8s.io/kube-openapi => k8s.io/kube-openapi v0.0.0-20250318190949-c8a335a9a2ff

replace kubevirt.io/client-go => github.com/barani129/client-go v0.1.3

replace kubevirt.io/client-go/kubecli => github.com/barani129/client-go/kubecli v0.1.3

require (
	dario.cat/mergo v1.0.0 // indirect
	github.com/ProtonMail/go-crypto v1.1.5 // indirect
	github.com/cloudflare/circl v1.3.7 // indirect
	github.com/cyphar/filepath-securejoin v0.3.6 // indirect
	github.com/emicklei/go-restful/v3 v3.11.0 // indirect
	github.com/evanphx/json-patch/v5 v5.6.0 // indirect
	github.com/fxamacker/cbor/v2 v2.7.0 // indirect
	github.com/go-kit/log v0.2.1 // indirect
	github.com/go-openapi/jsonpointer v0.21.0 // indirect
	github.com/go-openapi/jsonreference v0.20.2 // indirect
	github.com/go-openapi/swag v0.23.0 // indirect
	github.com/go-task/slim-sprig/v3 v3.0.0 // indirect
	github.com/google/gnostic-models v0.6.9 // indirect
	github.com/google/pprof v0.0.0-20241029153458-d1b30febd7db // indirect
	github.com/josharian/intern v1.0.0 // indirect
	github.com/kubernetes-csi/external-snapshotter/client/v4 v4.2.0 // indirect
	github.com/mailru/easyjson v0.7.7 // indirect
	github.com/munnerz/goautoneg v0.0.0-20191010083416-a7dc8b61c822 // indirect
	github.com/mxk/go-flowrate v0.0.0-20140419014527-cca7078d478f // indirect
	github.com/pjbgf/sha1cd v0.3.2 // indirect
	github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring v0.68.0 // indirect
	github.com/skeema/knownhosts v1.3.0 // indirect
	github.com/x448/float16 v0.8.4 // indirect
	golang.org/x/exp v0.0.0-20240719175910-8a7402abbf56 // indirect
	golang.org/x/mod v0.21.0 // indirect
	golang.org/x/sync v0.11.0 // indirect
	golang.org/x/tools v0.26.0 // indirect
	gopkg.in/evanphx/json-patch.v4 v4.12.0 // indirect
	k8s.io/kube-openapi v0.0.0-20241105132330-32ad38e42d3f // indirect
	kubevirt.io/controller-lifecycle-operator-sdk/api v0.0.0-20220329064328-f3cc58c6ed90 // indirect
	sigs.k8s.io/randfill v1.0.0 // indirect
)

require (
	github.com/Masterminds/goutils v1.1.1 // indirect
	github.com/Masterminds/semver v1.5.0 // indirect
	github.com/Microsoft/go-winio v0.6.1 // indirect
	github.com/ajeddeloh/go-json v0.0.0-20170920214419-6a2fe990e083 // indirect
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/clarketm/json v1.14.1 // indirect
	github.com/coreos/fcct v0.5.0 // indirect
	github.com/coreos/go-json v0.0.0-20211020211907-c63f628265de // indirect
	github.com/coreos/go-semver v0.3.1 // indirect
	github.com/coreos/go-systemd v0.0.0-20190719114852-fd7a80b32e1f // indirect
	github.com/coreos/go-systemd/v22 v22.5.0 // indirect
	github.com/coreos/ign-converter v0.0.0-20201123214124-8dac862888aa // indirect
	github.com/coreos/ignition v0.35.0 // indirect
	github.com/coreos/ignition/v2 v2.14.0 // indirect
	github.com/coreos/vcontext v0.0.0-20211021162308-f1dbbca7bef4 // indirect
	github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc // indirect
	github.com/emirpasic/gods v1.18.1 // indirect
	github.com/fsnotify/fsnotify v1.7.0 // indirect
	github.com/go-git/gcfg v1.5.1-0.20230307220236-3a3c6141e376 // indirect
	github.com/go-git/go-billy/v5 v5.6.2 // indirect
	github.com/go-kit/kit v0.13.0 // indirect
	github.com/go-logfmt/logfmt v0.6.0 // indirect
	github.com/go-logr/zapr v1.2.4 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/golang/groupcache v0.0.0-20210331224755-41bb18bfe9da // indirect
	github.com/golang/mock v1.6.0 // indirect
	github.com/golang/protobuf v1.5.4 // indirect
	github.com/google/go-cmp v0.6.0 // indirect
	github.com/google/gofuzz v1.2.0 // indirect
	github.com/gorilla/websocket v1.5.0 // indirect
	github.com/huandu/xstrings v1.3.2 // indirect
	github.com/imdario/mergo v0.3.16 // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/jbenet/go-context v0.0.0-20150711004518-d14ea06fba99 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/kevinburke/ssh_config v1.2.0 // indirect
	github.com/metal3-io/baremetal-operator/pkg/hardwareutils v0.0.0 // indirect
	github.com/mitchellh/copystructure v1.0.0 // indirect
	github.com/mitchellh/reflectwalk v1.0.1 // indirect
	github.com/moby/spdystream v0.5.0 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/openshift/api v0.0.1 // indirect
	github.com/openshift/client-go v0.0.0-20211209144617-7385dd6338e3 // indirect
	github.com/openshift/custom-resource-status v1.1.2 // indirect
	github.com/openshift/machine-config-operator v0.0.1-0.20230323140348-afb47c916680 // indirect
	github.com/prometheus/client_golang v1.19.1 // indirect
	github.com/prometheus/client_model v0.6.1 // indirect
	github.com/prometheus/common v0.55.0 // indirect
	github.com/prometheus/procfs v0.15.1 // indirect
	github.com/sergi/go-diff v1.3.2-0.20230802210424-5b0b94c5c0d3 // indirect
	github.com/spf13/pflag v1.0.6 // indirect
	github.com/tidwall/match v1.1.1 // indirect
	github.com/tidwall/pretty v1.2.0 // indirect
	github.com/vincent-petithory/dataurl v1.0.0 // indirect
	github.com/xanzy/ssh-agent v0.3.3 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	go.uber.org/zap v1.27.0
	go4.org v0.0.0-20200104003542-c7e774b10ea0 // indirect
	golang.org/x/net v0.34.0 // indirect
	golang.org/x/oauth2 v0.23.0 // indirect
	golang.org/x/sys v0.30.0 // indirect
	golang.org/x/term v0.29.0 // indirect
	golang.org/x/text v0.22.0 // indirect
	golang.org/x/time v0.7.0 // indirect
	gomodules.xyz/jsonpatch/v2 v2.4.0 // indirect
	google.golang.org/protobuf v1.35.1 // indirect
	gopkg.in/inf.v0 v0.9.1 // indirect
	gopkg.in/warnings.v0 v0.1.2 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.1
	k8s.io/apiextensions-apiserver v0.31.0 // indirect
	k8s.io/component-base v0.31.0 // indirect
	k8s.io/klog v1.0.0 // indirect
	k8s.io/klog/v2 v2.130.1
	sigs.k8s.io/json v0.0.0-20241010143419-9aa6b5e7a4b3 // indirect
	sigs.k8s.io/kube-storage-version-migrator v0.0.5
	sigs.k8s.io/structured-merge-diff/v4 v4.6.0 // indirect
)

replace (
	github.com/containerd/containerd => github.com/containerd/containerd v1.5.7
	// dependabot fixes
	github.com/dgrijalva/jwt-go => github.com/golang-jwt/jwt v3.2.1+incompatible
	github.com/irifrance/gini => github.com/go-air/gini v1.0.4

	// dependabot (NULL Pointer Dereference in Kubernetes CSI snapshot-controller)
	github.com/kubernetes-csi/external-snapshotter/v2 => github.com/kubernetes-csi/external-snapshotter/v2 v2.1.3
	github.com/metal3-io/baremetal-operator/apis => github.com/openshift/baremetal-operator/apis v0.0.0-20220127144325-36eec3619228 // release-4.10
	github.com/metal3-io/baremetal-operator/pkg/hardwareutils => github.com/openshift/baremetal-operator/pkg/hardwareutils v0.0.0-20220124032151-aa596d5a5bdd // release-4.10

	// required by client-go, prometheus-operator..
	// Bump to avoid CVE detection with v1.1.22. https://bugzilla.redhat.com/show_bug.cgi?id=1786761
	github.com/miekg/dns => github.com/miekg/dns v1.1.43

	// controller runtime
	github.com/openshift/api => github.com/openshift/api v0.0.1 // release-4.10
	github.com/openshift/client-go => github.com/openshift/client-go v0.0.0-20211209144617-7385dd6338e3 // release-4.10

	// CDI
	github.com/operator-framework/operator-lifecycle-manager => github.com/operator-framework/operator-lifecycle-manager v0.0.0-20190128024246-5eb7ae5bdb7a

	// CDI. Bump to avoid CVE with v0.5.7
	github.com/ulikunitz/xz => github.com/ulikunitz/xz v0.5.10

	google.golang.org/grpc => google.golang.org/grpc v1.26.0

	// pin to v0.24.13
	// required by Microsoft/hcsshim, containers/storage, sriov-network-operator
	// Not used within this Operator.
	// Bump to avoid CVE detection with earlier versions (v1.5.4).
	// * https://bugzilla.redhat.com/show_bug.cgi?id=1899487
	// * https://bugzilla.redhat.com/show_bug.cgi?id=1982681
	// * https://bugzilla.redhat.com/show_bug.cgi?id=2011007
	k8s.io/client-go => k8s.io/client-go v0.32.2
	k8s.io/cloud-provider => k8s.io/cloud-provider v0.24.13
	k8s.io/cluster-bootstrap => k8s.io/cluster-bootstrap v0.24.13
	k8s.io/code-generator => k8s.io/code-generator v0.24.13
	k8s.io/component-base => k8s.io/component-base v0.24.13
	k8s.io/cri-api => k8s.io/cri-api v0.24.13
	k8s.io/csi-translation-lib => k8s.io/csi-translation-lib v0.24.13
	k8s.io/kube-aggregator => k8s.io/kube-aggregator v0.24.13
	k8s.io/kube-controller-manager => k8s.io/kube-controller-manager v0.24.13
	k8s.io/kube-proxy => k8s.io/kube-proxy v0.24.13
	k8s.io/kube-scheduler => k8s.io/kube-scheduler v0.24.13
	k8s.io/kubectl => k8s.io/kubectl v0.24.13
	k8s.io/kubelet => k8s.io/kubelet v0.24.13

	// required by kubernetes-csi/external-snapshotter, kubevirt.io/client-go. Bump to avoid CVE detection with v1.14.0: https://bugzilla.redhat.com/show_bug.cgi?id=1757701
	// Not used within this Operator.

	k8s.io/legacy-cloud-providers => k8s.io/legacy-cloud-providers v0.24.13
	k8s.io/metrics => k8s.io/metrics v0.24.13
	k8s.io/node-api => k8s.io/node-api v0.24.13
	k8s.io/sample-apiserver => k8s.io/sample-apiserver v0.24.13
	k8s.io/sample-cli-plugin => k8s.io/sample-cli-plugin v0.24.13
	k8s.io/sample-controller => k8s.io/sample-controller v0.24.13

	// pinned because no tag supports 1.18 yet
	sigs.k8s.io/structured-merge-diff => sigs.k8s.io/structured-merge-diff v1.0.1-0.20191108220359-b1b620dd3f06

)
