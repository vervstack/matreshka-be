CREATE SCHEMA IF NOT EXISTS matreshka;
CREATE ROLE configurator;

DO
$$BEGIN
    IF NOT EXISTS (
        SELECT
        FROM   pg_catalog.pg_roles
        WHERE  rolname = 'makosh'
    ) THEN
CREATE USER makosh WITH PASSWORD 'da8331e37a';
GRANT configurator TO makosh;
END IF;
END$$;

ALTER ROLE configurator SET search_path TO matreshka;
ALTER ROLE configurator SET DEFAULT_TABLESPACE TO matreshka;

GRANT USAGE ON SCHEMA matreshka TO configurator;
GRANT SELECT, INSERT, UPDATE, DELETE ON ALL TABLES IN SCHEMA matreshka TO configurator;
GRANT USAGE, SELECT ON ALL SEQUENCES IN SCHEMA matreshka TO configurator;

ALTER DEFAULT PRIVILEGES IN SCHEMA matreshka
    GRANT SELECT, INSERT, UPDATE, DELETE ON TABLES TO
    configurator;
ALTER DEFAULT PRIVILEGES IN SCHEMA matreshka
    GRANT USAGE, SELECT ON SEQUENCES TO
    configurator;