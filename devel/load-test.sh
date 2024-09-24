#!/bin/bash
# Our custom function
func_acc(){
  echo "Creating token number: $1"
  vault write ${ENGINE_PATH}/account/wfe-ops ttl=20s
}

func_proj(){
  echo "Creating token number: $1"
  vault write ${ENGINE_PATH}/project/wfecd-stg-unprotected/role/unprotected-role ttl=20s
}

for i in {1..40}
do
	func_proj $i &
	sleep 0.50 # comment out this line to load test with higher concurrent requests
done

wait
echo "All done"