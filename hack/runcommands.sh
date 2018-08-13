#!/bin/bash

set -x

operator-sdk generate k8s

operator-sdk build docker.io/surajnarwade/website-operator

docker push docker.io/surajnarwade/website-operator