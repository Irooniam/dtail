CREATE OR REPLACE FUNCTION dtail_%s_update() RETURNS TRIGGER AS $$
	DECLARE
		row json;
		op text;

    	BEGIN
		op = '{"op":"' || TG_OP || '", "table":"%s"}';
		row = op::jsonb || row_to_json(NEW)::jsonb;
    		PERFORM pg_notify('dtail_table_update', row::text);
		RETURN NEW;
    	END;

$$ LANGUAGE plpgsql;
