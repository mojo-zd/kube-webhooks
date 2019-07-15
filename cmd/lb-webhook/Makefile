ROOT=../..

PROG = lb-webhook
REPO?=registry.cn-hangzhou.aliyuncs.com/mojo
GIT_COMMIT=`git rev-parse --short HEAD`
RELEASE ?= `git symbolic-ref --short -q HEAD`

IMG=${REPO}/${PROG}:${RELEASE}

default: clean
	CGO_ENABLED=0 go build -a -tags netgo -installsuffix netgo -ldflags "-s -w -X github.com/mojo-zd/kube-webhooks/version.gitCommit=${GIT_COMMIT}" -o ${PROG}

docker: clean
	docker build --label gitCommit=${GIT_COMMIT} --build-arg gitCommit=${GIT_COMMIT} -t ${IMG} -f Dockerfile ${ROOT}
	docker push ${IMG}
	@echo "<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<"
	@echo "image_name: ${IMG}"
	@echo "$$(docker images ${IMG} --format 'image_id: {{ .ID }}')"
	@echo "$$(docker images ${IMG} --format 'created:  {{ .CreatedSince }}')"
	@echo "$$(docker images ${IMG} --format 'size:     {{ .Size }} ')"

clean:
	-@rm ${PROG}

test:
	go test ./test/...

.PHONY: clean