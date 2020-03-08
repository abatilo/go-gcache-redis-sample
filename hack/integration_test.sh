#!/bin/bash

seq 1 1000 | xargs -n1 -P4 bash -c 'i=$0; url="http://localhost:8000/$i"; curl $url'
