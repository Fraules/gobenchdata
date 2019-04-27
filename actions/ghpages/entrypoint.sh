#!/bin/bash
set -e

cd "${GITHUB_WORKSPACE}"

echo '📊 Running benchmarks...'
RUN_OUTPUT="/data/gobenchdata/benchmarks.json"
go test \
  -bench "${GO_BENCHMARKS:-"."}" \
  -benchmem \
  ${GO_BENCHMARK_FLAGS} \
  ${GO_BENCHMARK_PKGS:-"./..."} \
  | gobenchdata --json "${OUTPUT}" -v "${GITHUB_SHA}" -t "ref=${GITHUB_REF}"

echo '📚 Checkout out gh-pages...'
mkdir -p build
cd build
git clone https://${GITHUB_TOKEN}@github.com/${GITHUB_REPOSITORY}.git
git checkout gh-pages

FINAL_OUTPUT="${GO_BENCHMARK_OUT:-"benchmarks.json"}"
if [[ -f "${FINAL_OUTPUT}" ]]; then
  echo '📈 Existing report found - merging...'
  gobenchdata merge "${RUN_OUTPUT}" "${FINAL_OUTPUT}" --json "${FINAL_OUTPUT}"
else
  cp "${RUN_OUTPUT}" "${FINAL_OUTPUT}"
fi

echo '📷 Committing new benchmark data...'
git add .
git commit -m "${GIT_COMMIT_MESSAGE:-"add new benchmark run"}"
git push origin gh-pages
cd ../

echo '🚀 Done!'
