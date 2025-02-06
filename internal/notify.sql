CREATE TRIGGER trigger_my_table_update
  AFTER INSERT OR UPDATE OR DELETE
  ON my_table
  FOR EACH ROW
  EXECUTE PROCEDURE notify_my_table_update();
