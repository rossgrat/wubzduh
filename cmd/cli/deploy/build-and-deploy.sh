if [ -z $KEY ]; then
    echo "Key is unset, please specify KEY= when invoking this script."
    exit 1
fi

USER="ubuntu"
HOST="wubzduh.grattafiori.dev"
URL="$USER@$HOST"

GOARCH=amd64 GOOS=linux go build ../

mkdir build-cli
mv cli build-cli/
cp ./run-cli.sh build-cli/
cp ../../../src/config/env.txt build-cli/

tar -czvf build-cli.tar.gz build-cli
rm -rf build-cli

scp -i $KEY build-cli.tar.gz $URL:~/
ssh -i $KEY $URL "tar -xvzf build-cli.tar.gz"
ssh -i $KEY $URL "rm build-cli.tar.gz"
