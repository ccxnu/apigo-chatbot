-- Parameters Module - SQL Functions & Procedures

-- =====================================================
-- Function: fn_get_all_parameters
-- Description: Retrieves all active parameters
-- =====================================================
create or replace function fn_get_all_parameters()
returns table (
    prm_id          int,
    prm_name        varchar,
    prm_code        varchar,
    prm_data        jsonb,
    prm_description varchar,
    prm_active      boolean,
    prm_created_at  timestamp,
    prm_updated_at  timestamp
) as $$
begin
    return query
    select
        p.prm_id,
        p.prm_name,
        p.prm_code,
        p.prm_data,
        p.prm_description,
        p.prm_active,
        p.prm_created_at,
        p.prm_updated_at
    from public.cht_parameters p
    where p.prm_active = true
    order by p.prm_name;
end;
$$
language plpgsql;

-- =====================================================
-- Function: fn_get_parameter_by_code
-- Description: Get specific parameter by code
-- =====================================================
create or replace function fn_get_parameter_by_code(
    p_code varchar
)
returns table (
    prm_id int,
    prm_name varchar,
    prm_code varchar,
    prm_data jsonb,
    prm_description varchar,
    prm_active boolean,
    prm_created_at timestamp,
    prm_updated_at timestamp
) as $$
begin
    return query
    select
        p.prm_id,
        p.prm_name,
        p.prm_code,
        p.prm_data,
        p.prm_description,
        p.prm_active,
        p.prm_created_at,
        p.prm_updated_at
    from public.cht_parameters p
    where p.prm_code = p_code
    and p.prm_active = true;
end;
$$ language plpgsql;

-- =====================================================
-- Procedure: sp_create_parameter
-- Description: Creates a new parameter
-- Returns: success (boolean), code (varchar)
-- =====================================================
create or replace procedure sp_create_parameter(
    out success boolean,
    out code varchar,
    in p_name varchar,
    in p_code varchar,
    in p_data jsonb,
    in p_description varchar
)
language plpgsql
as $$
declare
    v_exists boolean;

begin
    success := true;
    code := 'OK';

    -- Check if code already exists
    select exists(
        select 1
        from public.cht_parameters
        where prm_code = p_code
    ) into v_exists;

    if v_exists then
        success := false;
        code := 'ERR_PARAM_CODE_EXISTS';
        return;
    end if;

    -- Insert new parameter
    insert into public.cht_parameters (
        prm_name,
        prm_code,
        prm_data,
        prm_description,
        prm_active
    ) values (
        p_name,
        p_code,
        p_data,
        p_description,
        true
    );

exception
    when others then
        success := false;
        code := 'ERR_CREATE_PARAMETER';
        raise notice 'Error creating parameter: %', sqlerrm;
end;
$$;

-- =====================================================
-- Procedure: sp_update_parameter
-- Description: Updates an existing parameter
-- Returns: success (boolean), code (varchar)
-- =====================================================
create or replace procedure sp_update_parameter(
    out success boolean,
    out code varchar,
    in p_code varchar,
    in p_name varchar,
    in p_data jsonb,
    in p_description varchar
)
language plpgsql
as $$
declare
    v_prm_id int;

begin
    success := true;
    code := 'OK';

    -- Check if parameter exists
    select prm_id into v_prm_id
    from public.cht_parameters
    where prm_code = p_code
    and prm_active = true;

    if v_prm_id is null then
        success := false;
        code := 'ERR_PARAM_NOT_FOUND';
        return;
    end if;

    -- Update parameter
    update public.cht_parameters
    set
        prm_name = p_name,
        prm_data = p_data,
        prm_description = p_description
    where prm_id = v_prm_id;

exception
    when others then
        success := false;
        code := 'ERR_UPDATE_PARAMETER';
        raise notice 'Error updating parameter: %', sqlerrm;
end;
$$;

-- =====================================================
-- Procedure: sp_delete_parameter
-- Description: Soft delete parameter (set active = false)
-- Returns: success (boolean), code (varchar)
-- =====================================================
create or replace procedure sp_delete_parameter(
    out success boolean,
    out code varchar,
    in p_code varchar
)
language plpgsql
as $$
declare
    v_prm_id int;

begin
    success := true;
    code := 'OK';

    -- Check if parameter exists
    select prm_id into v_prm_id
    from public.cht_parameters
    where prm_code = p_code
    and prm_active = true;

    if v_prm_id is null then
        success := false;
        code := 'ERR_PARAM_NOT_FOUND';
        return;
    end if;

    -- Soft delete
    update public.cht_parameters
    set prm_active = false
    where prm_id = v_prm_id;

exception
    when others then
        success := false;
        code := 'ERR_DELETE_PARAMETER';
        raise notice 'Error deleting parameter: %', sqlerrm;
end;
$$;

-- =====================================================
-- Function: fn_get_parameter_value
-- Description: Get only the data value by code
-- =====================================================
create or replace function fn_get_parameter_value(
    p_code varchar
)
returns jsonb as $$
declare
    v_data jsonb;
begin
    select prm_data into v_data
    from public.cht_parameters
    where prm_code = p_code
    and prm_active = true;

    return coalesce(v_data, '{}'::jsonb);
end;
$$ language plpgsql;

-- =====================================================
-- Function: fn_get_parameters_by_name
-- Description: Search parameters by name pattern
-- =====================================================
create or replace function fn_get_parameters_by_name(
    p_name_pattern varchar
)
returns table (
    prm_id int,
    prm_name varchar,
    prm_code varchar,
    prm_data jsonb,
    prm_description varchar,
    prm_active boolean,
    prm_created_at timestamp,
    prm_updated_at timestamp
) as $$
begin
    return query
    select
        p.prm_id,
        p.prm_name,
        p.prm_code,
        p.prm_data,
        p.prm_description,
        p.prm_active,
        p.prm_created_at,
        p.prm_updated_at
    from public.cht_parameters p
    where p.prm_name ilike '%' || p_name_pattern || '%'
    and p.prm_active = true
    order by p.prm_name;
end;
$$ language plpgsql;

-- =====================================================
-- Function: fn_get_parameters_by_codes
-- Description: Get multiple parameters by their codes (bulk)
-- =====================================================
create or replace function fn_get_parameters_by_codes(
    p_codes varchar[]
)
returns table (
    prm_id int,
    prm_name varchar,
    prm_code varchar,
    prm_data jsonb,
    prm_description varchar,
    prm_active boolean,
    prm_created_at timestamp,
    prm_updated_at timestamp
) as $$
begin
    return query
    select
        p.prm_id,
        p.prm_name,
        p.prm_code,
        p.prm_data,
        p.prm_description,
        p.prm_active,
        p.prm_created_at,
        p.prm_updated_at
    from public.cht_parameters p
    where p.prm_code = any(p_codes)
    and p.prm_active = true
    order by p.prm_name;
end;
$$ language plpgsql;

-- Comments
comment on function fn_get_all_parameters() is 'Retrieves all active system parameters';
comment on function fn_get_parameter_by_code(varchar) is 'Get specific parameter by unique code';
comment on function fn_get_parameter_value(varchar) is 'Get only the JSON data value of a parameter';
comment on function fn_get_parameters_by_name(varchar) is 'Search parameters by name pattern using ILIKE';
comment on function fn_get_parameters_by_codes(varchar[]) is 'Get multiple parameters at once by array of codes';
comment on procedure sp_create_parameter is 'Creates a new system parameter. Returns success and code';
comment on procedure sp_update_parameter is 'Updates an existing parameter. Returns success and code';
comment on procedure sp_delete_parameter is 'Soft deletes a parameter. Returns success and code';
