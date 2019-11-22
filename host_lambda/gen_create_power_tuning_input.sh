#!/usr/bin/env bash

cat create_message  $(jq '.|tojson' create_message.json | sed 's:^.\(.*\).$:\1:') > create_power_tuning_inpput.json