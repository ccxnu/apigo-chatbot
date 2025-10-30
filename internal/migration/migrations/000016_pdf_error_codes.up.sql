-- Add error codes for PDF processing functionality

do $$
begin
    -- ERR_PDF_PROCESSING
    if not exists (select 1 from cht_parameters where prm_code = 'ERR_PDF_PROCESSING') then
        insert into cht_parameters (prm_name, prm_code, prm_data, prm_description)
        values ('ERROR_CODES', 'ERR_PDF_PROCESSING', '{"message": "Error al procesar el archivo PDF"}'::jsonb, 'Error processing PDF file');
    end if;

    -- ERR_CHUNK_CREATION
    if not exists (select 1 from cht_parameters where prm_code = 'ERR_CHUNK_CREATION') then
        insert into cht_parameters (prm_name, prm_code, prm_data, prm_description)
        values ('ERROR_CODES', 'ERR_CHUNK_CREATION', '{"message": "Error al crear los fragmentos del documento"}'::jsonb, 'Error creating document chunks');
    end if;
end $$;
