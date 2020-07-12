rgExists=$(az group show --name $(resourceGroupName))
if [ -z "$rgExists" ]
then
echo "creating resourcegroup $(resourceGroupName)"
az group create --name $(resourceGroupName) --location $(location)
fi

ehExists=$(az eventhubs eventhub show 
    --resource-group $(resourceGroupName) \
    --namespace-name $(ehNamespace))

#EH Namespace and Event Hubs
if [ -z "$ehExists" ]
then  
    #Create Event Hub Namespace
    az eventhubs namespace create \
                --name $(ehNamespace) \
                --resource-group $(resourceGroupName) \
                -l $location
    #Create Event Hub
    # Create an event hub. Specify a name for the event hub. 
    az eventhubs eventhub create \
                --name $(ingestionEventhub) \
                --resource-group $(resourceGroupName) \
                --namespace-name $(ehNamespace) \
                --partition-count $(partitionCount)

    az eventhubs eventhub create \
                --name $(egressEventhub) \
                --resource-group $(resourceGroupName) \
                --namespace-name $(ehNamespace) \
                --partition-count $(partitionCount)
else
    # check if the Event Hubs exist in the Namespace
fi


