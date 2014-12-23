# cool server live reload script from
# https://medium.com/@olebedev/live-code-reloading-for-golang-web-projects-in-19-lines-8b2e8777b1ea
PID = /tmp/rklotz.pid

serve:
	@make restart 
	@fswatch -o -r . | xargs -n1 -I{}  make restart || make kill

kill:
	@echo ""
	@echo ""
	@echo ""
	@echo "Trying to kill old instance..."
	@kill `cat $(PID)` || true

stuff:
	@echo ""

restart:
	@make kill
	@make stuff
	@go run main.go --env dev & echo $$! > $(PID)

.PHONY: serve restart kill stuff
