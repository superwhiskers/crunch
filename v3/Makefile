#
# crunch - utilities for taking bytes out of things
# Copyright (c) 2019-2020 superwhiskers <whiskerdev@protonmail.com>
#
# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at https://mozilla.org/MPL/2.0/.
#

ifeq ($(origin .RECIPEPREFIX), undefined)
	$(error your make version does not support .RECIPEPREFIX. please use a newer version of make)
endif

.RECIPEPREFIX = >
.ONESHELL:
MAKEFLAGS += --warn-undefined-variables
MAKEFLAGS += --no-builtin-rules

# get path of makefile regardless of directory
MAKEFILE_DIR := $(realpath $(dir $(lastword $(MAKEFILE_LIST))))

.PHONY: all
all: *.generated.go

.PHONY: clean
clean:
> @cd $(MAKEFILE_DIR)
> @rm -f *.generated.go ../generation/generation

$(MAKEFILE_DIR)/generation/generation:
> @cd $(MAKEFILE_DIR)/../generation
> @go build

*.generated.go: $(MAKEFILE_DIR)/generation/generation
> @cd $(MAKEFILE_DIR)
> @../generation/generation -generation "complex"
