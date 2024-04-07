#!/bin/bash

# Fetch the image name
IMAGE_NAME=$(cat image_name.txt)

# Set the prefix of the existing instance template and the new image ID
EXISTING_TEMPLATE_PREFIX="webapp-instance"
NEW_IMAGE_ID=$(cat image_name.txt)
cat image_name.txt

# Fetch the name of the existing instance template
EXISTING_TEMPLATE_NAME=$(gcloud compute instance-templates list --format="get(name)" --filter="name~'webapp-instance'")

# Fetch the configuration of the existing instance template
#TEMPLATE_CONFIG=$(gcloud compute instance-templates describe $EXISTING_TEMPLATE_NAME --region us-east1 --format=json)

# Modify the image ID in the configuration
#MODIFIED_TEMPLATE_CONFIG=$(echo "$TEMPLATE_CONFIG" | jq --arg NEW_IMAGE_ID $NEW_IMAGE_ID '.properties.disks[0].initializeParams.sourceImage = $NEW_IMAGE_ID')

# Set the name of the new instance template
NEW_TEMPLATE_NAME="new-webapp-template"

# Write the modified configuration to a temporary file
#echo $MODIFIED_TEMPLATE_CONFIG > temp.json

# Create a new instance template that is identical to the existing one
# Create a new instance template with the new image
gcloud compute instance-templates create $NEW_TEMPLATE_NAME \
    --source-instance-template $EXISTING_TEMPLATE_NAME \
    --image-project cloud-csye6225-dev \
    --image $NEW_IMAGE_ID \
    --boot-disk-auto-delete

# Remove the temporary file
#rm temp.json


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