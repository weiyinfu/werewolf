backend:
	# 构建后端
	go build
front:
	# 构建前端js
	vue build index.vue
format:
	find . -name '*.go' | grep -Ev 'vendor|thrift_gen' | xargs goimports -w


build:format backend front
	# 同时构建前端和后端
linux:
	export CGO_ENABLED=0;export GOOS=linux;export GOARCH=amd64;go build

upload:linux
	# 上传到服务器,服务器一定是linux
	rsync -r --progress werewolf dist/ tencent:~/app/Werewolf
cloc:
	# 统计代码行数
	cloc --exclude-dir=node_modules,dist .
pro:
	vue-cli-service build --mode production
run:
	export GIN_MODE=release;./werewolf