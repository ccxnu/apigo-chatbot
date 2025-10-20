-- Knowledge Module - SQL Functions & Procedures
-- Handles documents, chunks, and chunk statistics for RAG system

-- =====================================================
-- DOCUMENTS SECTION
-- =====================================================

-- =====================================================
-- Function: fn_get_all_documents
-- Description: Retrieves all active documents with pagination
-- =====================================================
create or replace function fn_get_all_documents(
    p_limit int default 100,
    p_offset int default 0
)
returns table (
    doc_id int,
    doc_category varchar,
    doc_title varchar,
    doc_summary text,
    doc_source varchar,
    doc_published_at timestamp,
    doc_active boolean,
    doc_created_at timestamp,
    doc_updated_at timestamp
) as $$
begin
    return query
    select
        d.doc_id,
        d.doc_category,
        d.doc_title,
        d.doc_summary,
        d.doc_source,
        d.doc_published_at,
        d.doc_active,
        d.doc_created_at,
        d.doc_updated_at
    from public.cht_documents d
    where d.doc_active = true
    order by d.doc_created_at desc
    limit p_limit
    offset p_offset;
end;
$$ language plpgsql;

-- =====================================================
-- Function: fn_get_document_by_id
-- Description: Get specific document by ID
-- =====================================================
create or replace function fn_get_document_by_id(
    p_doc_id int
)
returns table (
    doc_id int,
    doc_category varchar,
    doc_title varchar,
    doc_summary text,
    doc_source varchar,
    doc_published_at timestamp,
    doc_active boolean,
    doc_created_at timestamp,
    doc_updated_at timestamp
) as $$
begin
    return query
    select
        d.doc_id,
        d.doc_category,
        d.doc_title,
        d.doc_summary,
        d.doc_source,
        d.doc_published_at,
        d.doc_active,
        d.doc_created_at,
        d.doc_updated_at
    from public.cht_documents d
    where d.doc_id = p_doc_id
    and d.doc_active = true;
end;
$$ language plpgsql;

-- =====================================================
-- Function: fn_get_documents_by_category
-- Description: Get documents filtered by category
-- =====================================================
create or replace function fn_get_documents_by_category(
    p_category varchar,
    p_limit int default 100,
    p_offset int default 0
)
returns table (
    doc_id int,
    doc_category varchar,
    doc_title varchar,
    doc_summary text,
    doc_source varchar,
    doc_published_at timestamp,
    doc_active boolean,
    doc_created_at timestamp,
    doc_updated_at timestamp
) as $$
begin
    return query
    select
        d.doc_id,
        d.doc_category,
        d.doc_title,
        d.doc_summary,
        d.doc_source,
        d.doc_published_at,
        d.doc_active,
        d.doc_created_at,
        d.doc_updated_at
    from public.cht_documents d
    where d.doc_category = p_category
    and d.doc_active = true
    order by d.doc_created_at desc
    limit p_limit
    offset p_offset;
end;
$$ language plpgsql;

-- =====================================================
-- Function: fn_search_documents_by_title
-- Description: Search documents by title pattern
-- =====================================================
create or replace function fn_search_documents_by_title(
    p_title_pattern varchar,
    p_limit int default 100
)
returns table (
    doc_id int,
    doc_category varchar,
    doc_title varchar,
    doc_summary text,
    doc_source varchar,
    doc_published_at timestamp,
    doc_active boolean,
    doc_created_at timestamp,
    doc_updated_at timestamp
) as $$
begin
    return query
    select
        d.doc_id,
        d.doc_category,
        d.doc_title,
        d.doc_summary,
        d.doc_source,
        d.doc_published_at,
        d.doc_active,
        d.doc_created_at,
        d.doc_updated_at
    from public.cht_documents d
    where d.doc_title ilike '%' || p_title_pattern || '%'
    and d.doc_active = true
    order by d.doc_created_at desc
    limit p_limit;
end;
$$ language plpgsql;

-- =====================================================
-- Procedure: sp_create_document
-- Description: Creates a new document
-- Returns: success (boolean), code (varchar), doc_id (int)
-- =====================================================
create or replace procedure sp_create_document(
    out success boolean,
    out code varchar,
    out o_doc_id int,
    in p_category varchar,
    in p_title varchar,
    in p_summary text,
    in p_source varchar,
    in p_published_at timestamp
)
language plpgsql
as $$
begin
    success := true;
    code := 'OK';
    o_doc_id := null;

    -- Validate required fields
    if p_category is null or p_title is null then
        success := false;
        code := 'ERR_REQUIRED_FIELDS';
        return;
    end if;

    -- Insert new document
    insert into public.cht_documents (
        doc_category,
        doc_title,
        doc_summary,
        doc_source,
        doc_published_at,
        doc_active
    ) values (
        p_category,
        p_title,
        p_summary,
        p_source,
        p_published_at,
        true
    )
    returning doc_id into o_doc_id;

exception
    when others then
        success := false;
        code := 'ERR_CREATE_DOCUMENT';
        raise notice 'Error creating document: %', sqlerrm;
end;
$$;

-- =====================================================
-- Procedure: sp_update_document
-- Description: Updates an existing document
-- Returns: success (boolean), code (varchar)
-- =====================================================
create or replace procedure sp_update_document(
    out success boolean,
    out code varchar,
    in p_doc_id int,
    in p_category varchar,
    in p_title varchar,
    in p_summary text,
    in p_source varchar,
    in p_published_at timestamp
)
language plpgsql
as $$
declare
    v_exists boolean;
begin
    success := true;
    code := 'OK';

    -- Check if document exists
    select exists(
        select 1
        from public.cht_documents
        where doc_id = p_doc_id
        and doc_active = true
    ) into v_exists;

    if not v_exists then
        success := false;
        code := 'ERR_DOCUMENT_NOT_FOUND';
        return;
    end if;

    -- Update document
    update public.cht_documents
    set
        doc_category = p_category,
        doc_title = p_title,
        doc_summary = p_summary,
        doc_source = p_source,
        doc_published_at = p_published_at
    where doc_id = p_doc_id;

exception
    when others then
        success := false;
        code := 'ERR_UPDATE_DOCUMENT';
        raise notice 'Error updating document: %', sqlerrm;
end;
$$;

-- =====================================================
-- Procedure: sp_delete_document
-- Description: Soft delete document (cascades to chunks)
-- Returns: success (boolean), code (varchar)
-- =====================================================
create or replace procedure sp_delete_document(
    out success boolean,
    out code varchar,
    in p_doc_id int
)
language plpgsql
as $$
declare
    v_exists boolean;
begin
    success := true;
    code := 'OK';

    -- Check if document exists
    select exists(
        select 1
        from public.cht_documents
        where doc_id = p_doc_id
        and doc_active = true
    ) into v_exists;

    if not v_exists then
        success := false;
        code := 'ERR_DOCUMENT_NOT_FOUND';
        return;
    end if;

    -- Soft delete (note: chunks reference this, so they stay but document is inactive)
    update public.cht_documents
    set doc_active = false
    where doc_id = p_doc_id;

exception
    when others then
        success := false;
        code := 'ERR_DELETE_DOCUMENT';
        raise notice 'Error deleting document: %', sqlerrm;
end;
$$;

-- =====================================================
-- CHUNKS SECTION
-- =====================================================

-- =====================================================
-- Function: fn_get_chunks_by_document
-- Description: Get all chunks for a specific document
-- =====================================================
create or replace function fn_get_chunks_by_document(
    p_doc_id int
)
returns table (
    chk_id int,
    chk_fk_document int,
    chk_content text,
    chk_created_at timestamp,
    chk_updated_at timestamp
) as $$
begin
    return query
    select
        c.chk_id,
        c.chk_fk_document,
        c.chk_content,
        c.chk_created_at,
        c.chk_updated_at
    from public.cht_chunks c
    where c.chk_fk_document = p_doc_id
    order by c.chk_id;
end;
$$ language plpgsql;

-- =====================================================
-- Function: fn_get_chunk_by_id
-- Description: Get specific chunk by ID with embedding
-- =====================================================
create or replace function fn_get_chunk_by_id(
    p_chk_id int
)
returns table (
    chk_id int,
    chk_fk_document int,
    chk_content text,
    chk_embedding vector(1536),
    chk_created_at timestamp,
    chk_updated_at timestamp
) as $$
begin
    return query
    select
        c.chk_id,
        c.chk_fk_document,
        c.chk_content,
        c.chk_embedding,
        c.chk_created_at,
        c.chk_updated_at
    from public.cht_chunks c
    where c.chk_id = p_chk_id;
end;
$$ language plpgsql;

-- =====================================================
-- Function: fn_similarity_search_chunks
-- Description: Vector similarity search for RAG
-- Returns chunks ordered by cosine similarity to query embedding
-- =====================================================
create or replace function fn_similarity_search_chunks(
    p_query_embedding vector(1536),
    p_limit int default 5,
    p_min_similarity float default 0.7
)
returns table (
    chk_id int,
    chk_fk_document int,
    chk_content text,
    similarity_score float,
    doc_title varchar,
    doc_category varchar
) as $$
begin
    return query
    select
        c.chk_id,
        c.chk_fk_document,
        c.chk_content,
        1 - (c.chk_embedding <=> p_query_embedding) as similarity_score,
        d.doc_title,
        d.doc_category
    from public.cht_chunks c
    inner join public.cht_documents d on c.chk_fk_document = d.doc_id
    where d.doc_active = true
    and c.chk_embedding is not null
    and (1 - (c.chk_embedding <=> p_query_embedding)) >= p_min_similarity
    order by c.chk_embedding <=> p_query_embedding
    limit p_limit;
end;
$$ language plpgsql;

-- =====================================================
-- Procedure: sp_create_chunk
-- Description: Creates a new chunk with optional embedding
-- Returns: success (boolean), code (varchar), chk_id (int)
-- =====================================================
create or replace procedure sp_create_chunk(
    out success boolean,
    out code varchar,
    out o_chk_id int,
    in p_doc_id int,
    in p_content text,
    in p_embedding vector(1536)
)
language plpgsql
as $$
declare
    v_doc_exists boolean;
begin
    success := true;
    code := 'OK';
    o_chk_id := null;

    -- Validate document exists
    select exists(
        select 1
        from public.cht_documents
        where doc_id = p_doc_id
        and doc_active = true
    ) into v_doc_exists;

    if not v_doc_exists then
        success := false;
        code := 'ERR_DOCUMENT_NOT_FOUND';
        return;
    end if;

    -- Insert chunk
    insert into public.cht_chunks (
        chk_fk_document,
        chk_content,
        chk_embedding
    ) values (
        p_doc_id,
        p_content,
        p_embedding
    )
    returning chk_id into o_chk_id;

    -- Initialize statistics record for this chunk
    insert into public.cht_chunk_statistics (
        cst_fk_chunk,
        cst_usage_count
    ) values (
        o_chk_id,
        0
    );

exception
    when others then
        success := false;
        code := 'ERR_CREATE_CHUNK';
        raise notice 'Error creating chunk: %', sqlerrm;
end;
$$;

-- =====================================================
-- Procedure: sp_update_chunk_embedding
-- Description: Updates the embedding for a chunk
-- Returns: success (boolean), code (varchar)
-- =====================================================
create or replace procedure sp_update_chunk_embedding(
    out success boolean,
    out code varchar,
    in p_chk_id int,
    in p_embedding vector(1536)
)
language plpgsql
as $$
declare
    v_exists boolean;
begin
    success := true;
    code := 'OK';

    -- Check if chunk exists
    select exists(
        select 1
        from public.cht_chunks
        where chk_id = p_chk_id
    ) into v_exists;

    if not v_exists then
        success := false;
        code := 'ERR_CHUNK_NOT_FOUND';
        return;
    end if;

    -- Update embedding
    update public.cht_chunks
    set chk_embedding = p_embedding
    where chk_id = p_chk_id;

exception
    when others then
        success := false;
        code := 'ERR_UPDATE_CHUNK_EMBEDDING';
        raise notice 'Error updating chunk embedding: %', sqlerrm;
end;
$$;

-- =====================================================
-- Procedure: sp_delete_chunk
-- Description: Hard delete chunk (cascades to statistics)
-- Returns: success (boolean), code (varchar)
-- =====================================================
create or replace procedure sp_delete_chunk(
    out success boolean,
    out code varchar,
    in p_chk_id int
)
language plpgsql
as $$
declare
    v_exists boolean;
begin
    success := true;
    code := 'OK';

    -- Check if chunk exists
    select exists(
        select 1
        from public.cht_chunks
        where chk_id = p_chk_id
    ) into v_exists;

    if not v_exists then
        success := false;
        code := 'ERR_CHUNK_NOT_FOUND';
        return;
    end if;

    -- Hard delete (cascades to statistics via FK)
    delete from public.cht_chunks
    where chk_id = p_chk_id;

exception
    when others then
        success := false;
        code := 'ERR_DELETE_CHUNK';
        raise notice 'Error deleting chunk: %', sqlerrm;
end;
$$;

-- =====================================================
-- Procedure: sp_bulk_create_chunks
-- Description: Creates multiple chunks for a document at once
-- Returns: success (boolean), code (varchar), chunks_created (int)
-- =====================================================
create or replace procedure sp_bulk_create_chunks(
    out success boolean,
    out code varchar,
    out o_chunks_created int,
    in p_doc_id int,
    in p_contents text[],
    in p_embeddings vector(1536)[]
)
language plpgsql
as $$
declare
    v_doc_exists boolean;
    v_content text;
    v_embedding vector(1536);
    i_index int;
    v_new_chk_id int;
begin
    success := true;
    code := 'OK';
    o_chunks_created := 0;

    -- Validate document exists
    select exists(
        select 1
        from public.cht_documents
        where doc_id = p_doc_id
        and doc_active = true
    ) into v_doc_exists;

    if not v_doc_exists then
        success := false;
        code := 'ERR_DOCUMENT_NOT_FOUND';
        return;
    end if;

    -- Loop through contents and create chunks
    for i_index in 1..array_length(p_contents, 1)
    loop
        v_content := p_contents[i_index];

        -- Get embedding if provided
        if p_embeddings is not null and array_length(p_embeddings, 1) >= i_index then
            v_embedding := p_embeddings[i_index];
        else
            v_embedding := null;
        end if;

        -- Insert chunk
        insert into public.cht_chunks (
            chk_fk_document,
            chk_content,
            chk_embedding
        ) values (
            p_doc_id,
            v_content,
            v_embedding
        )
        returning chk_id into v_new_chk_id;

        -- Initialize statistics
        insert into public.cht_chunk_statistics (
            cst_fk_chunk,
            cst_usage_count
        ) values (
            v_new_chk_id,
            0
        );

        o_chunks_created := o_chunks_created + 1;
    end loop;

exception
    when others then
        success := false;
        code := 'ERR_BULK_CREATE_CHUNKS';
        o_chunks_created := 0;
        raise notice 'Error bulk creating chunks: %', sqlerrm;
end;
$$;

-- =====================================================
-- CHUNK STATISTICS SECTION
-- =====================================================

-- =====================================================
-- Function: fn_get_chunk_statistics
-- Description: Get statistics for a specific chunk
-- =====================================================
create or replace function fn_get_chunk_statistics(
    p_chk_id int
)
returns table (
    cst_id int,
    cst_fk_chunk int,
    cst_usage_count int,
    cst_last_used_at timestamp,
    cst_precision_atk float,
    cst_recall_atk float,
    cst_f1_atk float,
    cst_mrr float,
    cst_map float,
    cst_ndcg float,
    cst_staleness_days int,
    cst_last_refresh_at timestamp,
    cst_curriculum_coverage_pct float
) as $$
begin
    return query
    select
        cs.cst_id,
        cs.cst_fk_chunk,
        cs.cst_usage_count,
        cs.cst_last_used_at,
        cs.cst_precision_atk,
        cs.cst_recall_atk,
        cs.cst_f1_atk,
        cs.cst_mrr,
        cs.cst_map,
        cs.cst_ndcg,
        cs.cst_staleness_days,
        cs.cst_last_refresh_at,
        cs.cst_curriculum_coverage_pct
    from public.cht_chunk_statistics cs
    where cs.cst_fk_chunk = p_chk_id;
end;
$$ language plpgsql;

-- =====================================================
-- Function: fn_get_top_chunks_by_usage
-- Description: Get most frequently used chunks
-- =====================================================
create or replace function fn_get_top_chunks_by_usage(
    p_limit int default 10
)
returns table (
    chk_id int,
    chk_content text,
    doc_title varchar,
    usage_count int,
    last_used_at timestamp,
    f1_score float
) as $$
begin
    return query
    select
        c.chk_id,
        c.chk_content,
        d.doc_title,
        cs.cst_usage_count,
        cs.cst_last_used_at,
        cs.cst_f1_atk
    from public.cht_chunk_statistics cs
    inner join public.cht_chunks c on cs.cst_fk_chunk = c.chk_id
    inner join public.cht_documents d on c.chk_fk_document = d.doc_id
    where d.doc_active = true
    order by cs.cst_usage_count desc
    limit p_limit;
end;
$$ language plpgsql;

-- =====================================================
-- Procedure: sp_increment_chunk_usage
-- Description: Increments usage count for a chunk
-- Returns: success (boolean), code (varchar)
-- =====================================================
create or replace procedure sp_increment_chunk_usage(
    out success boolean,
    out code varchar,
    in p_chk_id int
)
language plpgsql
as $$
declare
    v_exists boolean;
begin
    success := true;
    code := 'OK';

    -- Check if statistics record exists
    select exists(
        select 1
        from public.cht_chunk_statistics
        where cst_fk_chunk = p_chk_id
    ) into v_exists;

    if not v_exists then
        success := false;
        code := 'ERR_CHUNK_STATS_NOT_FOUND';
        return;
    end if;

    -- Increment usage and update last used timestamp
    update public.cht_chunk_statistics
    set
        cst_usage_count = cst_usage_count + 1,
        cst_last_used_at = current_timestamp
    where cst_fk_chunk = p_chk_id;

exception
    when others then
        success := false;
        code := 'ERR_INCREMENT_CHUNK_USAGE';
        raise notice 'Error incrementing chunk usage: %', sqlerrm;
end;
$$;

-- =====================================================
-- Procedure: sp_update_chunk_quality_metrics
-- Description: Updates quality metrics for a chunk
-- Returns: success (boolean), code (varchar)
-- =====================================================
create or replace procedure sp_update_chunk_quality_metrics(
    in p_chk_id int,
    out success boolean,
    out code varchar,
    in p_precision_atk float default null,
    in p_recall_atk float default null,
    in p_f1_atk float default null,
    in p_mrr float default null,
    in p_map float default null,
    in p_ndcg float default null
)
language plpgsql
as $$
declare
    v_exists boolean;
begin
    success := true;
    code := 'OK';

    -- Check if statistics record exists
    select exists(
        select 1
        from public.cht_chunk_statistics
        where cst_fk_chunk = p_chk_id
    ) into v_exists;

    if not v_exists then
        success := false;
        code := 'ERR_CHUNK_STATS_NOT_FOUND';
        return;
    end if;

    -- Update quality metrics (only non-null values)
    update public.cht_chunk_statistics
    set
        cst_precision_atk = coalesce(p_precision_atk, cst_precision_atk),
        cst_recall_atk = coalesce(p_recall_atk, cst_recall_atk),
        cst_f1_atk = coalesce(p_f1_atk, cst_f1_atk),
        cst_mrr = coalesce(p_mrr, cst_mrr),
        cst_map = coalesce(p_map, cst_map),
        cst_ndcg = coalesce(p_ndcg, cst_ndcg)
    where cst_fk_chunk = p_chk_id;

exception
    when others then
        success := false;
        code := 'ERR_UPDATE_CHUNK_METRICS';
        raise notice 'Error updating chunk metrics: %', sqlerrm;
end;
$$;

-- =====================================================
-- Procedure: sp_update_chunk_staleness
-- Description: Updates staleness tracking for a chunk
-- Returns: success (boolean), code (varchar)
-- =====================================================
create or replace procedure sp_update_chunk_staleness(
    out success boolean,
    out code varchar,
    in p_chk_id int,
    in p_staleness_days int
)
language plpgsql
as $$
declare
    v_exists boolean;
begin
    success := true;
    code := 'OK';

    -- Check if statistics record exists
    select exists(
        select 1
        from public.cht_chunk_statistics
        where cst_fk_chunk = p_chk_id
    ) into v_exists;

    if not v_exists then
        success := false;
        code := 'ERR_CHUNK_STATS_NOT_FOUND';
        return;
    end if;

    -- Update staleness
    update public.cht_chunk_statistics
    set
        cst_staleness_days = p_staleness_days,
        cst_last_refresh_at = current_timestamp
    where cst_fk_chunk = p_chk_id;

exception
    when others then
        success := false;
        code := 'ERR_UPDATE_CHUNK_STALENESS';
        raise notice 'Error updating chunk staleness: %', sqlerrm;
end;
$$;

-- =====================================================
-- Comments
-- =====================================================

-- Documents
comment on function fn_get_all_documents(int, int) is 'Retrieves all active documents with pagination';
comment on function fn_get_document_by_id(int) is 'Get specific document by ID';
comment on function fn_get_documents_by_category(varchar, int, int) is 'Get documents filtered by category with pagination';
comment on function fn_search_documents_by_title(varchar, int) is 'Search documents by title pattern using ILIKE';
comment on procedure sp_create_document is 'Creates a new document. Returns success, code, and doc_id';
comment on procedure sp_update_document is 'Updates an existing document. Returns success and code';
comment on procedure sp_delete_document is 'Soft deletes a document. Returns success and code';

-- Chunks
comment on function fn_get_chunks_by_document(int) is 'Get all chunks for a specific document';
comment on function fn_get_chunk_by_id(int) is 'Get specific chunk by ID including embedding';
comment on function fn_similarity_search_chunks(vector, int, float) is 'Vector similarity search for RAG - returns top K similar chunks';
comment on procedure sp_create_chunk is 'Creates a new chunk with optional embedding. Returns success, code, and chk_id';
comment on procedure sp_update_chunk_embedding is 'Updates the embedding vector for a chunk. Returns success and code';
comment on procedure sp_delete_chunk is 'Hard deletes a chunk (cascades to statistics). Returns success and code';
comment on procedure sp_bulk_create_chunks is 'Bulk creates multiple chunks for a document. Returns success, code, and count';

-- Chunk Statistics
comment on function fn_get_chunk_statistics(int) is 'Get all statistics for a specific chunk';
comment on function fn_get_top_chunks_by_usage(int) is 'Get most frequently used chunks';
comment on procedure sp_increment_chunk_usage is 'Increments usage count and updates last used timestamp. Returns success and code';
comment on procedure sp_update_chunk_quality_metrics is 'Updates RAG quality metrics (precision, recall, F1, MRR, MAP, NDCG). Returns success and code';
comment on procedure sp_update_chunk_staleness is 'Updates staleness tracking for content freshness. Returns success and code';
