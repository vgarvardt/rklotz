# cool server live reload script from
# https://medium.com/@olebedev/live-code-reloading-for-golang-web-projects-in-19-lines-8b2e8777b1ea
PID = /tmp/rklotz.pid

vendor:
	@gb vendor restore
	@bower install

build:
	@echo "Building..."
	@gb build

serve:
	@make restart
	@fswatch -o -r ./src/ -r ./templates/ | xargs -n1 -I{}  make restart || make kill

kill:
	@echo ""
	@echo ""
	@echo ""
	@echo "Trying to kill old instance..."
	@kill `cat $(PID)` || true

restart:
	@make kill
	@make build
	@bin/rklotz --env dev & echo $$! > $(PID)

.PHONY: vendor build serve kill restart
