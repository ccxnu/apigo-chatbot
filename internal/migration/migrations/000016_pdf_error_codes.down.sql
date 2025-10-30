-- Remove PDF processing error codes

delete from cht_parameters where prm_code in ('ERR_PDF_PROCESSING', 'ERR_CHUNK_CREATION');
