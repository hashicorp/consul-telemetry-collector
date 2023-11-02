#!/usr/bin/env bash
# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0


SHA=${1}
if [[ -z $SHA ]]; then
  echo "Missing build sha"
  echo "Usage: trigger-promotion.bash SHA"
fi

set -euo pipefail
IFS=$'\n\t'

source .release/varz.bash

bob trigger-promotion \
  --product-name=${PRODUCT} \
  --org=${ORG} \
  --repo=${REPO} \
  --branch=${BRANCH} \
  --product-version="${PRODUCT_VERSION}" \
  --environment=${ENVIRONMENT} \
  --slack-channel=${SLACK_CHANNEL} \
  --sha="${SHA}" production
