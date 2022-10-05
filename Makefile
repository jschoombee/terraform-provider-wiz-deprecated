GOOS   	   ?= darwin
GOARCH 	   ?= arm64
PLUGINDIR  = ${HOME}/.terraform.d/plugins/terraform.shell.com/deepblue/wiz/0.0.1/$(GOOS)_$(GOARCH)
OUTPUTFILE = terraform-provider-wiz

default: build

# Build the provider
build:
	@go build -o ${OUTPUTFILE} .

# In order to make the provider "globally" available to the system there are a
# few places where we can put the binary.
# See Implied Local Mirror Directories [1] for more information.
# We have decided to place it in the home directory as it's the safest and cleanest
# bet for local development.
#
# The `install` target first builds the binary and then places it in the home
# directory. See the `${PLUGINDIR}` to get the exact location.
#
# This will allow us to test the provider locally without having to publish to
# the public registry, which is the suggested method for In-House Providers [2].
#
# [1]: https://www.terraform.io/cli/config/config-file#implied-local-mirror-directories
# [2]: https://www.terraform.io/language/providers/requirements#in-house-providers
install: build
	@mkdir -p ${PLUGINDIR}
	@mv ${OUTPUTFILE} ${PLUGINDIR}/${OUTPUTFILE}

# Run acceptance tests
.PHONY: testacc
testacc:
	TF_ACC=1 go test ./... -v $(TESTARGS) -timeout 120m

# Lint
lint:
	@gofmt -d .
