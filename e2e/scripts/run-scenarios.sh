date=$(date +%s)
results_file="results/$date.txt"
total_errors=0

_jq() {
 echo ${row} | base64 --decode
}

_get_assertion() {
 echo $(cat $TEST_FILE | jq -r ".assert.${1}")
}

_get_scenario() {
 echo $(cat $TEST_FILE | jq -r ".scenario")
}

_get_command() {
 echo $(cat $TEST_FILE | jq -r ".command")
}

run_scenario() {
    for TEST_FILE in ./*.json; do
        echo "" | tee -a $results_file
        echo "Running scenario: $(_get_scenario $(_jq))" | tee -a $results_file
        echo "Running command $(_get_command $(_jq))"
        response=$(eval $(_get_command $(_jq)))
        error_count=0

        for row in $(cat $TEST_FILE | jq -r '.assert | [paths | join(".")] | . - ["data"]' | jq -r '.[] | @base64'); do

            assertion=$(_get_assertion $(_jq) )

            response_value=$(echo $response | jq -r ".$(_jq)")

           if [ "$assertion" = "$response_value" ]; then
             echo "✅ Assertion $assertion equals response" >> $results_file
           else
             echo "❌ Assertion $assertion does not equal response" | tee -a $results_file

             echo "Expected object key: $(_jq) ...to have value: $assertion ...but was $response_value" | tee -a $results_file

             echo "Full response was $response" >> $results_file

             error_count=$((error_count+1))
           fi
        done

        if (( error_count > 0 )); then
            echo "❌ Scenario $(_get_scenario $(_jq)) failed with $error_count errors" | tee -a $results_file
        else
            echo "✅ Scenario $(_get_scenario $(_jq)) passed" | tee -a $results_file
        fi
        total_errors=$((error_count+total_errors))
    done
}

echo "Running tests from /e2e/scenarios"

pushd ./e2e/scenarios

mkdir -p results && touch $results_file

run_scenario

popd


if (( total_errors > 0 )); then
    echo "❌ Tests Failed with $total_errors errors"
    exit 1
else
    echo "✅ Tests passed"
fi