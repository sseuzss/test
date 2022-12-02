#!/bin/bash

echo "Creating user for vxserver"
mysql --host=${DB_HOST} --user=${MYSQL_ROOT_USER} --password=${MYSQL_ROOT_PASSWORD} --port=${DB_PORT} --execute="CREATE DATABASE IF NOT EXISTS ${DB_NAME};"
mysql --host=${DB_HOST} --user=${MYSQL_ROOT_USER} --password=${MYSQL_ROOT_PASSWORD} --port=${DB_PORT} --execute="ALTER DATABASE ${DB_NAME} DEFAULT CHARACTER SET utf8 DEFAULT COLLATE utf8_unicode_ci;"
mysql --host=${DB_HOST} --user=${MYSQL_ROOT_USER} --password=${MYSQL_ROOT_PASSWORD} --port=${DB_PORT} --execute="CREATE DATABASE IF NOT EXISTS ${AGENT_SERVER_DB_NAME};"
mysql --host=${DB_HOST} --user=${MYSQL_ROOT_USER} --password=${MYSQL_ROOT_PASSWORD} --port=${DB_PORT} --execute="ALTER DATABASE ${AGENT_SERVER_DB_NAME} DEFAULT CHARACTER SET utf8 DEFAULT COLLATE utf8_unicode_ci;"
mysql --host=${DB_HOST} --user=${MYSQL_ROOT_USER} --password=${MYSQL_ROOT_PASSWORD} --port=${DB_PORT} --execute="CREATE USER IF NOT EXISTS '${AGENT_SERVER_DB_USER}' IDENTIFIED BY '${AGENT_SERVER_DB_PASS}';"
mysql --host=${DB_HOST} --user=${MYSQL_ROOT_USER} --password=${MYSQL_ROOT_PASSWORD} --port=${DB_PORT} --execute="GRANT ALL PRIVILEGES ON ${AGENT_SERVER_DB_NAME}.* TO ${AGENT_SERVER_DB_USER}@'%';"
mysql --host=${DB_HOST} --user=${MYSQL_ROOT_USER} --password=${MYSQL_ROOT_PASSWORD} --port=${DB_PORT} "$DB_NAME" < /opt/vxbinaries/seed.sql
echo "done"

mc config host add vxm "${MINIO_ENDPOINT}" ${MINIO_ACCESS_KEY} ${MINIO_SECRET_KEY} 2>/dev/null
mc mb --ignore-existing vxm/${MINIO_BUCKET_NAME}
mc cp --recursive /opt/vxdbmigrate/utils vxm/${MINIO_BUCKET_NAME}
mc config host add vxinst "${MINIO_ENDPOINT}" ${MINIO_ACCESS_KEY} ${MINIO_SECRET_KEY} 2>/dev/null
mc mb --ignore-existing vxinst/${AGENT_SERVER_MINIO_BUCKET_NAME}
mc cp --recursive /opt/vxdbmigrate/utils vxinst/${AGENT_SERVER_MINIO_BUCKET_NAME}

cat <<EOT > /opt/vxbinaries/upload_binaries.sql
INSERT
IGNORE INTO \`tenants\` (
    \`id\`,
    \`hash\`,
    \`status\`,
    \`description\`
 ) VALUES (1, MD5(RAND()), 'active', 'First user tenant');

INSERT
IGNORE INTO \`services\` (
    \`id\`,
    \`hash\`,
    \`tenant_id\`,
    \`name\`,
    \`type\`,
    \`status\`,
    \`info\`
) VALUES (
    1,
    MD5('localhost'),
    1,
    'localhost',
    'vxmonitor',
    'active',
    '{
      "db": {
        "host": "$DB_HOST",
        "name": "$AGENT_SERVER_DB_NAME",
        "pass": "$AGENT_SERVER_DB_PASS",
        "port": $DB_PORT,
        "user": "$AGENT_SERVER_DB_USER"
      },
      "s3": {
        "access_key": "$MINIO_ACCESS_KEY",
        "secret_key": "$MINIO_SECRET_KEY",
        "bucket_name": "$MINIO_BUCKET_NAME",
        "endpoint": "$MINIO_ENDPOINT"
      },
      "server": {
        "host": "vxserver.local",
        "port": 8443,
        "proto": "wss"
      }
    }'
);

INSERT
IGNORE INTO \`users\` (
    \`id\`,
    \`hash\`,
    \`status\`,
    \`role_id\`,
    \`tenant_id\`,
    \`mail\`,
    \`name\`,
    \`password\`
) VALUES (
    1,
    MD5(RAND()),
    'active',
    1,
    1,
    'admin@vxcontrol.com',
    'admin',
    '\$2a\$10\$deVOk0o1nYRHpaVXjIcyCuRmaHvtoMN/2RUT7w5XbZTeiWKEbXx9q'
);

EOT
