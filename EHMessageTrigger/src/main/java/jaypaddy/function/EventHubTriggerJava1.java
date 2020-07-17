package jaypaddy.function;

import com.microsoft.azure.functions.annotation.*;
import com.microsoft.azure.functions.*;
import java.util.*;

/**
 * Azure Functions with Event Hub trigger.
 */
public class EventHubTriggerJava1 {
    /**
     * This function will be invoked when an event is received from Event Hub.
     */
    @FunctionName("EventHubTriggerJava1")
    public void run(
        @EventHubTrigger(name = "message", 
        eventHubName = "merged", 
        connection = "functrigger_merged_policy", 
        consumerGroup = "functrigger", cardinality = Cardinality.ONE) 
        String message,
        @EventHubOutput(name = "final", 
        eventHubName = "final", 
        connection = "funcwriter_merged_policy") OutputBinding<Object> toEH,
        final ExecutionContext context
    ) {
        context.getLogger().info(message.toString());
        toEH.setValue(message);    }
}
