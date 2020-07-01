#!/usr/bin/env bash
# exit immediately when a command fails
set -e
# only exit with zero if all commands of the pipeline exit successfully
set -o pipefail
# error on unset variables
set -u
# print each command before executing it
set -x

function runTest()
{
    declare -A map
    map["testsuite-1"]="ginkgo --focus=\"\[Revision\]\" -skip=\"\[Slow\]|\[Serial\]\" -r test"
    map["testsuite-2"]="ginkgo --focus=\"\[CRD\]\" -skip=\"\[Slow\]|\[Serial\]\" -r test"
    map["testsuite-3"]="ginkgo --focus=\"\[CONTROLLER-1\]\" -skip=\"\[Slow\]|\[Serial\]\" -r test"
    map["testsuite-4"]="ginkgo --skip=\"\[Serial\]|\[Slow\]|\[Revision\]|\[CRD\]|\[CONTROLLER-1\]|\[CONTROLLER-2\]\" -r test"
    map["testsuite-5"]="ginkgo --focus=\"\[CONTROLLER-2\]\" -skip=\"\[Slow\]|\[Serial\]\" -r test"
    map["testsuite-6"]="ginkgo --focus=\"\[Serial\]\" -r test"
    map["testsuite-7"]="ginkgo --focus=\"\[Slow\]\" -r test"

    set +u
    set +e

    res="0"

    # shellcheck disable=SC2068
    for key in ${!map[@]}
    do
      if [[ ${1} == $key ]]; then
        eval ${map[$key]}
        sub_res=`echo $?`
        if [ $sub_res != "0" ]; then
          res=$sub_res
        fi
      fi
    done

    set -e
    set -u

    if [ $res != "0" ]; then
        echo "ERROR: run e2e test case failed"
    fi

    echo "::set-env name=CI_RES::$res"
}

TS=""
if [[ ${1} != "" ]];then
  TS=${1}
  runTest "$TS"
else
  ginkgo -r test
fi
