-include .makefiles/Makefile
-include .makefiles/pkg/protobuf/v1/Makefile
-include .makefiles/pkg/go/v1/Makefile

.makefiles/%:
	@curl -sfL https://makefiles.dev/v1 | bash /dev/stdin "$@"
