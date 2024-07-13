CREATE OR REPLACE FUNCTION lower_email_address()
RETURNS TRIGGER AS
$$
BEGIN
  NEW.email_address = LOWER(NEW.email_address);
  RETURN NEW;
END;
$$
LANGUAGE plpgsql;

CREATE TRIGGER trg_lower_email_address
BEFORE INSERT OR UPDATE ON email_addresses
FOR EACH ROW EXECUTE PROCEDURE lower_email_address();

CREATE OR REPLACE FUNCTION upsert_address_geom() RETURNS TRIGGER AS
$$
BEGIN
  -- we only want to do an st_point if it's not been set yet. the likely hood
  -- of this happening anyway (creating new area codes geolication) won't
  -- happen on a single basis but rather batch, I think so anyway.
  IF NEW.longitude IS NOT NULL AND NEW.latitude IS NOT NULL AND NEW.geom IS NULL THEN
    NEW.geom = ST_SetSRID(ST_MakePoint(NEW.longitude, NEW.latitude), 4326);
  END IF;
  RETURN NEW;
END;
$$
LANGUAGE plpgsql;

CREATE TRIGGER trg_upsert_address_geom
BEFORE INSERT ON addresses
FOR EACH ROW EXECUTE PROCEDURE upsert_address_geom();

CREATE TRIGGER trg_update_address_geom
BEFORE UPDATE ON addresses
FOR EACH ROW
WHEN (
    (old.latitude IS DISTINCT FROM new.latitude)
    AND (old.longitude IS DISTINCT FROM new.longitude)
)
EXECUTE PROCEDURE upsert_address_geom();

CREATE FUNCTION create_customer(
    v_last_name VARCHAR,
    v_first_name VARCHAR,
    v_contact_info JSONB,
    v_address JSONB
) RETURNS TABLE (
    customer_id INTEGER,
    last_name VARCHAR,
    first_name VARCHAR,
    contact_info JSONB,
    address JSONB
) AS
$$
declare
    _customer_id INTEGER;
    _address_id INTEGER;
BEGIN
  INSERT INTO customers AS c (
      first_name,
      last_name
  ) VALUES (
      v_first_name,
      v_last_name
  ) RETURNING c.customer_id INTO _customer_id;


  INSERT INTO email_addresses (
      customer_id,
      email_address
  ) VALUES (
      _customer_id,
      v_contact_info->>'email_address'
  );

  INSERT INTO addresses as a (
      street_number,
      route,
      unit_number,
      locality,
      administrative_area_level_1,
      country,
      postal_code,
      latitude,
      longitude,
      geodata
  ) VALUES (
      v_address->>'street_number',
      v_address->>'route',
      v_address->>'unit_number',
      v_address->>'locality',
      v_address->>'administrative_area_level_1',
      v_address->>'country',
      v_address->>'postal_code',
      CAST(v_address->>'latitude' AS NUMERIC),
      CAST(v_address->>'longitude' AS NUMERIC),
      CAST(v_address->>'geodata' AS JSONB)
  ) RETURNING a.address_id INTO _address_id;

  INSERT INTO customer_addresses (
    customer_id,
    address_id
  ) VALUES (
    _customer_id,
    _address_id
  );

  RETURN QUERY
  SELECT
      c.customer_id,
      c.last_name,
      c.first_name,
      ea.contact_info,
      a.address
  FROM customers c
  JOIN (
      SELECT
        _ea.customer_id,
        jsonb_build_object(
            'email_address_id', _ea.email_address_id,
            'email_address', _ea.email_address
        ) AS contact_info
      FROM email_addresses _ea
      WHERE _ea.customer_id = _customer_id
  ) ea USING (customer_id)
  JOIN (
      SELECT
          _ca.customer_id,
          jsonb_build_object(
			  'address_id', _a.address_id,
			  'street_number', _a.street_number,
			  'route', _a.route,
			  'unit_number', _a.unit_number,
			  'locality', _a.locality,
			  'administrative_area_level_1', _a.administrative_area_level_1,
			  'country', _a.country,
			  'postal_code', _a.postal_code,
			  'latitude', _a.latitude,
			  'longitude', _a.longitude,
			  'geodata', _a.geodata
          ) AS address
      FROM customer_addresses _ca
      JOIN addresses _a using (address_id)
      WHERE _ca.customer_id = _customer_id
      AND _ca.address_id = _address_id
  ) a USING (customer_id);
END;
$$
LANGUAGE plpgsql;
