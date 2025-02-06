CREATE OR REPLACE FUNCTION notify_my_table_update() RETURNS TRIGGER AS $$
	DECLARE
		row json;
		op text;

    	BEGIN
		op = '{"op":"' || TG_OP || '"}';
		row = op::jsonb || row_to_json(NEW)::jsonb;
    		PERFORM pg_notify('my_table_update', row::text);
		RETURN NEW;
    	END;

$$ LANGUAGE plpgsql;
