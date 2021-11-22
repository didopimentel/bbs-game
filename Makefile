GOPATH ?= ~/go
GO ?= go
GORUN ?= $(GO) run
GOIMPORTS ?= $(GORUN) golang.org/x/tools/cmd/goimports
PG_ADDR ?= 'postgres://ps_user:ps_password@localhost:7002/bbs-game?sslmode=disable'

################################################################################
## Go Tools
################################################################################

.PHONY: setup
setup:
	@echo "==> Setup: installing tools"
	go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@v4.14.1

################################################################################
## Migration and persistence make targets
################################################################################

.PHONY: migrations/up migrations/down migrations/down/yes db/dump-schema migrations/create

migrations/up:
	@command -v migrate >/dev/null 2>&1 || { echo >&2 "Setup requires migrate but it's not installed.  Aborting."; exit 1; }
	migrate -source="file:$$PWD/gateways/persistence/migrations" -database $(PG_ADDR) up
	make db/dump-schema

migrations/down:
	@command -v migrate >/dev/null 2>&1 || { echo >&2 "Setup requires migrate but it's not installed.  Aborting."; exit 1; }
	migrate -source="file:$$PWD/gateways/persistence/migrations" -database $(PG_ADDR) down

db/dump-schema:
	if [ -z "$(CI)" ]; then \
		pg_dump $(PG_ADDR) -sO --no-comments --no-tablespaces | sed -e '/^--/d' | cat -s > gateways/persistence/schema.sql; \
	fi

################################################################################
## Deps targets
################################################################################

.PHONY: deps undeps

deps:
	docker-compose up -d
	until docker exec bbs-game-postgres pg_isready; do echo 'Waiting for Postgres...' && sleep 1; done
	make migrations/up

undeps:
	docker-compose down

################################################################################
## Linters and formatters
################################################################################

.PHONY: goimports

SOURCES := $(shell \
	find . -name '*.go' | \
	grep -Ev './(proto|protogen|third_party|vendor|dal|.history)/' | \
	xargs)

goimports:
	$(GOIMPORTS) -w $(SOURCES)
