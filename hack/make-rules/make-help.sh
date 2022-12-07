#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

readonly red=$(tput setaf 1)
readonly reset=$(tput sgr0)

KUBE_ROOT=$(dirname "${BASH_SOURCE[0]}")/../..
ALL_TARGETS=$(make -C "${KUBE_ROOT}" PRINT_HELP=y -rpn | sed -n -e '/^$/ { n ; /^[^ .#][^ ]*:/ { s/:.*$// ; p ; } ; }' | sort)
CMD_TARGETS=$(cd "${KUBE_ROOT}/cmd"; find . -mindepth 1 -maxdepth 1 -type d | cut -c 3-)
CMD_FLAG=false

echo "--------------------------------------------------------------------------------"
for tar in ${ALL_TARGETS}; do
	for cmdtar in ${CMD_TARGETS}; do
		if [ "${tar}" = "${cmdtar}" ]; then
			if [ ${CMD_FLAG} = true ]; then
				continue 2;
			fi

			echo -e "${red}${CMD_TARGETS}${reset}"
			make -C "${KUBE_ROOT}" "${tar}" PRINT_HELP=y
			echo "---------------------------------------------------------------------------------"

			CMD_FLAG=true
			continue 2
		fi
	done

	echo -e "${red}${tar}${reset}"
	make -C "${KUBE_ROOT}" "${tar}" PRINT_HELP=y
	echo "---------------------------------------------------------------------------------"
done
