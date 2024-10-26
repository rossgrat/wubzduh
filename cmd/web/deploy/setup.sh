yes | sudo apt install postgresql
sudo -u postgres psql -c 'CREATE DATABASE wubz;'

# This conditional block should allow us
# to run the setup script locally or remotely
if ![ -e "init-db.sql" ]; then
    cp ../src/db/init-db.sql ./
fi
sudo -u postgres psql wubz < init-db.sql

if [ -e "init-db.sql" ]; then
    rm init-db.sql
fi