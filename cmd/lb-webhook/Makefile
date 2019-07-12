ROOT=../..

PROG = lb-webhook
DEV_REPO?=registry.cn-hangzhou.aliyuncs.com/wise2c-dev
RLS_REPO?=registry.cn-hangzhou.aliyuncs.com/wise2c-test
RELEASE ?= `git symbolic-ref --short -q HEAD`
GIT_COMMIT=`git rev-parse --short HEAD`

DEV_IMG=${DEV_REPO}/${PROG}:${RELEASE}
RLS_IMG=${RLS_REPO}/${PROG}:${RELEASE}

default: clean
	CGO_ENABLED=0 go build -a -tags netgo -installsuffix netgo -ldflags "-s -w -X gitee.com/wisecloud/kube-webhooks/version.gitCommit=${GIT_COMMIT}" -o ${PROG}

docker:
	docker build --label gitCommit=${GIT_COMMIT} --build-arg gitCommit=${GIT_COMMIT} -t ${DEV_IMG} -f Dockerfile ${ROOT}
	docker push ${DEV_IMG}

release:
	docker build --label gitCommit=${GIT_COMMIT} --build-arg gitCommit=${GIT_COMMIT} -t ${RLS_IMG} -f Dockerfile.release ${ROOT}
	docker push ${RLS_IMG}
	@echo "<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<"
	@echo "image_name: ${RLS_IMG}"
	@echo "$$(docker images ${RLS_IMG} --format 'image_id: {{ .ID }}')"
	@echo "$$(docker images ${RLS_IMG} --format 'created:  {{ .CreatedSince }}')"
	@echo "$$(docker images ${RLS_IMG} --format 'size:     {{ .Size }} ')"

clean:
	-@rm ${PROG}

test:
	go test ./test/...

.PHONY: clean