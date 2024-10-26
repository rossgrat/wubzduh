
if [ -z $KEY ]; then
    echo "Key is unset, please specify KEY= when invoking this script."
    exit 1
fi

USER="ubuntu"
HOST="wubzduh.grattafiori.dev"
URL="$USER@$HOST"

# Install postgres and initialize the database
scp -i $KEY ../src/db/init-db.sql $URL:~/
scp -i $KEY ../../../src/config/setup.sh $URL:~/
ssh -i $KEY $URL "./setup.sh"
ssh -i $KEY $URL "rm init-db.sql"
ssh -i $KEY $URL "rm setup.sh"

# Copy the service over to the system
KEY=$KEY ./deploy-daemon.sh
# Bulld the web cmd, copy cmd, templates, and environment to server
KEY=$KEY ./build-and-deploy.sh
