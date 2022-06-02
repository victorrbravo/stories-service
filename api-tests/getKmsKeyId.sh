#!/bin/bash

function get_kms_key_lbd_config_id() {
   local query='Aliases[?AliasName==`alias/aclaraciones-dev`].TargetKeyId'
   local kms_aliases=$(aws kms list-aliases --query $query)
   local actual_count=$(echo $kms_aliases | jq 'length')
   local expected_count=1

   if [ $actual_count -ne $expected_count ]; then
      echo "Expected KMS aliases: $expected_count. Retrieved: $actual_count." >>/dev/stderr
      return $GET_KMS_ALIASES_ENDPOINT_FAIL
   fi

   echo $kms_aliases | jq --raw-output '.[0]'
}

kms_lbd_config_id=$(get_kms_key_lbd_config_id) || exit $?

echo $kms_lbd_config_id
