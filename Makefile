dev: local.start.deps
	nix-shell

local.start.deps:
	docker compose -f docker-compose.dev.yml up -d

local.stop.deps:
	docker compose -f docker-compose.dev.yml down



local.run:
	go run .

test.mail:
	swaks --from=decima@zeus --to=do@localhost --server=localhost:2525 --data='Subject: success mail\n\nThis is a test mail.'

test.mail.full:
	swaks --from=$(from) --to=$(to) --server=localhost:2525 --data='Subject: failure mail\n\nThis is a test mail.'


docker.release:
	docker build -t decima/mailhook:latest -t decima/mailhook:$(version) .
	docker push decima/mailhook:latest
	docker push decima/mailhook:$(version)
	git tag -a $(version) -m "Release version $(version)"
	git push --tags
