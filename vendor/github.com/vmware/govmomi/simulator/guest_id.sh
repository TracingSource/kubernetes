#!/bin/bash -e

pushd "$(dirname "$0")" >/dev/null

{
  cat <<EOF
package simulator

// GuestID is the list of valid types.VirtualMachineGuestOsIdentifier
var GuestID = []types.VirtualMachineGuestOsIdentifier{
EOF

  ids=($(grep 'VirtualMachineGuestOsIdentifier(' ../vim25/types/enum.go | grep = | awk '{print $1}'))
  printf "types.%s,\n" "${ids[@]}"

  echo "}"
} > guest_id.go

goimports -w guest_id.go
