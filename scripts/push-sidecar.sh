#!/usr/bin/env bash
docker tag sidecar ${KOKI_SIDECAR_IMAGE}
docker push ${KOKI_SIDECAR_IMAGE}
