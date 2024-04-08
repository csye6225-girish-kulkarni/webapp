#!/bin/bash

# Fetch the image name
IMAGE_NAME=$(cat image_name.txt)
echo $IMAGE_NAME


gcloud compute instance-templates create "new-webapp-template" \
    --machine-type=$MACHINE_TYPE \
    --image=$NEW_IMAGE \
    --image-project=$NEW_IMAGE_PROJECT_ID \
    --boot-disk-size=${DISK_SIZE_GB}GB \
    --boot-disk-type=$DISK_TYPE \
    --network-interface=network=$NETWORK,subnet=$SUBNET \
    --instance-template-region=$REGION \
    --tags=$TAGS \
    --no-address \
    --metadata-from-file=startup-script=startup-script.sh \
    --service-account=$SERVICE_ACCOUNT_EMAIL \
    --scopes=$SCOPES

sleep 30

MANAGED_INSTANCE_GROUP_NAME="instance-group-manager"
gcloud compute instance-groups managed set-instance-template $MANAGED_INSTANCE_GROUP_NAME \
    --template=projects/cloud-csye6225-dev/regions/us-east1/instanceTemplates/$NEW_TEMPLATE_NAME \
    --region=us-east1

gcloud compute instance-groups managed rolling-action start-update $MANAGED_INSTANCE_GROUP_NAME \
    --version template=projects/cloud-csye6225-dev/regions/us-east1/instanceTemplates/$NEW_TEMPLATE_NAME \
    --type proactive \
    --max-unavailable 0 \
    --max-surge 3 \
    --region=us-east1

while [[ $(gcloud compute instance-groups managed list --filter="name=$MANAGED_INSTANCE_GROUP_NAME" --format="get(status)") != "STABLE" ]]; do
  echo "Waiting for instance group refresh to complete..."
  sleep 20
done

# Create a new instance template
#NEW_INSTANCE_TEMPLATE_NAME="instance-template-$IMAGE_NAME"
#gcloud compute instance-templates create $NEW_INSTANCE_TEMPLATE_NAME \
#    --image-project=$GCP_PROJECT_ID \
#    --image=$IMAGE_NAME \
#    --machine-type=MACHINE_TYPE
#
## Configure the managed instance group to use the new template
#gcloud compute instance-groups managed set-instance-template MANAGED_INSTANCE_GROUP_NAME \
#    --template=$NEW_INSTANCE_TEMPLATE_NAME
#
## Recreate instances in the managed instance group
#gcloud compute instance-groups managed rolling-action start-update MANAGED_INSTANCE_GROUP_NAME \
#    --version template=$NEW_INSTANCE_TEMPLATE_NAME \
#    --type proactive \
#    --max-unavailable 0 \
#    --max-surge 3
#
## Wait for the instance group refresh to complete
#while [[ $(gcloud compute instance-groups managed list --filter="name=MANAGED_INSTANCE_GROUP_NAME" --format="get(status)") != "STABLE" ]]; do
#  echo "Waiting for instance group refresh to complete..."
#  sleep 20
#done