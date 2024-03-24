CREATE TABLE IF NOT EXISTS clicks
(
    ID Int32,
    Projectid Int32,
    Name String,
    Description String,
    Priority Int32,
    Removed UInt8,
    EventTime DateTime
)
    ENGINE = MergeTree()
    ORDER BY ID;

ALTER TABLE clicks DROP INDEX IF EXISTS id_index;
ALTER TABLE clicks DROP INDEX IF EXISTS projectid_index;
ALTER TABLE clicks DROP INDEX IF EXISTS name_index;


ALTER TABLE clicks ADD INDEX id_index(ID) TYPE minmax GRANULARITY 1;
ALTER TABLE clicks ADD INDEX projectid_index(Projectid) TYPE minmax GRANULARITY 1;
ALTER TABLE clicks ADD INDEX name_index(Name) TYPE set(0) GRANULARITY 1;