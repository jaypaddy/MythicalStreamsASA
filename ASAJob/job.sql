WITH   
STEP1 AS  
(  
  SELECT    Count(*) as membercount, 
            Id,
            Counter,
            DtTm,
            HHmm,
            Status,
            City,
            EventEnqueuedUtcTime
  FROM newmembers TIMESTAMP BY EventEnqueuedUtcTime 
  PARTITION BY City  
  GROUP BY  City,
            Id,
            Counter,
            DtTm,
            HHmm,
            Status,
            City,
            EventEnqueuedUtcTime,
            TumblingWindow(Duration(minute, 1), Offset(millisecond, -1))
)  

SELECT    CONCAT(   STEP1.City,
                    '-',
                DATEPART(yyyy,STEP1.EventEnqueuedUtcTime),
                '-',
                DATEPART(mm,STEP1.EventEnqueuedUtcTime),
                '-',
                DATEPART(dd,STEP1.EventEnqueuedUtcTime),
                ' ',
                DATEPART(hh,STEP1.EventEnqueuedUtcTime),
                ':',
                DATEPART(mi,STEP1.EventEnqueuedUtcTime),
                ':',
                DATEPART(ss,STEP1.EventEnqueuedUtcTime),
                ':'
                )   AS DISCRIMINATOR1,
              *
INTO newmembersagg FROM STEP1
/*INTO azfunc FROM STEP1*/