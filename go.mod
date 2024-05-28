module github.com/fabien-marty/slog-helpers

go 1.21.7

require (
	github.com/mattn/go-isatty v0.0.20
	github.com/ztrue/tracerr v0.4.0
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

require (
	github.com/stretchr/testify v1.9.0
	github.com/vlad-tokarev/sloggcp v0.0.0-20230820053939-1b7dbb8c7b58
	golang.org/x/sys v0.6.0 // indirect
)

replace github.com/ztrue/tracerr => github.com/fabien-marty/tracerr v0.0.0-20240521193902-a136106672ba
