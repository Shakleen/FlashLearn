## Data Storage Solution Considerations and Direction

### User's Needs and Considerations

* **Initial Data Type:** Primarily structured data easily mappable to fields.
* **Future Data Type:** Need to store larger, less structured items like long text, images, and audio.
* **Schema Importance:** Maintaining a schema is important for data pipelines.
* **Data Relationships:** Mostly one-to-one and one-to-many, with few many-to-many.
* **Query Patterns:** Primarily querying based on time, but also using string matching, IDs, and relationships.
* **Data Volume:** Initial tens of thousands, growing quickly to millions of records.
* **Performance (Latency):** Minimal delay is required for both reads and writes.
* **Large Item Relationships:** One-to-many, with larger items shared between rows.
* **Consistency (Large Items):** Immediate updates are preferred for larger items in relation to structured data.
* **Team Familiarity:** More comfortable with SQL.
* **Transactional Needs:** No immediate complex transactional requirements.
* **Willingness for Multiple Databases:** Yes, open to using more than one type of database.
* **Top Priorities:** Cost and latency.
* **Cloud Platform:** Planning to use AWS for cloud compute and storage.
* **Large File Access:** Currently unsure about the exact access patterns for larger files, potentially requiring some processing.

### Chosen Direction

Given the above considerations, the recommended direction is a **hybrid approach** utilizing two main AWS services:

* **Amazon Relational Database Service (RDS)** for the core structured data.
* **Amazon Simple Storage Service (S3)** for the larger, less structured items (BLOBs).

### Reasons for Choosing This Direction

* **Structured Data and Schema Enforcement:** RDS (specifically considering PostgreSQL, MySQL, or MariaDB) is well-suited for managing structured data with defined schemas, which is crucial for your data pipelines. SQL databases excel at maintaining data integrity and consistency for relational data.
* **Querying Capabilities:** RDS offers robust querying capabilities using SQL, which you are already familiar with. It can efficiently handle time-based queries (with proper indexing), string matching, and queries based on IDs and relationships.
* **Scalability and Performance for Structured Data:** RDS allows you to choose instance types optimized for performance and scale as your data volume grows to millions of records. Read replicas can further improve read latency.
* **Cost-Effective Storage for Large Items:** S3 is a highly scalable and cost-effective solution for storing large amounts of unstructured data like images and audio. Its tiered storage options can help optimize costs based on access frequency.
* **Handling Large, Shared Items:** S3 is designed for storing and serving such content efficiently. You would store metadata (like S3 object URLs) in your RDS database to link the larger items to your structured data.
* **Preference for Immediate Consistency:** While using separate services requires application-level management for immediate consistency between RDS and S3, it provides more flexibility and potentially better performance and cost-effectiveness for storing and serving large binary objects compared to storing them directly in the SQL database.
* **AWS Integration:** Both RDS and S3 are core AWS services and integrate seamlessly with other AWS offerings, which can be beneficial for your overall architecture.
* **Alignment with Priorities:** This hybrid approach aims to balance cost-effectiveness (by using S3 for bulk storage) and low latency (by optimizing RDS and potentially using S3's CDN capabilities if direct serving is needed).

### Diagrams and Details

1. [Postgresql ERD](diagrams/erd_rdbms.png)

### Next Steps

Further refinement would involve determining the specific RDS engine that best fits your needs and cost considerations, and deciding on the exact access patterns and potential processing requirements for the larger files in S3. This will influence how you manage the links between the two services and ensure the desired level of consistency.