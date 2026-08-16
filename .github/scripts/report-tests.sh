#!/bin/bash
# Runs the unit test suite with coverage, then renders a markdown summary
# (job summary + sticky PR comment on failure) via tparse.
set -o pipefail

go test -coverprofile=coverage.out -json -v ./... > test-output.json
TEST_EXIT_CODE=$?

tparse -all -file test-output.json -format markdown > test-summary.md || true

{
  echo ""
  echo "<!-- Sticky Pull Request Comment: test-failure-summary -->"
} >> test-summary.md

cat test-summary.md >> "$GITHUB_STEP_SUMMARY"

exit $TEST_EXIT_CODE
