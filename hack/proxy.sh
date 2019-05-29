#!/usr/bin/env bash

echo "WARNING: Make sure ports 5432 is open locally!"
echo "Starting proxy...(Only supporting us-east-1)"

# Get Jump
jump=$(aws ec2 describe-instances --region us-east-1 --filter "Name=tag-value,Values=jump*"|grep -i "publicipaddress"|cut -d: -f2|sed 's/[,"]//g'|tr -d " ")
rdsdv=$(aws ec2 describe-instances --region us-east-1 --filter "Name=tag-value,Values=rds-data-validator"|grep -i "privateipaddress"|head -1|cut -d: -f2|sed 's/[,"]//g'|tr -d " ")

if [ "${jump}" = "" ] || [ "${rdsdv}" = "" ]; then
  echo "Could not locate either the Jump Server or the RDS DV Server. Are you in a production vaulted shell?"
  exit 1
fi

ssh -L 5432:localhost:5432 -J $jump ec2-user@$rdsdv -N
