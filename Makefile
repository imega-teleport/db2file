CWD=/go/src/github.com/imega-teleport/db2file

build: test
	@mkdir $(CURDIR)/test
	@docker run --rm \
		-v $(CURDIR):$(CWD) \
		-w $(CWD) \
		-e GOOS=linux \
		-e GOARCH=amd64 \
		-e CGO_ENABLED=0 \
		golang:1.8-alpine \
		sh -c 'go build -v -o db2file && ./db2file -db test_teleport -path $(CWD)/test'
	ls -l $(CURDIR)/test

db:
	@touch $(CURDIR)/mysql.log
	@docker run -d \
		--name "server_db" \
		-v $(CURDIR)/sql/cnf:/etc/mysql/conf.d \
		-v $(CURDIR)/mysql.log:/var/log/mysql/mysql.log \
		imega/mysql
	@docker run --rm \
		-v $(CURDIR)/sql:/sql \
		--link server_db:s \
		imega/mysql-client \
		mysql --host=s -e "source /sql/schema.sql"
	@docker run --rm \
		-v $(CURDIR)/sql:/sql \
		--link server_db:s \
		imega/mysql-client \
		mysql --host=s --database=test_teleport -e "source /sql/dump.sql"

clean:
	@-docker stop server_db
	@-docker rm -fv server_db


test: clean db
	@docker run --rm \
		-e DB_USER=root \
		-e DB_PASS=1 \
		-e DB_HOST="server_db:3306" \
		-e DB_NAME=test_teleport \
		-v $(CURDIR):$(CWD) \
		--link server_db:server_db \
		-w $(CWD) \
		golang:1.8-alpine sh -c "go test github.com/imega-teleport/db2file/integration"

dep:
	@docker run --rm \
		-v $(CURDIR):$(CWD) \
		-w $(CWD) \
		golang:1.8-alpine sh -c 'apk add --update git && go get -u github.com/golang/dep/cmd/dep && dep ensure -v'
