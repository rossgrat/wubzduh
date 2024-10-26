
if [ -z $KEY ]; then
    echo "Key is unset, please specify KEY= when invoking this script."
    exit 1
fi

USER="ubuntu"
HOST="wubzduh.grattafiori.dev"
URL="$USER@$HOST"

# Install postgres and initialize the database
ssh -i $KEY $URL 'yes | sudo apt install postgresql'
ssh -i $KEY $URL "sudo -u postgres psql -c 'CREATE DATABASE wubz;'"
scp -i $KEY ../src/db/init-db.sql $URL:~/
ssh -i $KEY $URL "sudo -u postgres psql wubz < init-db.sql"
ssh -i $KEY $URL "rm init-db.sql"

# Copy the service over to the system
scp -i $KEY ./wubz.service $URL:~/
ssh -i $KEY $URL "sudo mv ~/wubz.service /etc/systemd/system/"
ssh -i $KEY $URL "sudo systemctl daemon-reload"
ssh -i $KEY $URL "sudo systemctl enable wubz.service"


# Bulld the web cmd, copy cmd, templates, and environment to server
KEY=$KEY ./build-and-deploy.sh
