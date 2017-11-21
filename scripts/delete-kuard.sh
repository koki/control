#!/usr/bin/env bash
kubectl delete pod kuard
kubectl delete cm kuard
kubectl delete deployment koki-controller
