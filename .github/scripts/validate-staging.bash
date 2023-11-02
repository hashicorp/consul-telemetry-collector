#!/usr/bin/env bash
# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0


SHA=${1}
if [[ -z $SHA ]]; then
  echo "Missing build sha"
  echo "Usage: trigger-promotion.bash SHA"
fi

set -euo pipefail

source .release/varz.bash

CRT_STAGING_REGISTRY="crt-core-staging-docker-local.artifactory.hashicorp.engineering"
TARGET=release-default

TAG="${CRT_STAGING_REGISTRY}/${REPO}/${TARGET}:${PRODUCT_VERSION}_${SHA}"
echo pulling $TAG
docker pull $TAG
