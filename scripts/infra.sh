#!/bin/sh

cd "/service/terraform"

echo $1

if [ $1 = "create" ]; then
    echo "creating infrastructure ..."
    export TF_LOG_PATH="error.log"
    export TF_LOG=TRACE
    terraform init
    terraform plan -out=tfplan
    terraform apply tfplan
    # cat error.log | grep "[ERROR"
else
    echo "deleting ..."
    terraform apply -destroy -auto-approve
fi


