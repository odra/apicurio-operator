SHELL = /bin/bash
REG = quay.io
ORG = integreatly
IMAGE = apicurio-operator
TAG = latest
RESOURCES_DIR = ./res
TMPL_VERSIONS = '0.2.22.Final,0.2.21.Final,0.2.20.Final'
DEPLOY_DIR = deploy
OUT_STATIC_DIR = build/_output
OUTPUT_BIN_NAME = ${IMAGE}
TARGET_BIN = cmd/manager/main.go
NS = apicurio
TEST_FOLDER = ./test/e2e
TEST_POD_NAME = apicurio-operator-test
KC_HOST =
APPS_HOST =

.PHONY: setup/dep
setup/dep:
	@echo Installing golang dependencies
	@go get golang.org/x/sys/unix
	@go get golang.org/x/crypto/ssh/terminal
	@go get -u github.com/gobuffalo/packr/packr
	@echo Installing dep
	@curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh
	@echo setup complete

.PHONY: setup/travis
setup/travis:
	@echo Installing Operator SDK
	@curl -Lo operator-sdk https://github.com/operator-framework/operator-sdk/releases/download/v0.1.0/operator-sdk-v0.1.0-x86_64-linux-gnu && chmod +x operator-sdk && sudo mv operator-sdk /usr/local/bin/

.PHONY: dep/ensure
dep/ensure:
	@dep ensure -v

.PHONY: code/run
code/run:
	@operator-sdk up local --namespace=${NAMESPACE}

.PHONY: code/compile
code/compile:
	@packr
	@go build -o ${OUTPUT_BIN_NAME} ${TARGET_BIN}
	@packr clean
	@rm ${OUTPUT_BIN_NAME}

.PHONY: code/gen
code/gen:
	@operator-sdk generate k8s

.PHONY: code/check
code/check:
	@diff -u <(echo -n) <(gofmt -d `find . -type f -name '*.go' -not -path "./vendor/*"`)

.PHONY: code/fix
code/fix:
	@gofmt -w `find . -type f -name '*.go' -not -path "./vendor/*"`

hack/templates:
	./hack/get-template.sh -v ${TMPL_VERSIONS}

.PHONY: image/build
image/build: code/compile
	@packr
	@operator-sdk build ${REG}/${ORG}/${IMAGE}:${TAG}
	@packr clean

.PHONY: image/push
image/push:
	docker push ${REG}/${ORG}/${IMAGE}:${TAG}

.PHONY: image/build/push
image/build/push: image/build
	@docker push ${REG}/${ORG}/${IMAGE}:${TAG}

.PHONY: test/unit
test/unit:
	@go test -v -race -cover ./pkg/...

.PHONY: test/e2e/prepare
test/e2e/prepare:
	@kubectl create secret generic apicurio-operator-test-env --from-literal="apicurio-apps-host=${APPS_HOST}" --from-literal="apicurio-kc-host=${KC_HOST}" -n ${NS}

.PHONY: test/e2e/local
test/e2e/local: image/build/test image/push
	@operator-sdk test local ${TEST_FOLDER} --go-test-flags "-v"

.PHONY: test/e2e/clear
test/e2e/clear:
	@kubectl delete secret/apicurio-operator-test-env -n ${NS}

.PHONY: test/e2e/cluster
test/e2e/cluster: image/build/test image/push
	@kubectl apply -f deploy/test-pod.yaml -n ${NS}
	${SHELL} ./scripts/stream-pod ${TEST_POD_NAME} ${NS}

.PHONY: image/build/test
image/build/test:
	@packr
	@operator-sdk build --enable-tests ${REG}/${ORG}/${IMAGE}:${TAG}
	@packr clean

.PHONY: cluster/prepare
cluster/prepare:
	@kubectl create namespace ${NS}  || true
	@kubectl apply -f ${DEPLOY_DIR}/role.yaml -n ${NS}
	@kubectl apply -f ${DEPLOY_DIR}/role_binding.yaml -n ${NS}
	@kubectl apply -f ${DEPLOY_DIR}/service_account.yaml -n ${NS}
	@kubectl apply -f ${DEPLOY_DIR}/crds/integreatly_v1alpha1_apicurio_crd.yaml -n ${NS}

.PHONY: cluster/deploy/stanalone
cluster/deploy/standalone:
	@kubectl apply -f ${DEPLOY_DIR}/crds/integreatly_v1alpha1_apicurio_cr.full.yaml -n ${NS}
	@kubectl apply -f ${DEPLOY_DIR}/operator.yaml -n ${NS}

.PHONY: cluster/deploy/external-auth
cluster/deploy/external-auth:
	@kubectl apply -f ${DEPLOY_DIR}/crds/integreatly_v1alpha1_apicurio_cr.external_kc.yaml -n ${NS}
	@kubectl apply -f ${DEPLOY_DIR}/operator.yaml -n ${NS}

.PHONY: cluster/clean
cluster/clean:
	@kubectl delete apicurio -l 'app=apicurio' -n ${NS}
	@kubectl delete deployment/apicurio-operator -n ${NS}
