APPTAG="rklotz/dev:latest"
WORKDIR="/go/src/github.com/vgarvardt/rklotz"
VOLSRC="`pwd`/src:/go/src"
VOLBIN="`pwd`/bin:/data/bin"
VOLDB="`pwd`/db:/data/db"
VOLTPL="`pwd`/templates:/data/templates"
VOLSTC="`pwd`/static:/data/static"
VOLUMES=--volume "${VOLSRC}" --volume "${VOLBIN}" --volume "${VOLDB}" --volume "${VOLTPL}" --volume "${VOLSTC}"


init:
	@docker build -f ./Dockerfile.dev -t ${APPTAG} .
	@docker run ${VOLUMES} --workdir "${WORKDIR}" ${APPTAG} glide install
	@docker run ${VOLUMES} --workdir "/data/static" ${APPTAG} bower --allow-root install

cli:
	@docker run --interactive --tty ${VOLUMES} --workdir "${WORKDIR}" ${APPTAG} /bin/bash

build:
	@echo "Building amd64/linux..."
	@docker run \
		${VOLUMES} \
		--workdir "${WORKDIR}" \
		${APPTAG} \
		env GOOS=linux GOARCH=amd64 \
			go build -ldflags "-X github.com/vgarvardt/rklotz/app.version=`<./VERSION`" -v -o /data/bin/rklotz.linux
	@echo "Building amd64/darwin..."
	@docker run \
		${VOLUMES} \
		--workdir "${WORKDIR}" \
		${APPTAG} \
		env GOOS=darwin GOARCH=amd64 \
			go build -ldflags "-X github.com/vgarvardt/rklotz/app.version=`<./VERSION`" -v -o /data/bin/rklotz.darwin

run:
	@docker run --interactive --tty \
		${VOLUMES} \
		--workdir "${WORKDIR}" \
		--publish 8080:8080 \
		--hostname 127.0.0.1 \
		--env-file ./env.dev.txt \
		${APPTAG} \
		go run main.go --root="/data"

rund:
	@docker run --detach \
		${VOLUMES} \
		--workdir "${WORKDIR}" \
		--publish 8080:8080 \
		--hostname 127.0.0.1 \
		--env-file ./env.dev.txt \
		${APPTAG} \
		go run main.go --root="/data"

serve:
	@make restart
	@fswatch -o -r ./src/ -r ./templates/ -r ./static/ | xargs -n1 -I{} make restart || make kill

kill:
	@echo ""
	@echo ""
	@echo ""
	@echo "Trying to kill old instance..."
	@docker ps | grep ${APPTAG} | awk '{print $$1}' | xargs docker stop || true

restart:
	@make kill
	@make rund

test:
	@docker run --tty \
	${VOLUMES} \
	--workdir "${WORKDIR}" \
	--env-file ./env.dev.txt \
	${APPTAG} \
	/bin/bash -c "go list ./... | grep -v /vendor/ | xargs go test"

.PHONY: init cli build run rund serve kill restart test
