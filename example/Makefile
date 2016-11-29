all: timestamp build
	@echo

build:
	@echo
	@echo "* Generate assets.go"
	go generate

timestamp:
	@echo
	@echo "* Ensure timestamp of local assets"
	find -type f -exec touch -t `date --date=2016-11-21T01:25:30+0800 +%Y%m%d%H%M.%S` {} \;

.PHONY: all build timestamp
