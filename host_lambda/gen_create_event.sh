#!/usr/bin/env bash

sam local generate-event sns notification --message $(jq '.|tojson' create_message.json | sed 's:^.\(.*\).$:\1:') > create_event.json