#!/bin/bash
# Our custom function
func(){
  echo "Crating token number: $1"
  vault write wfecd-stg/account/repo-reporting ttl=20s
}
for i in {1..500}
do
	func $i &
	sleep 0.50
done

wait
echo "All done"