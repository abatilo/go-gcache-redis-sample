#!/bin/bash

function install_tools {
  grep _ tools/tools.go | awk -F'"' '{print $2}' | xargs -tI % go install %
}
install_tools
