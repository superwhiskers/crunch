#
# crunch - utilities for taking bytes out of things
# copyright (c) 2019-2020 superwhiskers <whiskerdev@protonmail.com>
#
# this source code form is subject to the terms of the mozilla public
# license, v. 2.0. if a copy of the mpl was not distributed with this
# file, you can obtain one at http://mozilla.org/MPL/2.0/.
#

ifeq ($(origin .RECIPEPREFIX), undefined)
	$(error your make version does not support .RECIPEPREFIX. please use a newer version of make)
endif

.RECIPEPREFIX = >
.ONESHELL:
MAKEFLAGS += --warn-undefined-variables
MAKEFLAGS += --no-builtin-rules

.PHONY: all
all: *.generated.go

.PHONY: clean
clean:
> rm *.generated.go ../generation/generation

../generation/generation:
> cd ../generation
> go build

*.generated.go: ../generation/generation
> ../generation/generation -generation "complex"
