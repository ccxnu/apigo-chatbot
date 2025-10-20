 Guía de  Prácticas para Scripts SQL

---

## Reglas de Nomenclatura

### Variables en SQL/PLpgSQL

- **Parámetros de entrada:**
  Prefijo `p_`
  Ejemplo: `p_id_usuario`, `p_estado`
- **Variables de salida:**
  Prefijo `o_`
  Ejemplo: `o_id_parametro`
  Espaciales: success boolean, code varchar
- **Variables locales, declaradas:**
  Prefijo `v_`
  Ejemplo: `v_mod_id`, `v_pram_id`
- **Registros temporales:**
  Prefijo `r_`
  Ejemplo: `r_modulo`
- **Contadores:**
  Prefijo `i_`
  Ejemplo: `i_contador`
- **Booleanos:**
  Prefijo `bl_` o `is_`
  Ejemplo: `bl_existe`, `is_active`
- **Constantes:**
  Prefijo `c_`
  Ejemplo: `c_estado_activo`

### Funciones y Procedimientos

- Usa prefijos:
  - **sp_** para procedimientos (ej. `sp_add_usuario`)
  - **fn_** para funciones (ej. `fn_get_parametros`)
  - **vw_** para vistas (ej. `vw_usuarios_permisos`)
  - **tr_** para triggers (ej. `vw_usuarios_permisos`)
- Los nombres deben ser descriptivos, usar inglés o español de forma consistente.
- Mantener `snake_case` en todos los identificadores.
- Los procedimientos siempre retornan:
  - success: boolean
  - code: varchar (Es el código unico del error, ejm. 'OK', "ERR_NOT_FOUND')

---

## Convenciones Generales

- **Siempre comentar** cada bloque relevante (`-- Comentario explicativo`)
- **Indentación de 4 espacios**, tabs.
- Palabras clave SQL en **minúsculas**.
- Usar `do $$ ... $$;` solo cuando es necesario (lógica condicional, variables, etc).
- Si el script puede ejecutarse múltiples veces, usa `if not exists` para inserts y creaciones.
- **No uses abreviaturas ambiguas.**
- Scripts deben ser **idempotentes** (repetirlos no debe causar errores ni duplicados).

---

## Ejemplos de Script

Ejemplo de agregar.

```sql
do $$
declare
    v_mod_id      int;
    v_prm_id      int;
begin
    -- Obtener ID del módulo SISTEMA
    select mod_id into v_mod_id from cht_modulos where mod_codigo = 'SISTEMA';

    -- Siguiente prm_id disponible
    select coalesce(max(prm_id), 0) + 1 into v_prm_id from cht_parametros;

    -- Insertar parámetro si no existe
    if not exists (select 1 from cht_parametros where prm_nemonico = 'COD_OK') then
        insert into cht_parametros (
            prm_id, prm_fk_modulo, prm_nombre, prm_nemonico, prm_valor1, prm_valor2, prm_descripcion, prm_tipo, prm_activo
        ) values (
            v_prm_id, v_mod_id, 'Operación exitosa', 'COD_OK',
            'Operación realizada correctamente.', '', 'Mensaje estándar de éxito para todas las operaciones.', 'string', true
        );
    end if;
end
$$;
```

Ejemplo de función para actualizar el campo updated_at automáticamente.

```sql
-- Ejemplo adaptado de: https://github.com/pg-nano/pg-nano/blob/master/demos/exhaustive/sql/blog/post.pgsql
-- Modificado para seguir las convenciones del proyecto.
create or replace function fn_set_updated_at()
returns trigger
language plpgsql
as $$
begin
    new.updated_at := current_timestamp;
    return new;
end;
$$;

-- Trigger que llama a la función antes de cada UPDATE en la tabla post
create trigger tr_update_post
before update on post
for each row
execute function fn_set_updated_at();
```
