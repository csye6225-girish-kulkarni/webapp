#!/bin/bash

# Fetch the image name
IMAGE_NAME=$(cat image_name.txt)
echo $IMAGE_NAME

# Set the prefix of the existing instance template and the new image ID
EXISTING_TEMPLATE_PREFIX="webapp-instance"
cat image_name.txt

gcloud compute instances list --format="json" > temp.json

# Parse the JSON file to get the instance name and zone of the first instance
INSTANCE_NAME=$(jq -r '.[0].name' temp.json)
INSTANCE_ZONE_URL=$(jq -r '.[0].zone' temp.json)
INSTANCE_ZONE=$(basename $INSTANCE_ZONE_URL)

echo $INSTANCE_NAME
echo $INSTANCE_ZONE
# Extract the region from the zone
INSTANCE_REGION="${INSTANCE_ZONE%-*}"
echo $INSTANCE_REGION
INSTANTIATE_OPTIONS="custom-image,custom-image=projects/cloud-csye6225-dev/global/images/$IMAGE_NAME"
AUTO_DELETE="yes"
SOURCE_DISK="persistent-disk-0"

NEW_TEMPLATE_NAME="new-webapp-template"

# Create a new instance template
gcloud compute instance-templates create $NEW_TEMPLATE_NAME \
    --source-instance=$INSTANCE_NAME \
    --source-instance-zone=us-east1-b \
    --instance-template-region=us-east1 \
    --configure-disk=device-name=$SOURCE_DISK,instantiate-from=$INSTANTIATE_OPTIONS,auto-delete=$AUTO_DELETE

# Remove the temporary file
rm temp.json


MANAGED_INSTANCE_GROUP_NAME="instance-group-manager"
gcloud compute instance-groups managed set-instance-template $MANAGED_INSTANCE_GROUP_NAME \
    --template=$NEW_TEMPLATE_NAME

gcloud compute instance-groups managed rolling-action start-update $MANAGED_INSTANCE_GROUP_NAME \
    --version template=$NEW_TEMPLATE_NAME \
    --type proactive \
    --max-unavailable 0 \
    --max-surge 3

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