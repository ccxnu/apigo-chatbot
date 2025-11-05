-- Rollback: restore original function without category parameter

DROP FUNCTION IF EXISTS fn_similarity_search_chunks_hybrid(vector, text, int, float, float, varchar);

CREATE OR REPLACE FUNCTION fn_similarity_search_chunks_hybrid(
    p_query_embedding vector(1536),
    p_query_text text,
    p_limit int default 5,
    p_min_similarity float default 0.2,
    p_keyword_weight float default 0.15
)
RETURNS TABLE (
    chk_id int,
    chk_fk_document int,
    chk_content text,
    similarity_score float,
    keyword_score float,
    combined_score float,
    doc_title varchar,
    doc_category varchar
) AS $$
DECLARE
    v_tsquery tsquery;
BEGIN
    v_tsquery := plainto_tsquery('spanish', p_query_text);

    RETURN QUERY
    WITH ranked_chunks AS (
        SELECT
            c.chk_id,
            c.chk_fk_document,
            c.chk_content,
            (1 - (c.chk_embedding <=> p_query_embedding)) as semantic_score,
            ts_rank(c.chk_fts_vector, v_tsquery)::double precision as keyword_rank,
            d.doc_title,
            d.doc_category
        FROM public.cht_chunks c
        INNER JOIN public.cht_documents d ON c.chk_fk_document = d.doc_id
        WHERE d.doc_active = true
          AND c.chk_embedding IS NOT NULL
          AND ((1 - (c.chk_embedding <=> p_query_embedding)) >= p_min_similarity
               OR c.chk_fts_vector @@ v_tsquery)
    )
    SELECT
        rc.chk_id,
        rc.chk_fk_document,
        rc.chk_content,
        rc.semantic_score,
        rc.keyword_rank,
        (rc.semantic_score * (1 - p_keyword_weight)) + (rc.keyword_rank * p_keyword_weight) as combined,
        rc.doc_title,
        rc.doc_category
    FROM ranked_chunks rc
    ORDER BY combined DESC
    LIMIT p_limit;
END;
$$ LANGUAGE plpgsql STABLE;
