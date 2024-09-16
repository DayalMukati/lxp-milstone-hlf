#!/bin/bash

# Initialize score
score=0
source ./scripts/setOrgPeerContext.sh 1
export PATH=${PWD}/../bin:${PWD}:$PATH
export FABRIC_CFG_PATH=${PWD}/configtx

# Query the milestone payment
echo "Querying the milestone payment for order1..."
QUERY_OUTPUT=$(peer chaincode query -C mychannel -n milestonepayment -c '{"Args":["QueryMilestonePayments","order1"]}' 2>&1)

# Check if the milestone exists
if [[ $QUERY_OUTPUT == *"Error"* ]]; then
    echo "Milestone payment not found. It was not created."
    # Output the final score
    echo "Final Score: $score/50"
    exit 0
else
    echo "Milestone payment found."
fi

# Check if the milestone status is "created" (indicating creation)
if [[ $QUERY_OUTPUT == *"pending"* ]]; then
    echo "Milestone payment creation successful."
    score=$((score + 20))
fi

# Check if the milestone status is "completed" (indicating status update)
if [[ $QUERY_OUTPUT == *"completed"* ]]; then
    echo "Milestone status update successful."
    score=$((score + 30))
fi

# Check if the milestone status is "paid" (indicating payment release)
if [[ $QUERY_OUTPUT == *"paid"* ]]; then
    echo "Milestone payment release successful."
    score=$((score + 50))
fi

# Final score output
echo "Final Score: $score/50"

# Exit with success
exit 0
