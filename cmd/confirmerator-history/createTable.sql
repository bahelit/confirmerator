-- Create datebase then connect to database to create table.
-- create database confirmerator;

CREATE TABLE ethtxns
(
    "time" integer,
    txnfrom bytea COLLATE pg_catalog."default",
    txnto bytea COLLATE pg_catalog."default",
    gas bigint,
    gasprice bigint,
    block integer,
    txnhash text COLLATE pg_catalog."default",
    value numeric,
    contract_to bytea COLLATE pg_catalog."default",
    contract_value bytea COLLATE pg_catalog."default"
)

TABLESPACE pg_default;

CREATE INDEX block_index
    ON ethtxns USING btree
    (block)
    TABLESPACE pg_default;

CREATE INDEX contract_to_index
    ON ethtxns USING btree
    (contract_to COLLATE pg_catalog."default")
    TABLESPACE pg_default;

CREATE INDEX txnfrom_index
    ON ethtxns USING btree
    (txnfrom COLLATE pg_catalog."default")
    TABLESPACE pg_default;

CREATE INDEX txnto_index
    ON ethtxns USING btree
    (txnto COLLATE pg_catalog."default")
    TABLESPACE pg_default;