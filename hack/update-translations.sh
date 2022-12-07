#!/usr/bin/env bash

KUBE_ROOT=$(dirname "${BASH_SOURCE[0]}")/..
source "${KUBE_ROOT}/hack/lib/util.sh"

KUBECTL_FILES="pkg/kubectl/cmd/*.go pkg/kubectl/cmd/*/*.go"

generate_pot="false"
generate_mo="false"

while getopts "hf:xg" opt; do
  case ${opt} in
    h)
      echo "$0 [-f files] [-x] [-g]"
      echo " -f <file-path>: Files to process"
      echo " -x extract strings to a POT file"
      echo " -g sort .po files and generate .mo files"
      exit 0
      ;;
    f)
      KUBECTL_FILES="${OPTARG}"
      ;;
    x)
      generate_pot="true"
      ;;
    g)
      generate_mo="true"
      ;;
    \?)
      echo "[-f <files>] -x -g" >&2
      exit 1
      ;;
  esac
done

if ! which go-xgettext > /dev/null; then
  echo 'Can not find go-xgettext, install with:'
  echo 'go get github.com/gosexy/gettext/go-xgettext'
  exit 1
fi

if ! which msgfmt > /dev/null; then
  echo 'Can not find msgfmt, install with:'
  echo 'apt-get install gettext'
  exit 1
fi

if [[ "${generate_pot}" == "true" ]]; then
  echo "Extracting strings to POT"
  go-xgettext -k=i18n.T "${KUBECTL_FILES}" > tmp.pot
  perl -pi -e 's/CHARSET/UTF-8/' tmp.pot
  perl -pi -e 's/\\\(/\\\\\(/g' tmp.pot
  perl -pi -e 's/\\\)/\\\\\)/g' tmp.pot
  kube::util::ensure-temp-dir
  if msgcat -s tmp.pot > "${KUBE_TEMP}/template.pot"; then
    mv "${KUBE_TEMP}/template.pot" translations/kubectl/template.pot
    rm tmp.pot
  else
    echo "Failed to update template.pot"
    exit 1
  fi
fi

if [[ "${generate_mo}" == "true" ]]; then
  echo "Generating .po and .mo files"
  for x in translations/*/*/*/*.po; do
    msgcat -s "${x}" > tmp.po
    mv tmp.po "${x}"
    echo "generating .mo file for: ${x}"
    msgfmt "${x}" -o "$(dirname "${x}")/$(basename "${x}" .po).mo"
  done
fi

./hack/generate-bindata.sh
