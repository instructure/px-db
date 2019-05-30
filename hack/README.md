# Practice DB Sanitization Process

It should be noted this process is somewhat manual and can be optimized later if deemed valuable
Unfortunately, this process sucks because of the restrictions around RDS Permissions - even after granting `rds_superuser` table constraints cannot be altered.
As such we must perform multiple `pg_dumps`

## Steps

Replace the RDS Endpoint below, with whatever it actually is - although it may be the same

Obtaining Jump Instance and Active RDS DV Instance:

```
jump=$(aws ec2 describe-instances --region us-east-1 --filter "Name=tag-value,Values=jump*"|grep -i "publicipaddress"|cut -d: -f2|sed 's/[,"]//g'|tr -d " ")
rdsdv=$(aws ec2 describe-instances --region us-east-1 --filter "Name=tag-value,Values=rds-data-validator"|grep -i "privateipaddress"|head -1|cut -d: -f2|sed 's/[,"]//g'|tr -d " ")
```
1. SSH to the RDS DV Instance: `ssh -J $jump ec2-user@$rdsdv`
1. On the RDS DV Instance Run: `db_sanitization_rds_restore.sh` --- output will be sent to slack, after the restore is complete proceed below.
1. On the RDS DV Instance Run: `docker run -d --name px-db -e POSTGRES_DB=practice_development -e POSTGRES_PASSWORD=practice -e POSTGRES_USER=practice -p 5432:5432 postgres:9.6.11`
1. On the RDS DV Instance Run: `docker exec -i px-db bash -c "psql -h iad-px-db-prod-sanitization-restore.ceta0kkbian5.us-east-1.rds.amazonaws.com -U master -d practice -c \"ALTER USER practice WITH PASSWORD 'practice'\""` (You'll need the current production master RDS password)
1. On the RDS DV Instance Run: `docker exec -i -e POSTGRES_PASSWORD=practice px-db bash -c "pg_dump -v -h iad-px-db-prod-sanitization-restore.ceta0kkbian5.us-east-1.rds.amazonaws.com -U practice -d practice|gzip > /var/tmp/dump.gz"`
1. On the RDS DV Instance Run: `docker exec -i px-db bash -c "cat /var/tmp/dump.gz |gunzip| psql -h localhost -U practice -d practice_development"
1. Locally Run: `proxy.sh` to proxy the RDS DV implementation of PostgreSQL
1. Locally Run: `run.sh` to sanitize the DB
1. On the RDS DV Instance Run: `docker exec -i -e POSTGRES_PASSWORD=practice px-db bash -c "pg_dump -v -h localhost -U practice -d practice_development|gzip > /var/tmp/sanitized-dump.gz"`
1. On the RDS DV Instance Run: `docker cp px-db:/var/tmp/sanitized-dump.gz .`
1. Locally use `scp` to copy the dump to your local machine: `scp -o "ProxyJump ${jump}" ec2-user@"${rdsdv}":sanitized-dump.gz .`
1. Purge the container and sanitized dump on the RDS DV Instance: `docker rm -f px-db && rm -f sanitized-dump.gz`
1. Cleanup the RDS Restored Instance: `iad-px-db-prod-sanitization-restore`
