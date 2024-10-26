if [ -z $KEY ]; then
    echo "Key is unset, please specify KEY= when invoking this script."
    exit 1
fi

USER="ubuntu"
HOST="wubzduh.grattafiori.dev"
URL="$USER@$HOST"

GOARCH=amd64 GOOS=linux go build ../

mkdir build
cp -r ../static build/
cp -r ../templates build/
cp env.txt build/
cp wubz.service build/
mv web build/

tar -czvf build.tar.gz build
rm -rf build

scp -i $KEY build.tar.gz $URL:~/
ssh -i $KEY $URL "tar -xvzf build.tar.gz"
ssh -i $KEY $URL "rm build.tar.gz"
ssh -i $KEY $URL "sudo systemctl restart wubz"