#!/usr/bin/env bash
docker tag controller ${KOKI_CONTROLLER_IMAGE}
docker push ${KOKI_CONTROLLER_IMAGE}
