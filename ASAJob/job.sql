WITH   
STEP1 AS  
(  
  SELECT
        NAXMLPOSJournal.TransmissionHeader.StoreLocationID AS STOREID,
        [NAXMLPOSJournal].[TransmissionHeader] AS TRANSMISSIONHEADER,
        [NAXMLPOSJournal].[JournalReport] AS JOURNALREPORT,
        [NAXMLPOSJournal].[PricingUID] AS PRICINGUID,
        COUNT(*),
        EventEnqueuedUtcTime
  FROM pos TIMESTAMP BY EventEnqueuedUtcTime 
  PARTITION BY STOREID  
  GROUP BY  STOREID,
            [NAXMLPOSJournal].TransmissionHeader.StoreLocationID, 
            [NAXMLPOSJournal].[TransmissionHeader],
            [NAXMLPOSJournal].[JournalReport],
            [NAXMLPOSJournal].[PricingUID],
            EventEnqueuedUtcTime,
            TumblingWindow(Duration(minute, 3), Offset(millisecond, -1))
)  

SELECT  
        CONCAT(
            STEP1.TRANSMISSIONHEADER.StoreLocationID,
            'D',
            DATEPART(yyyy,STEP1.EventEnqueuedUtcTime),
            DATEPART(mm,STEP1.EventEnqueuedUtcTime),
            DATEPART(dd,STEP1.EventEnqueuedUtcTime),
            'H',
            DATEPART(hh,STEP1.EventEnqueuedUtcTime),
            'M',
            udf.getbucket(DATEPART(mi,STEP1.EventEnqueuedUtcTime)))           
              AS DISCRIMINATOR1,
            STEP1.EventEnqueuedUtcTime,
              *
INTO posstore FROM STEP1
/*INTO azfunc FROM STEP1*/