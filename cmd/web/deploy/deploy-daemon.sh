if [ -z $KEY ]; then
    echo "Key is unset, please specify KEY= when invoking this script."
    exit 1
fi

USER="ubuntu"
HOST="wubzduh.grattafiori.dev"
URL="$USER@$HOST"

scp -i $KEY ./wubz.service $URL:~/
ssh -i $KEY $URL "sudo rm /etc/systemd/system/wubz.service"
ssh -i $KEY $URL "sudo mv ~/wubz.service /etc/systemd/system/"
ssh -i $KEY $URL "sudo systemctl daemon-reload"
ssh -i $KEY $URL "sudo systemctl enable wubz.service"