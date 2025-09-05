#!/usr/bin/env bash

set -e

# extract each bash `run` step from all workflows and run them through `bash -n` for
# validity and shellcheck for safety. requires `yq` and `shellcheck` in $PATH

IFS="
"

find .github/workflows -type f -name '*y*ml' -print \
  | while read -r YAMLFILE
    do

      echo "Checking $YAMLFILE"

      for n in $(yq '.. | select(has("steps")) | .steps[] | select(.run != null) | .name' < "$YAMLFILE")
      do
        # do not error out while in this inner loop: we want to find all problems
        set +e

        # prepare a safe location to create a temporary script for evaluation
        tmpscript=$(mktemp)

        # extract the jobs.*.steps[].run element into the temp file
        yq ".. | select(has(\"steps\")) | .steps[] | select(.name == \"$n\") | .run" < "$YAMLFILE" > "$tmpscript"

        # run the checks
        bash -n "$tmpscript" || echo "the script in $n did not pass the validity check"
        shellcheck --shell bash -S warning "$tmpscript" || echo "the script in $n did not pass shellcheck"

        rm "$tmpscript"
        set -e
      done
    done