

# Oracle Data Hygiene: English Character Extraction

This guide covers the most efficient ways to retrieve only English letters ($A-Z$, $a-z$) and numbers ($0-9$) from an Oracle Database.

---

## 1. The Core Methods

### A. The `REGEXP_REPLACE` Method
Best for readability and one-off queries.
```sql
SELECT REGEXP_REPLACE(column_name, '[^a-zA-Z0-9]', '') AS clean_data
FROM your_table;
```

### B. The `TRANSLATE` Method
Best for high-performance and large datasets. It is significantly faster than Regex because it avoids the overhead of the regular expression engine.
```sql
SELECT TRANSLATE(column_name, 
                 'abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789' || column_name, 
                 'abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789') AS clean_data
FROM your_table;
```

---

## 2. Automating the Process

### The Reusable Function
Instead of typing the long string of letters every time, create this function once:

```sql
CREATE OR REPLACE FUNCTION clean_alphanumeric(p_str IN VARCHAR2) 
RETURN VARCHAR2 DETERMINISTIC IS
    v_keep VARCHAR2(62) := 'abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789';
BEGIN
    RETURN TRANSLATE(p_str, v_keep || p_str, v_keep);
END;
/
```

### The "Auto-View" Generator
Run this script to generate a `CREATE VIEW` statement for any table. It automatically applies the cleanup to every `VARCHAR2` and `CHAR` column by querying the Data Dictionary. High-Speed parallel indexes and a series of GRANT statements in one go.

```sql
SELECT 
    DBMS_XMLGEN.CONVERT(
        XMLAGG(
            XMLELEMENT(e, 
                -- View Definition (Mirroring all columns)
                'CREATE OR REPLACE VIEW v_clean_' || t.table_name || ' AS SELECT ' || 
                (SELECT LISTAGG(
                    CASE 
                        WHEN data_type IN ('VARCHAR2', 'CHAR') 
                        THEN 'TRANSLATE(' || column_name || ', CHR(39) || ''abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789 '' || ' || column_name || ', CHR(39) || ''abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789 '')'
                        ELSE column_name 
                    END || ' AS ' || column_name, ', ')
                 WITHIN GROUP (ORDER BY column_id)
                 FROM all_tab_columns 
                 WHERE table_name = t.table_name) || 
                ' FROM ' || t.table_name || ';' || CHR(10) ||
                
                -- High-Speed Parallel Indexes
                (SELECT LISTAGG('CREATE INDEX idx_clean_' || ic.table_name || '_' || SUBSTR(ic.column_name,1,10) || ' ON ' || ic.table_name || 
                 '(TRANSLATE(' || ic.column_name || ', CHR(39) || ''abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789 '' || ' || ic.column_name || ', CHR(39) || ''abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789 '')) PARALLEL 4;' || 
                 CHR(10) || 'ALTER INDEX idx_clean_' || ic.table_name || '_' || SUBSTR(ic.column_name,1,10) || ' NOPARALLEL;', CHR(10))
                 WITHIN GROUP (ORDER BY ic.column_name)
                 FROM all_ind_columns ic
                 JOIN all_tab_columns tc ON ic.table_name = tc.table_name AND ic.column_name = tc.column_name
                 WHERE ic.table_name = t.table_name 
                   AND tc.data_type IN ('VARCHAR2', 'CHAR')
                   AND ic.column_position = 1) || CHR(10) ||

                -- Grants
                (SELECT LISTAGG('GRANT SELECT ON v_clean_' || table_name || ' TO ' || grantee || ';', CHR(10))
                 WITHIN GROUP (ORDER BY grantee)
                 FROM all_tab_privs 
                 WHERE table_name = t.table_name AND privilege = 'SELECT')
            ).EXTRACT('//text()')
        ).GETCLOBVAL(), 1) AS deploy_script
FROM (SELECT DISTINCT table_name FROM all_tab_columns WHERE table_name IN ('YOUR_TABLE_NAME')) t
GROUP BY table_name;
```

---
## 2.5. The "Auto-View" Generator + Special Characters

```sql
SELECT 
    DBMS_XMLGEN.CONVERT(
        XMLAGG(
            XMLELEMENT(e, 
                -- View Definition
                'CREATE OR REPLACE VIEW v_clean_' || t.table_name || ' AS SELECT ' || 
                (SELECT LISTAGG(
                    CASE 
                        WHEN data_type IN ('VARCHAR2', 'CHAR') 
                        THEN 'TRANSLATE(' || column_name || ', CHR(39) || ''abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789 !@#$%^&*()-_=+[]{}|;:,.<>/?'' || ' || column_name || ', CHR(39) || ''abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789 !@#$%^&*()-_=+[]{}|;:,.<>/?'')'
                        ELSE column_name 
                    END || ' AS ' || column_name, ', ')
                 WITHIN GROUP (ORDER BY column_id)
                 FROM all_tab_columns 
                 WHERE table_name = t.table_name) || 
                ' FROM ' || t.table_name || ';' || CHR(10) ||
                
                -- High-Speed Parallel Indexes
                (SELECT LISTAGG('CREATE INDEX idx_clean_' || ic.table_name || '_' || SUBSTR(ic.column_name,1,10) || ' ON ' || ic.table_name || 
                 '(TRANSLATE(' || ic.column_name || ', CHR(39) || ''abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789 !@#$%^&*()-_=+[]{}|;:,.<>/?'' || ' || ic.column_name || ', CHR(39) || ''abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789 !@#$%^&*()-_=+[]{}|;:,.<>/?'')) PARALLEL 4;' || 
                 CHR(10) || 'ALTER INDEX idx_clean_' || ic.table_name || '_' || SUBSTR(ic.column_name,1,10) || ' NOPARALLEL;', CHR(10))
                 WITHIN GROUP (ORDER BY ic.column_name)
                 FROM all_ind_columns ic
                 JOIN all_tab_columns tc ON ic.table_name = tc.table_name AND ic.column_name = tc.column_name
                 WHERE ic.table_name = t.table_name 
                   AND tc.data_type IN ('VARCHAR2', 'CHAR')
                   AND ic.column_position = 1) || CHR(10) ||

                -- Grants
                (SELECT LISTAGG('GRANT SELECT ON v_clean_' || table_name || ' TO ' || grantee || ';', CHR(10))
                 WITHIN GROUP (ORDER BY grantee)
                 FROM all_tab_privs 
                 WHERE table_name = t.table_name AND privilege = 'SELECT')
            ).EXTRACT('//text()')
        ).GETCLOBVAL(), 1) AS deploy_script
FROM (SELECT DISTINCT table_name FROM all_tab_columns WHERE table_name IN ('YOUR_TABLE_NAME')) t
GROUP BY table_name;
```


---
## 2.5.1. Rollback Script

```sql
SELECT 
    DBMS_XMLGEN.CONVERT(
        XMLAGG(
            XMLELEMENT(e, 
                'DROP VIEW v_clean_' || t.table_name || ';' || CHR(10) ||
                (SELECT LISTAGG('DROP INDEX idx_clean_' || ic.table_name || '_' || SUBSTR(ic.column_name,1,10) || ';', CHR(10))
                 WITHIN GROUP (ORDER BY ic.column_name)
                 FROM all_ind_columns ic
                 JOIN all_tab_columns tc ON ic.table_name = tc.table_name AND ic.column_name = tc.column_name
                 WHERE ic.table_name = t.table_name 
                   AND tc.data_type IN ('VARCHAR2', 'CHAR')
                   AND ic.column_position = 1)
            ).EXTRACT('//text()')
        ).GETCLOBVAL(), 1) AS undo_script
FROM (SELECT DISTINCT table_name FROM all_tab_columns WHERE table_name IN ('YOUR_TABLE_NAME')) t
GROUP BY table_name;
```

## 3. Performance Summary

| Feature | `REGEXP_REPLACE` | `TRANSLATE` |
| :--- | :--- | :--- |
| **Speed** | Slower (Engine overhead) | **Very Fast** (Native mapping) |
| **Ease of Use** | High (Short syntax) | Moderate (Requires long strings) |
| **Flexibility** | High (Supports ranges) | Low (Must list every char) |

---

## 4. Best Practices & Tips

* **Case Standardization:** To force everything to Uppercase while cleaning, wrap the result: `UPPER(TRANSLATE(...))`.
* **Function-Based Indexes:** If you frequently query the "cleaned" version of a column, create a function-based index to prevent full table scans.
* **Accented Characters:** Note that characters like `é` or `ñ` will be removed by these methods. If you need to keep them, use the `[:alpha:]` regex class instead.

---
*Generated by Gemini*
