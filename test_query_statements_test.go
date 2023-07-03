package sqlstmt

var insertTestingDataTypeQuery = []byte(`
	INSERT INTO testing_datatypes (
		testing_datatype_uuid,
		word,
		paragraph,
		metadata,
		created_at
	) VALUES (
		:uuid,
		:word,
		:paragraph,
		:metadata,
		:created_at
	)
`)

var deleteTestingDataTypeQuery = []byte(`
	DELETE FROM testing_datatypes td
	WHERE (nullif(:uuid, NULL) IS NULL OR td.testing_datatype_uuid = :uuid)
	AND (nullif(:uuids, '{}') IS NULL OR td.testing_datatype_uuid = ANY(:uuids))
`)

var selectTestingDataTypeQuery = []byte(`
	SELECT
		td.testing_datatype_id,
		td.testing_datatype_uuid,
		td.word,
		td.paragraph,
		td.metadata,
		td.created_at
	FROM testing_datatypes td
	WHERE (nullif(:id, NULL) IS NULL OR td.testing_datatype_id = :id)
	AND (nullif(:ids, '{}') IS NULL OR td.testing_datatype_id = ANY(:ids))
	AND (nullif(:uuid, NULL) IS NULL OR td.testing_datatype_uuid = :uuid)
	AND (nullif(:uuids, '{}') IS NULL OR td.testing_datatype_uuid = ANY(:uuids))
	ORDER BY td.created_at
`)

var createCustomerQuery = []byte(`
	SELECT
		c.customer_id,
		c.last_name,
		c.first_name,
		c.contact_info,
		c.address
	FROM create_customer(
		:last_name,
		:first_name,
		:contact_info,
		:address
	) c
`)

var deleteCustomerQuery = []byte(`
	DELETE FROM customers c
	WHERE c.customer_id = (
		SELECT _ea.customer_id
	    FROM email_addresses _ea
	    WHERE _ea.customer_id = :customer_id
	    AND _ea.email_address = :email_address
	)
	RETURNING customer_id;
`)

var selectCustomerQuery = []byte(`
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
	  WHERE (NULLIF(:customer_id, NULL) IS NULL OR _ea.customer_id = :customer_id)
	  AND (NULLIF(:email_address, NULL) IS NULL OR _ea.email_address = :email_address)
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
	  WHERE (NULLIF(:customer_id, NULL) IS NULL OR _ca.customer_id = :customer_id)
	  AND (NULLIF(:address_id, NULL) IS NULL OR _ca.address_id = :address_id)
  ) a USING (customer_id);
`)
