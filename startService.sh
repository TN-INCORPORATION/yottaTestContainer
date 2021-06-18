#!/bin/bash

cd /application

. /ydbdir/ydbenv
if [ $? -ne 0 ]
then
  echo "Load YDBENV_FILE Fail!!!"
  exit 1
fi
exec ./app
exit $?
