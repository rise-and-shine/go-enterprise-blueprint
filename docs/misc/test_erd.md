# ERP System ERD (Test Visualization)

This is a sample ERP system ERD for testing Mermaid rendering with ~50 tables.

```mermaid
erDiagram
    %% ==========================================
    %% CORE / COMPANY MANAGEMENT
    %% ==========================================

    companies {
        BIGSERIAL id PK
        VARCHAR name
        VARCHAR tax_id UK
        VARCHAR address
        VARCHAR phone
        VARCHAR email
        BOOLEAN is_active
        TIMESTAMPTZ created_at
        TIMESTAMPTZ updated_at
    }

    branches {
        BIGSERIAL id PK
        BIGINT company_id FK
        VARCHAR name
        VARCHAR code UK
        VARCHAR address
        VARCHAR phone
        BOOLEAN is_headquarters
        TIMESTAMPTZ created_at
    }

    departments {
        BIGSERIAL id PK
        BIGINT branch_id FK
        BIGINT parent_id FK
        VARCHAR name
        VARCHAR code
        TIMESTAMPTZ created_at
    }

    currencies {
        BIGSERIAL id PK
        VARCHAR code UK
        VARCHAR name
        VARCHAR symbol
        DECIMAL exchange_rate
        BOOLEAN is_default
    }

    fiscal_years {
        BIGSERIAL id PK
        BIGINT company_id FK
        VARCHAR name
        DATE start_date
        DATE end_date
        BOOLEAN is_closed
    }

    %% ==========================================
    %% HUMAN RESOURCES
    %% ==========================================

    employees {
        BIGSERIAL id PK
        BIGINT department_id FK
        BIGINT position_id FK
        BIGINT manager_id FK
        VARCHAR employee_number UK
        VARCHAR first_name
        VARCHAR last_name
        VARCHAR email
        DATE hire_date
        DATE termination_date
        VARCHAR status
        TIMESTAMPTZ created_at
    }

    positions {
        BIGSERIAL id PK
        BIGINT department_id FK
        VARCHAR title
        VARCHAR level
        DECIMAL min_salary
        DECIMAL max_salary
    }

    employee_contracts {
        BIGSERIAL id PK
        BIGINT employee_id FK
        VARCHAR contract_type
        DATE start_date
        DATE end_date
        DECIMAL salary
        VARCHAR pay_frequency
    }

    attendance {
        BIGSERIAL id PK
        BIGINT employee_id FK
        DATE work_date
        TIMESTAMPTZ check_in
        TIMESTAMPTZ check_out
        DECIMAL hours_worked
        VARCHAR status
    }

    leave_types {
        BIGSERIAL id PK
        VARCHAR name
        INTEGER default_days
        BOOLEAN is_paid
    }

    leave_requests {
        BIGSERIAL id PK
        BIGINT employee_id FK
        BIGINT leave_type_id FK
        BIGINT approved_by FK
        DATE start_date
        DATE end_date
        VARCHAR status
        TEXT reason
    }

    payroll_periods {
        BIGSERIAL id PK
        BIGINT company_id FK
        DATE start_date
        DATE end_date
        VARCHAR status
    }

    payroll_items {
        BIGSERIAL id PK
        BIGINT payroll_period_id FK
        BIGINT employee_id FK
        DECIMAL base_salary
        DECIMAL overtime
        DECIMAL deductions
        DECIMAL net_pay
    }

    benefits {
        BIGSERIAL id PK
        VARCHAR name
        VARCHAR type
        DECIMAL company_contribution
        DECIMAL employee_contribution
    }

    employee_benefits {
        BIGSERIAL id PK
        BIGINT employee_id FK
        BIGINT benefit_id FK
        DATE enrollment_date
        VARCHAR status
    }

    %% ==========================================
    %% FINANCE & ACCOUNTING
    %% ==========================================

    chart_of_accounts {
        BIGSERIAL id PK
        BIGINT company_id FK
        BIGINT parent_id FK
        VARCHAR account_number UK
        VARCHAR name
        VARCHAR account_type
        BOOLEAN is_active
    }

    journal_entries {
        BIGSERIAL id PK
        BIGINT fiscal_year_id FK
        BIGINT created_by FK
        VARCHAR entry_number UK
        DATE entry_date
        TEXT description
        VARCHAR status
        TIMESTAMPTZ posted_at
    }

    journal_lines {
        BIGSERIAL id PK
        BIGINT journal_entry_id FK
        BIGINT account_id FK
        DECIMAL debit
        DECIMAL credit
        TEXT memo
    }

    budgets {
        BIGSERIAL id PK
        BIGINT fiscal_year_id FK
        BIGINT department_id FK
        VARCHAR name
        DECIMAL total_amount
        VARCHAR status
    }

    budget_lines {
        BIGSERIAL id PK
        BIGINT budget_id FK
        BIGINT account_id FK
        DECIMAL amount
        VARCHAR period
    }

    tax_rates {
        BIGSERIAL id PK
        BIGINT company_id FK
        VARCHAR name
        VARCHAR code
        DECIMAL rate
        BOOLEAN is_active
    }

    bank_accounts {
        BIGSERIAL id PK
        BIGINT company_id FK
        BIGINT account_id FK
        VARCHAR bank_name
        VARCHAR account_number
        VARCHAR routing_number
        DECIMAL balance
    }

    bank_transactions {
        BIGSERIAL id PK
        BIGINT bank_account_id FK
        BIGINT journal_entry_id FK
        DATE transaction_date
        VARCHAR type
        DECIMAL amount
        VARCHAR reference
        BOOLEAN is_reconciled
    }

    %% ==========================================
    %% INVENTORY & WAREHOUSE
    %% ==========================================

    warehouses {
        BIGSERIAL id PK
        BIGINT branch_id FK
        VARCHAR name
        VARCHAR code UK
        VARCHAR address
        BOOLEAN is_active
    }

    warehouse_zones {
        BIGSERIAL id PK
        BIGINT warehouse_id FK
        VARCHAR name
        VARCHAR zone_type
    }

    storage_locations {
        BIGSERIAL id PK
        BIGINT zone_id FK
        VARCHAR aisle
        VARCHAR rack
        VARCHAR shelf
        VARCHAR bin
        INTEGER capacity
    }

    product_categories {
        BIGSERIAL id PK
        BIGINT parent_id FK
        VARCHAR name
        VARCHAR code
        TEXT description
    }

    products {
        BIGSERIAL id PK
        BIGINT category_id FK
        BIGINT default_supplier_id FK
        VARCHAR sku UK
        VARCHAR name
        TEXT description
        VARCHAR unit_of_measure
        DECIMAL weight
        DECIMAL cost_price
        DECIMAL selling_price
        INTEGER reorder_point
        INTEGER reorder_quantity
        BOOLEAN is_active
    }

    product_variants {
        BIGSERIAL id PK
        BIGINT product_id FK
        VARCHAR sku UK
        VARCHAR variant_name
        DECIMAL price_adjustment
        BOOLEAN is_active
    }

    inventory_stock {
        BIGSERIAL id PK
        BIGINT product_id FK
        BIGINT warehouse_id FK
        BIGINT location_id FK
        INTEGER quantity_on_hand
        INTEGER quantity_reserved
        INTEGER quantity_available
        TIMESTAMPTZ last_counted
    }

    inventory_movements {
        BIGSERIAL id PK
        BIGINT product_id FK
        BIGINT from_warehouse_id FK
        BIGINT to_warehouse_id FK
        BIGINT created_by FK
        VARCHAR movement_type
        INTEGER quantity
        VARCHAR reference_type
        BIGINT reference_id
        TIMESTAMPTZ created_at
    }

    inventory_adjustments {
        BIGSERIAL id PK
        BIGINT warehouse_id FK
        BIGINT approved_by FK
        VARCHAR adjustment_number UK
        DATE adjustment_date
        TEXT reason
        VARCHAR status
    }

    adjustment_lines {
        BIGSERIAL id PK
        BIGINT adjustment_id FK
        BIGINT product_id FK
        INTEGER quantity_before
        INTEGER quantity_after
        INTEGER variance
    }

    %% ==========================================
    %% SUPPLIERS & PURCHASING
    %% ==========================================

    suppliers {
        BIGSERIAL id PK
        VARCHAR code UK
        VARCHAR name
        VARCHAR contact_name
        VARCHAR email
        VARCHAR phone
        VARCHAR address
        VARCHAR payment_terms
        BOOLEAN is_active
    }

    supplier_products {
        BIGSERIAL id PK
        BIGINT supplier_id FK
        BIGINT product_id FK
        VARCHAR supplier_sku
        DECIMAL cost
        INTEGER lead_time_days
        INTEGER min_order_qty
    }

    purchase_orders {
        BIGSERIAL id PK
        BIGINT supplier_id FK
        BIGINT warehouse_id FK
        BIGINT created_by FK
        BIGINT approved_by FK
        VARCHAR po_number UK
        DATE order_date
        DATE expected_date
        VARCHAR status
        DECIMAL subtotal
        DECIMAL tax
        DECIMAL total
    }

    purchase_order_lines {
        BIGSERIAL id PK
        BIGINT purchase_order_id FK
        BIGINT product_id FK
        INTEGER quantity_ordered
        INTEGER quantity_received
        DECIMAL unit_cost
        DECIMAL line_total
    }

    goods_receipts {
        BIGSERIAL id PK
        BIGINT purchase_order_id FK
        BIGINT warehouse_id FK
        BIGINT received_by FK
        VARCHAR receipt_number UK
        DATE receipt_date
        TEXT notes
    }

    goods_receipt_lines {
        BIGSERIAL id PK
        BIGINT goods_receipt_id FK
        BIGINT po_line_id FK
        BIGINT location_id FK
        INTEGER quantity_received
        VARCHAR condition
    }

    %% ==========================================
    %% CUSTOMERS & SALES
    %% ==========================================

    customers {
        BIGSERIAL id PK
        VARCHAR code UK
        VARCHAR name
        VARCHAR contact_name
        VARCHAR email
        VARCHAR phone
        VARCHAR billing_address
        VARCHAR shipping_address
        VARCHAR payment_terms
        DECIMAL credit_limit
        BOOLEAN is_active
    }

    customer_contacts {
        BIGSERIAL id PK
        BIGINT customer_id FK
        VARCHAR name
        VARCHAR title
        VARCHAR email
        VARCHAR phone
        BOOLEAN is_primary
    }

    price_lists {
        BIGSERIAL id PK
        VARCHAR name
        BIGINT currency_id FK
        DATE effective_from
        DATE effective_to
        BOOLEAN is_default
    }

    price_list_items {
        BIGSERIAL id PK
        BIGINT price_list_id FK
        BIGINT product_id FK
        DECIMAL price
        DECIMAL min_quantity
    }

    sales_quotes {
        BIGSERIAL id PK
        BIGINT customer_id FK
        BIGINT sales_rep_id FK
        VARCHAR quote_number UK
        DATE quote_date
        DATE valid_until
        VARCHAR status
        DECIMAL subtotal
        DECIMAL discount
        DECIMAL tax
        DECIMAL total
    }

    sales_quote_lines {
        BIGSERIAL id PK
        BIGINT quote_id FK
        BIGINT product_id FK
        INTEGER quantity
        DECIMAL unit_price
        DECIMAL discount
        DECIMAL line_total
    }

    sales_orders {
        BIGSERIAL id PK
        BIGINT customer_id FK
        BIGINT quote_id FK
        BIGINT sales_rep_id FK
        BIGINT warehouse_id FK
        VARCHAR order_number UK
        DATE order_date
        DATE requested_date
        VARCHAR status
        DECIMAL subtotal
        DECIMAL discount
        DECIMAL tax
        DECIMAL shipping
        DECIMAL total
    }

    sales_order_lines {
        BIGSERIAL id PK
        BIGINT sales_order_id FK
        BIGINT product_id FK
        INTEGER quantity_ordered
        INTEGER quantity_shipped
        DECIMAL unit_price
        DECIMAL discount
        DECIMAL line_total
    }

    shipments {
        BIGSERIAL id PK
        BIGINT sales_order_id FK
        BIGINT warehouse_id FK
        BIGINT shipped_by FK
        VARCHAR shipment_number UK
        DATE ship_date
        VARCHAR carrier
        VARCHAR tracking_number
        VARCHAR status
    }

    shipment_lines {
        BIGSERIAL id PK
        BIGINT shipment_id FK
        BIGINT order_line_id FK
        INTEGER quantity_shipped
    }

    sales_returns {
        BIGSERIAL id PK
        BIGINT sales_order_id FK
        BIGINT customer_id FK
        VARCHAR return_number UK
        DATE return_date
        TEXT reason
        VARCHAR status
        DECIMAL refund_amount
    }

    %% ==========================================
    %% INVOICING & PAYMENTS
    %% ==========================================

    invoices {
        BIGSERIAL id PK
        BIGINT customer_id FK
        BIGINT sales_order_id FK
        BIGINT journal_entry_id FK
        VARCHAR invoice_number UK
        DATE invoice_date
        DATE due_date
        VARCHAR status
        DECIMAL subtotal
        DECIMAL tax
        DECIMAL total
        DECIMAL amount_paid
    }

    invoice_lines {
        BIGSERIAL id PK
        BIGINT invoice_id FK
        BIGINT product_id FK
        VARCHAR description
        INTEGER quantity
        DECIMAL unit_price
        DECIMAL tax_amount
        DECIMAL line_total
    }

    payments {
        BIGSERIAL id PK
        BIGINT customer_id FK
        BIGINT bank_account_id FK
        BIGINT journal_entry_id FK
        VARCHAR payment_number UK
        DATE payment_date
        VARCHAR payment_method
        DECIMAL amount
        VARCHAR reference
    }

    payment_allocations {
        BIGSERIAL id PK
        BIGINT payment_id FK
        BIGINT invoice_id FK
        DECIMAL amount_applied
    }

    %% ==========================================
    %% RELATIONSHIPS
    %% ==========================================

    %% Company Structure
    companies ||--o{ branches : "has"
    companies ||--o{ fiscal_years : "has"
    companies ||--o{ chart_of_accounts : "has"
    companies ||--o{ tax_rates : "has"
    companies ||--o{ bank_accounts : "has"
    companies ||--o{ payroll_periods : "has"
    branches ||--o{ departments : "contains"
    branches ||--o{ warehouses : "operates"
    departments ||--o{ departments : "parent of"
    departments ||--o{ employees : "employs"
    departments ||--o{ positions : "has"
    departments ||--o{ budgets : "has"

    %% HR
    employees ||--o{ employees : "manages"
    employees ||--o{ employee_contracts : "has"
    employees ||--o{ attendance : "records"
    employees ||--o{ leave_requests : "submits"
    employees ||--o{ payroll_items : "receives"
    employees ||--o{ employee_benefits : "enrolled in"
    positions ||--o{ employees : "filled by"
    leave_types ||--o{ leave_requests : "categorizes"
    payroll_periods ||--o{ payroll_items : "contains"
    benefits ||--o{ employee_benefits : "provides"

    %% Finance
    fiscal_years ||--o{ journal_entries : "contains"
    fiscal_years ||--o{ budgets : "plans"
    chart_of_accounts ||--o{ chart_of_accounts : "parent of"
    chart_of_accounts ||--o{ journal_lines : "debited/credited"
    chart_of_accounts ||--o{ budget_lines : "budgeted"
    chart_of_accounts ||--o{ bank_accounts : "linked to"
    journal_entries ||--o{ journal_lines : "contains"
    journal_entries ||--o{ bank_transactions : "reconciles"
    budgets ||--o{ budget_lines : "contains"
    bank_accounts ||--o{ bank_transactions : "has"
    bank_accounts ||--o{ payments : "receives"

    %% Inventory
    warehouses ||--o{ warehouse_zones : "divided into"
    warehouses ||--o{ inventory_stock : "stores"
    warehouses ||--o{ inventory_adjustments : "adjusted"
    warehouses ||--o{ purchase_orders : "receives"
    warehouses ||--o{ sales_orders : "ships from"
    warehouse_zones ||--o{ storage_locations : "contains"
    storage_locations ||--o{ inventory_stock : "holds"
    storage_locations ||--o{ goods_receipt_lines : "stores"
    product_categories ||--o{ product_categories : "parent of"
    product_categories ||--o{ products : "categorizes"
    products ||--o{ product_variants : "has"
    products ||--o{ inventory_stock : "stocked"
    products ||--o{ inventory_movements : "moved"
    products ||--o{ supplier_products : "supplied by"
    products ||--o{ purchase_order_lines : "ordered"
    products ||--o{ price_list_items : "priced in"
    products ||--o{ sales_quote_lines : "quoted"
    products ||--o{ sales_order_lines : "sold"
    products ||--o{ invoice_lines : "invoiced"
    products ||--o{ adjustment_lines : "adjusted"
    inventory_adjustments ||--o{ adjustment_lines : "contains"

    %% Suppliers & Purchasing
    suppliers ||--o{ supplier_products : "provides"
    suppliers ||--o{ purchase_orders : "receives"
    purchase_orders ||--o{ purchase_order_lines : "contains"
    purchase_orders ||--o{ goods_receipts : "received via"
    goods_receipts ||--o{ goods_receipt_lines : "contains"

    %% Customers & Sales
    customers ||--o{ customer_contacts : "has"
    customers ||--o{ sales_quotes : "receives"
    customers ||--o{ sales_orders : "places"
    customers ||--o{ invoices : "billed"
    customers ||--o{ payments : "pays"
    customers ||--o{ sales_returns : "returns"
    currencies ||--o{ price_lists : "denominates"
    price_lists ||--o{ price_list_items : "contains"
    sales_quotes ||--o{ sales_quote_lines : "contains"
    sales_quotes ||--o{ sales_orders : "converts to"
    sales_orders ||--o{ sales_order_lines : "contains"
    sales_orders ||--o{ shipments : "fulfilled by"
    sales_orders ||--o{ invoices : "billed via"
    sales_orders ||--o{ sales_returns : "returned"
    shipments ||--o{ shipment_lines : "contains"
    invoices ||--o{ invoice_lines : "contains"
    invoices ||--o{ payment_allocations : "paid by"
    payments ||--o{ payment_allocations : "applied to"
```

## Module Summary

| Module         | Tables | Description                                                           |
| -------------- | ------ | --------------------------------------------------------------------- |
| **Core**       | 5      | Company structure, branches, departments, currencies, fiscal years    |
| **HR**         | 10     | Employees, positions, contracts, attendance, leave, payroll, benefits |
| **Finance**    | 8      | Chart of accounts, journals, budgets, taxes, bank accounts            |
| **Inventory**  | 10     | Warehouses, zones, locations, products, stock, movements, adjustments |
| **Purchasing** | 6      | Suppliers, purchase orders, goods receipts                            |
| **Sales**      | 10     | Customers, quotes, orders, shipments, returns                         |
| **Invoicing**  | 4      | Invoices, payments, allocations                                       |

**Total: 53 tables**
