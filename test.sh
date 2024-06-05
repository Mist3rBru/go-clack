until [ ${RET} -eq 0 ]; do
  make test
  RET=$?
done