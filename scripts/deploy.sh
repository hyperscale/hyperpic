#!/bin/bash

set -e

PROJECT_WORKSPACE="$(dirname $0)/.."
PROJECT_WORKSPACE="$(cd $PROJECT_WORKSPACE; pwd)"

DOCKER_ORGANISATION="hypersclae"
DOCKER_REPO="hyperpic"

echo "Config:"
if [ -z "$TRAVIS_TAG" ]; then
    CI_BUILD_VERSION=$(git describe --match 'v[0-9]*' --dirty='-dev' --always)
else
    CI_BUILD_VERSION="${TRAVIS_TAG#v}"
    DOCKER_TAG="$CI_BUILD_VERSION"
fi

if [ -z "$TRAVIS_COMMIT"]; then
    CI_BUILD_COMMIT=$(git rev-parse HEAD)
else
    CI_BUILD_COMMIT="$TRAVIS_COMMIT"
fi

CI_BUILD_URL=$(git config --get remote.origin.url)
CI_BUILD_DATE=$(date +%Y-%m-%dT%T%z)

if [ -z "$TRAVIS_BRANCH" ]; then
    CI_BUILD_BRANCH=$(git rev-parse --abbrev-ref HEAD)
else
    CI_BUILD_BRANCH="$TRAVIS_BRANCH"
fi

if [ "$CI_BUILD_BRANCH" == "develop" ]; then
    DOCKER_TAG="dev"
fi

echo "  Docker Tag: $DOCKER_TAG"
echo "  Version: $CI_BUILD_VERSION"
echo "  VCS URL: $CI_BUILD_URL"
echo "  VCS Ref: $CI_BUILD_COMMIT"
echo "  VCS Branch: $CI_BUILD_BRANCH"
echo "  Build Date: $CI_BUILD_DATE"
echo "  Workspace: $PROJECT_WORKSPACE"
echo ""

echo "Building $DOCKER_ORGANISATION/$DOCKER_REPO..."
docker build --rm \
    --build-arg "VERSION=$CI_BUILD_VERSION" \
    --build-arg "VCS_URL=$CI_BUILD_URL" \
    --build-arg "VCS_REF=$CI_BUILD_COMMIT" \
    --build-arg "BUILD_DATE=$CI_BUILD_DATE" \
    -f "$PROJECT_WORKSPACE/Dockerfile" \
    -t "$DOCKER_ORGANISATION/$DOCKER_REPO:$DOCKER_TAG" \
    "$PROJECT_WORKSPACE"

# tagging latest only master branch
if [ "$CI_BUILD_BRANCH" == "master" ]; then
    echo "Tagging $DOCKER_ORGANISATION/$DOCKER_REPO:latest.."
    docker tag "$DOCKER_ORGANISATION/$DOCKER_REPO" "$DOCKER_ORGANISATION/$DOCKER_REPO:latest"
fi

# pushing only in CI mode
if [ "$CI" == "true" ]; then
    echo "Pushing $DOCKER_ORGANISATION/$DOCKER_REPO:$DOCKER_TAG..."
    docker push "$DOCKER_ORGANISATION/$DOCKER_REPO:$DOCKER_TAG"

    if [ "$CI_BUILD_BRANCH" == "master" ]; then
        echo "Pushing $DOCKER_ORGANISATION/$DOCKER_REPO:latest..."
        docker push "$DOCKER_ORGANISATION/$DOCKER_REPO:latest"
    fi
fi
