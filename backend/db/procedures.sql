CREATE OR REPLACE FUNCTION batch_label_process(
	labels_json JSONB,
	user_uuid UUID
) RETURNS JSONB AS $$
DECLARE
	label_record JSONB;
	label_id_val VARCHAR(255);
	existing_label_id VARCHAR(255);
	new_labels JSONB := '[]'::JSONB;
	duplicate_labels JSONB := '[]'::JSONB;
	new_count INTEGER := 0;
	duplicate_count INTEGER := 0;
BEGIN
	FOR label_record IN SELECT * FROM jsonb_array_elements(labels_json)
	LOOP
		label_id_val := label_record->>'ID';

		SELECT label_id INTO existing_label_id 
		FROM labels 
		WHERE label_id = label_id_val;

		IF existing_label_id IS NULL THEN
			INSERT INTO labels (
				label_id, location, bundle_no, pqd, unit, time, length,
				heat_no, product_heading, isi_bottom, isi_top, charge_dtm,
				mill, grade, url_apikey, weight, section, date, user_id,
				status, is_duplicate
			) VALUES (
				label_id_val,
				label_record->>'LOCATION',
				(label_record->>'BUNDLE_NO')::INTEGER,
				label_record->>'PQD',
				label_record->>'UNIT',
				label_record->>'TIME',
				(label_record->>'LENGTH')::INTEGER,
				label_record->>'HEAT_NO',
				label_record->>'PRODUCT_HEADING',
				label_record->>'ISI_BOTTOM',
				label_record->>'ISI_TOP',
				label_record->>'CHARGE_DTM',
				label_record->>'MILL',
				label_record->>'GRADE',
				label_record->>'URL_APIKEY',
				label_record->>'WEIGHT',
				label_record->>'SECTION',
				label_record->>'DATE',
				user_uuid,
				'pending',
				false
			);

			new_labels := new_labels || label_record;
			new_count := new_count + 1;
		ELSE
			duplicate_labels := duplicate_labels || label_record;
			duplicate_count := duplicate_count + 1;
		END IF;
	END LOOP;

	RETURN jsonb_build_object(
		'new_labels', new_labels,
		'duplicate_labels', duplicate_labels,
		'total_processed', jsonb_array_length(labels_json),
		'new_count', new_count,
		'duplicate_count', duplicate_count
	);
END;
$$ LANGUAGE plpgsql;
