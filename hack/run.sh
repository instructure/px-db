#!/usr/bin/env bash


function check_rc() {
  rc=$1

  if [ "$rc" != "0" ]; then
    exit $rc
  fi
}

DB_PASSWORD=$(printenv DB_PASSWORD)
if [ "$DB_PASSWORD" != "" ]; then
  export DB_PASSWORD=$DB_PASSWORD
fi

px-db --db-endpoint localhost --db-name practice_development --db-user practice sanitize delete-tables \
  --db-tables AccessCode,BulkGroupInvite,GroupInvitation,OrganizationInvitation,SamlIdentityProvider,SamlUserIdentity,Webhook
check_rc $?

px-db --db-endpoint localhost --db-name practice_development --db-user practice sanitize delete-tables \
  --cascade-mode \
  --db-tables LtiUser,LtiToolConsumer,LtiContextGroupMapping,LtiOutcomeServiceMapping,LtiResourceLink
check_rc $?

px-db --db-endpoint localhost --db-name practice_development --db-user practice plugin practice-pii
check_rc $?
