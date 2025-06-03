* what is kafka: 
  kafka is message broker supporting high throughput and parallism 

* what is the difference in kafka then rabbit
   true parallism in packet processing using partioning
   the order of of processing packet can be maintained within a partition 
   capable of handling more then 3million request/second 

* what are the most important things about kafka 
   kafka topic
   partitions of a topic : partioning is responsible to realize true parallesim for producer as as well as consumer
     producer can write parallely into different partitions. The consumers can truly consume parallely from different partitions.

   consumer group with a group id : to parally consume from different partitions we can have multiple consumer      instances with same group id; each instance will be associated with a partition to consume packet and upon restart that consumer having same group id will start consuming from where it left. But if the group id is different then    the consumer will start consuming from beginning.
        schenario : lets say we want all the packets of a particular topic to be processed again then we can make another consumer with a different group id and start it.

   offset : kakfa maintains offset/tracking numbner for each payload in a partition and tracks it 
    enable or disable auto commit offset ; if disabled then for each packet the consumer needs to commit that offset manually ; where this will be used : disables in banking system because each packet is most important; enabled in real time streaming video servies where we can manage and tolerate one missing packet

what happens in case of failue or we are not commiting the offset: 
  we can push the payload into another topic and process it with other consumer.

   
some patterns for failure cases :

a) retry then go to DLQ  and from DLQ (a separate consumer can process it later)