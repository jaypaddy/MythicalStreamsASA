rgExists=$(az group show --name $resourceGroupName)
if [ -z "$rgExists" ]
then
    echo "creating resourcegroup $resourceGroupName"
    az group create --name $resourceGroupName --location $location
fi

echo "checking eventhub $ehNamespace"
ehExists=$(az eventhubs eventhub show \
    --resource-group $resourceGroupName \
    --namespace-name $ehNamespace)

#EH Namespace and Event Hubs
if [ -z "$ehExists" ]
then  
    echo "creating eventhub $ehNamespace"
    #Create Event Hub Namespace
    az eventhubs namespace create \
                --name $ehNamespace \
                --resource-group $resourceGroupName \
                -l $location
    #Create Event Hub
    # Create ingress event hub. Specify a name for the event hub. 
    echo "creating eventhub $ehNamespace - $ingestionEventhub"
    az eventhubs eventhub create \
                --name $ingestionEventhub \
                --resource-group $resourceGroupName \
                --namespace-name $ehNamespace \
                --partition-count $partitionCount

    # Create egress event hub
    echo "creating eventhub $ehNamespace - $egressEventhub"
    az eventhubs eventhub create \
                --name $egressEventhub \
                --resource-group $resourceGroupName \
                --namespace-name $ehNamespace \
                --partition-count $partitionCount
    
    # Create a consumer group for function trigger
    echo "creating consumergroup - egressEHConsumerGroup"
    az eventhubs eventhub consumer-group create \
                --resource-group $resourceGroupName \
                --namespace-name $ehNamespace \ 
                --eventhub-name $egressEventhub \
                --name egressEHConsumerGroup

    #get the SAS and load it for functions
    echo "creating authorization rule - $ehNamespace"
    authorizationRuleName = $(az eventhubs namespace authorization-rule list \
                --namespace-name $ehNamespace \ 
                --resource-group $resourceGroupName \
                --query [0].name -o tsv )
    
    primaryConnectionString = $(az eventhubs namespace authorization-rule keys list \
                --name $authorizationRuleName \
                --namespace-name $ehNamespace \ 
                --resource-group $resourceGroupName \
                --query primaryConnectionString -o tsv)
                                                    
    #create function plan 
    #az functionapp plan create  \
    #            --resource-group $resourceGroupName \
    #            --name $(functionAppPlan) 
    #            --sku 'consumption' 

    #create function app (CONSUMPTION PLAN)

    az functionapp create \
                --name $functionApp \
                --consumption-plan-location $location \
                --os-type Linux  \
                --resource-group $resourceGroupName  \
                --runtime java  \
                --storage-account $funcstore   \
    
fi


