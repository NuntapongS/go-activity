#!/bin/sh

## bcrypt passwd generator ##
#############################
MONGO_USER="root"
MONGO_PASSWORD="password"
MONGO_HOST="127.0.0.1"
MONGO_DB="goecho"
MONGO_COLLECTION="staffs"

CMD=$(which htpasswd 2>/dev/null)
OPTS="-nB -C 4"
USERNAME=$1
ROLE=$2

usage() {
        local script=$(basename $0)
        cat <<EOF
$script: Generate Bcrypt Hashed Passwords using htpasswd

Usage: $script username role
       - role is ADMIN or STAFF
EOF
        exit 1
}

check_config() {
    if [ -z $CMD ]; then
        printf "Exiting: htpasswd is missing.\n"
        exit 1
    fi

    if [ -z "$USERNAME" ]; then
            usage
    fi

    if [ -z "$ROLE" ]; then
            usage
    fi

    if [[ "$ROLE" != "ADMIN" && "$ROLE" != "STAFF" ]]; then
            usage
    fi
}

check_config $USERNAME
printf "Generating Bcrypt hash for username: $USERNAME\n\n"
OUTPUT=$($CMD $OPTS "")

HASH_PASSWORD=$(echo $OUTPUT| tr -d ':\n' | sed 's/$2y/$2a/')

mongosh -u $MONGO_USER -p $MONGO_PASSWORD $MONGO_HOST/$MONGO_DB <<EOF
db.$MONGO_COLLECTION.insertOne({
    username: "$USERNAME",
    hashed_password: "$HASH_PASSWORD",
    role: "$ROLE",
})
EOF
exit $?