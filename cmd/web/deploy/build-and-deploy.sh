if [ -z $KEY ]; then
    echo "Key is unset, please specify KEY= when invoking this script."
    exit 1
fi

USER="ubuntu"
HOST="wubzduh.grattafiori.dev"
URL="$USER@$HOST"

GOARCH=amd64 GOOS=linux go build ../

mkdir build-web
cp -r ../static build-web/
cp -r ../templates build-web/
cp ../../../src/config/env.txt build-web/
cp wubz.service build-web/
mv web build-web/

tar -czvf build-web.tar.gz build-web
rm -rf build-web

scp -i $KEY build-web.tar.gz $URL:~/
ssh -i $KEY $URL "tar -xvzf build-web.tar.gz"
ssh -i $KEY $URL "rm build-web.tar.gz"
ssh -i $KEY $URL "sudo systemctl restart wubz"