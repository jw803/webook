.PHONY: docker
docker:
	@rm webook || true
	@GOOS=linux GOARCH=arm go build -tags=k8s -o webook .
	@docker rmi -f flycash/webook:v0.0.1
	@docker build -t flycash/webook-live:v0.0.1 .

.PHONY: mock
mock:
	@mockgen -source=internal/service/user.go -package=svcmocks -destination=internal/service/mocks/user.mock.go
	@mockgen -source=internal/service/article.go -package=svcmocks -destination=internal/service/mocks/article.mock.go
	@mockgen -source=internal/repository/user.go -package=repomocks -destination=internal/repository/mocks/user.mock.go
	@mockgen -source=internal/repository/article/article.go -package=repoarticlemocks -destination=internal/repository/article/mocks/article.mock.go
	@mockgen -source=internal/repository/article/article_author.go -package=repoarticlemocks -destination=internal/repository/article/mocks/article_author.mock.go
	@mockgen -source=internal/repository/article/article_reader.go -package=repoarticlemocks -destination=internal/repository/article/mocks/article_reader.mock.go
	@mockgen -source=internal/repository/dao/user.go -package=daomocks -destination=internal/repository/dao/mocks/user.mock.go
	@mockgen -source=internal/repository/cache/user.go -package=cachemocks -destination=internal/repository/cache/mocks/user.mock.go
	@mockgen -source=internal/service/sms/types.go -package=smsmocks -destination=internal/service/sms/mocks/sms.mock.go
	@mockgen -source=pkg/ratelimit/types.go -package=limitmocks -destination=pkg/ratelimit/mocks/ratelimit.mock.go
	@mockgen -package=redismocks -destination=internal/repository/cache/redismocks/cmdable.mock.go github.com/redis/go-redis/v9 Cmdable
	@go mod tidy