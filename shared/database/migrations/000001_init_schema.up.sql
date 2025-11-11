CREATE TABLE md_table_statuses (
                                   id BIGINT NOT NULL PRIMARY KEY
    ,code VARCHAR(15) NOT NULL UNIQUE
    ,NAME VARCHAR(100) NOT NULL UNIQUE
    ,name_en VARCHAR(100) NOT NULL UNIQUE
    ,created_at TIMESTAMP WITH TIME zone DEFAULT now() NOT NULL
    ,updated_at TIMESTAMP WITH TIME zone
);

CREATE TABLE md_session_extension_mode (
                                           id BIGINT NOT NULL PRIMARY KEY
    ,code VARCHAR(50) NOT NULL UNIQUE
    ,NAME VARCHAR(100) NOT NULL UNIQUE
    ,name_en VARCHAR(100) NOT NULL UNIQUE
    ,sort_order INT DEFAULT 1 NOT NULL UNIQUE
    ,created_at TIMESTAMP WITH TIME zone DEFAULT now() NOT NULL
    ,updated_at TIMESTAMP WITH TIME zone
);

CREATE TABLE md_session_extension_reason (
                                             id BIGINT NOT NULL PRIMARY KEY
    ,code VARCHAR(50) NOT NULL UNIQUE
    ,NAME VARCHAR(100) NOT NULL UNIQUE
    ,name_en VARCHAR(100) NOT NULL UNIQUE
    ,category VARCHAR(50)
    ,mode_code VARCHAR(50)
    ,is_active boolean DEFAULT true NOT NULL
    ,sort_order INT DEFAULT 1 NOT NULL
    ,created_at TIMESTAMP WITH TIME zone DEFAULT now() NOT NULL
    ,updated_at TIMESTAMP WITH TIME zone
);

ALTER TABLE md_session_extension_reason OWNER TO postgres;

CREATE TABLE md_categories (
                               id BIGINT NOT NULL PRIMARY KEY
    ,NAME VARCHAR(100) NOT NULL UNIQUE
    ,name_en VARCHAR(100) NOT NULL UNIQUE
    ,created_at TIMESTAMP WITH TIME zone DEFAULT now() NOT NULL
    ,updated_at TIMESTAMP WITH TIME zone
    ,icon_name VARCHAR(50)
    ,sort_order INT DEFAULT 1 NOT NULL UNIQUE
    ,is_visible boolean DEFAULT true NOT NULL
    ,code VARCHAR(50) NOT NULL UNIQUE
);

CREATE TABLE md_order_statuses (
                                   id BIGINT NOT NULL PRIMARY KEY
    ,code VARCHAR(15) NOT NULL UNIQUE
    ,NAME VARCHAR(100) NOT NULL UNIQUE
    ,name_en VARCHAR(100) NOT NULL UNIQUE
    ,sort_order INT NOT NULL UNIQUE
    ,is_final boolean DEFAULT false NOT NULL
    ,created_at TIMESTAMP WITH TIME zone DEFAULT now() NOT NULL
    ,updated_at TIMESTAMP WITH TIME zone
);

CREATE TABLE md_payment_methods (
                                    id BIGINT NOT NULL PRIMARY KEY
    ,code VARCHAR(15) NOT NULL UNIQUE
    ,NAME VARCHAR(100) NOT NULL UNIQUE
    ,name_en VARCHAR(100) NOT NULL UNIQUE
    ,enable boolean DEFAULT false NOT NULL
    ,created_at TIMESTAMP WITH TIME zone DEFAULT now() NOT NULL
    ,updated_at TIMESTAMP WITH TIME zone
);

CREATE TABLE md_payment_statuses (
                                     id BIGINT NOT NULL PRIMARY KEY
    ,code VARCHAR(15) NOT NULL UNIQUE
    ,NAME VARCHAR(100) NOT NULL UNIQUE
    ,name_en VARCHAR(100) NOT NULL UNIQUE
    ,is_final boolean DEFAULT false NOT NULL
    ,created_at TIMESTAMP WITH TIME zone DEFAULT now() NOT NULL
    ,updated_at TIMESTAMP WITH TIME zone
);

CREATE TABLE session_extension
(
    id                BIGINT NOT NULL PRIMARY KEY,
    session_id        UUID NOT NULL REFERENCES table_session (session_id),
    requested_minutes INTEGER NOT NULL,
    created_at        TIMESTAMP WITH time zone DEFAULT Now() NOT NULL,
    mode_id           BIGINT CONSTRAINT
        session_extension_md_session_extension_mode_id_fk
        REFERENCES
            md_session_extension_mode,
    reason_id         BIGINT CONSTRAINT
        session_extension_md_session_extension_reason_id_fk
        REFERENCES
            md_session_extension_reason
);

CREATE TABLE tables (
                        id BIGINT NOT NULL PRIMARY KEY
    ,table_number INT NOT NULL UNIQUE
    ,status_id BIGINT NOT NULL REFERENCES md_table_statuses
    ,seats INT DEFAULT 0 NOT NULL
    ,created_at TIMESTAMP WITH TIME zone DEFAULT now() NOT NULL
    ,updated_at TIMESTAMP WITH TIME zone
);


CREATE TABLE table_session (
       id                   bigint                                             not null
           primary key,
       table_id             bigint                                             not null
           references tables,
       session_id           uuid                     default gen_random_uuid() not null
           unique,
       number_of_people     integer                  default 1                 not null,
       status               table_session_status,
       started_at           timestamp with time zone default now()             not null,
       expires_at           timestamp with time zone                           not null,
       ended_at           timestamp with time zone                           ,
       max_extend_minutes   integer                  default 120               not null,
       extend_count         integer                  default 0                 not null,
       extend_total_minutes integer                  default 0                 not null,
       last_reason_code     text,
       lock_version         integer                  default 1                 not null
    ,CONSTRAINT chk_time_order CHECK (
        (started_at <= expires_at)
        )
);

CREATE TABLE products (
                          id BIGINT NOT NULL PRIMARY KEY
    ,NAME VARCHAR(255) NOT NULL UNIQUE
    ,name_en VARCHAR(255) NOT NULL UNIQUE
    ,categories BIGINT NOT NULL REFERENCES md_categories
    ,description TEXT
    ,price NUMERIC(10, 2) DEFAULT 0 NOT NULL
    ,is_available boolean DEFAULT false NOT NULL
    ,is_visible boolean DEFAULT true NOT NULL
    ,image_url TEXT
    ,created_at TIMESTAMP WITH TIME zone DEFAULT now() NOT NULL
    ,updated_at TIMESTAMP WITH TIME zone
);

CREATE TABLE orders (
                        id BIGINT NOT NULL PRIMARY KEY
    ,order_number VARCHAR(50) NOT NULL UNIQUE
    ,session_id uuid
    ,table_id BIGINT NOT NULL REFERENCES tables
    ,status_id BIGINT NOT NULL REFERENCES md_order_statuses
    ,total_amount NUMERIC(10, 2) DEFAULT 0 NOT NULL
    ,created_at TIMESTAMP WITH TIME zone DEFAULT now() NOT NULL
    ,updated_at TIMESTAMP WITH TIME zone
);

CREATE TABLE order_items (
                             id BIGINT NOT NULL PRIMARY KEY
    ,order_id BIGINT NOT NULL REFERENCES orders
    ,product_id BIGINT NOT NULL REFERENCES products
    ,status_id BIGINT NOT NULL REFERENCES md_order_statuses
    ,product_name VARCHAR(255) NOT NULL
    ,product_name_en VARCHAR(255) NOT NULL
    ,price NUMERIC(10, 2) DEFAULT 0 NOT NULL
    ,quantity INT DEFAULT 1 NOT NULL
    ,note TEXT
    ,prepared_at TIMESTAMP WITH TIME zone
    ,created_at TIMESTAMP WITH TIME zone DEFAULT now() NOT NULL
    ,updated_at TIMESTAMP WITH TIME zone
    ,product_image_url TEXT
    ,is_visible boolean DEFAULT true NOT NULL
);

comment ON COLUMN order_items.prepared_at IS 'เวลาที่ทำอาหารเสร็จ';

CREATE TABLE payments (
                          id BIGINT NOT NULL PRIMARY KEY
    ,order_id BIGINT NOT NULL REFERENCES orders
    ,amount NUMERIC(10, 2) DEFAULT 0 NOT NULL
    ,method BIGINT NOT NULL REFERENCES md_payment_methods
    ,STATUS BIGINT NOT NULL REFERENCES md_payment_statuses
    ,paid_at TIMESTAMP WITH TIME zone
    ,transaction_id TEXT NOT NULL UNIQUE
    ,ref_code VARCHAR(150) NOT NULL UNIQUE
    ,note TEXT
    ,created_at TIMESTAMP WITH TIME zone DEFAULT now() NOT NULL
    ,updated_at TIMESTAMP WITH TIME zone
);

create table product_time_extension
(
    id               bigint                                 not null
        constraint product_time_extension_pk
            primary key,
    duration_minutes integer                  default 0     not null,
    products_id      bigint
        constraint product_time_extension_products_id_fk
            references products,
    created_at       timestamp with time zone default now() not null,
    updated_at       timestamp with time zone
);

ALTER TABLE md_table_statuses OWNER TO postgres;

CREATE INDEX md_table_statuses_id_idx ON md_table_statuses (id);

CREATE INDEX md_table_statuses_code_idx ON md_table_statuses (code);

ALTER TABLE md_session_extension_mode OWNER TO postgres;

CREATE INDEX md_session_extension_mode_id_idx ON md_session_extension_mode (id);

ALTER TABLE session_extension OWNER TO postgres;

CREATE INDEX session_extension_session_id_idx ON session_extension (session_id);

CREATE INDEX session_extension_created_at_idx ON session_extension (created_at);

CREATE INDEX md_session_extension_reason_id_idx ON md_session_extension_reason (id);

CREATE INDEX md_session_extension_reason_code_idx ON md_session_extension_reason (code);

ALTER TABLE md_categories OWNER TO postgres;

CREATE INDEX md_categories_id_idx ON md_categories (id);

ALTER TABLE md_order_statuses OWNER TO postgres;

CREATE INDEX md_order_statuses_id_idx ON md_order_statuses (id);

CREATE INDEX md_order_statuses_code_idx ON md_order_statuses (code);

ALTER TABLE md_payment_methods OWNER TO postgres;

CREATE INDEX md_payment_methods_id_idx ON md_payment_methods (id);

CREATE INDEX md_payment_methods_code_idx ON md_payment_methods (code);

ALTER TABLE md_payment_statuses OWNER TO postgres;

CREATE INDEX md_payment_statuses_id_idx ON md_payment_statuses (id);

CREATE INDEX md_payment_statuses_code_idx ON md_payment_statuses (code);

ALTER TABLE tables OWNER TO postgres;

CREATE INDEX tables_id_idx ON tables (id);

CREATE INDEX tables_table_number_idx ON tables (table_number);

CREATE INDEX tables_status_id_idx ON tables (status_id);

ALTER TABLE table_session OWNER TO postgres;

CREATE INDEX table_session_id_idx ON table_session (id);

CREATE INDEX table_session_table_id_idx ON table_session (table_id);

CREATE INDEX table_session_session_id_idx ON table_session (session_id);

CREATE INDEX idx_session_status_expires_at ON table_session (
                                                             STATUS
    ,expires_at
    );

ALTER TABLE products OWNER TO postgres;

CREATE INDEX products_id_idx ON products (id);

CREATE INDEX products_name_idx ON products (NAME);

CREATE INDEX products_name_en_idx ON products (name_en);

ALTER TABLE orders OWNER TO postgres;

CREATE INDEX orders_id_idx ON orders (id);

CREATE INDEX orders_order_number_idx ON orders (order_number);

CREATE INDEX orders_status_id_idx ON orders (status_id);

CREATE INDEX orders_created_at_idx ON orders (created_at);

ALTER TABLE order_items OWNER TO postgres;

CREATE INDEX order_items_id_idx ON order_items (id);

CREATE INDEX order_items_order_id ON order_items (
                                                  order_id
    ,id
    );

CREATE INDEX order_items_order_id_idx ON order_items (order_id);

CREATE TABLE order_sequences (
                                 order_date DATE NOT NULL PRIMARY KEY
    ,current_number INT NOT NULL
);

ALTER TABLE order_sequences OWNER TO postgres;

ALTER TABLE payments OWNER TO postgres;

CREATE INDEX payments_id_idx ON payments (id);

CREATE INDEX payments_order_id_idx ON payments (order_id);

CREATE INDEX payments_transaction_id_idx ON payments (transaction_id);

CREATE INDEX transaction_id_status ON payments (
                                                transaction_id
    ,STATUS
    );

create type public.table_session_status as enum ('active', 'closed', 'expired');

alter type public.table_session_status owner to postgres;



-- Insert Master Data

insert into public.md_categories (id, name, name_en, icon_name, sort_order, created_at, updated_at)
values  (1921143886227443712, 'อาหาร', 'Food', 'soup', 1, '2025-10-10 05:34:36.643495 +00:00', null),
        (1921144050476388352, 'เครื่องดื่ม', 'Drink', 'glass-water', 2, '2025-10-10 05:34:36.643495 +00:00', null),
        (1921144250070732800, 'ขนม', 'Dessert', 'cake-slice', 3, '2025-10-10 05:34:36.643495 +00:00', null);

insert into public.md_order_statuses (id, code, name, name_en, sort_order, is_final, created_at, updated_at)
values  (1921868485739155456, 'PENDING', 'รอยืนยันออเดอร์', 'Pending', 1, false, '2025-09-18 11:43:59.552632 +00:00', null),
        (1921868485739155457, 'CONFIRMED', 'ยืนยันออเดอร์', 'Confirmed', 2, false, '2025-09-18 11:43:59.552632 +00:00', null),
        (1921868485739155458, 'PREPARING', 'กำลังเตรียมอาหาร', 'Preparing', 3, false, '2025-09-18 11:43:59.552632 +00:00', null),
        (1921868485739155459, 'SERVED', 'เสิร์ฟอาหารแล้ว', 'Served', 4, false, '2025-09-18 11:43:59.552632 +00:00', null),
        (1921868485739155460, 'WAITING_PAYMENT', 'รอชำระเงิน', 'Waiting for Payment', 5, false, '2025-09-18 11:43:59.552632 +00:00', null),
        (1921868485739155461, 'COMPLETED', 'เสร็จสิ้น', 'Completed', 6, true, '2025-09-18 11:43:59.552632 +00:00', null),
        (1921868485739155462, 'CANCELLED', 'ยกเลิก', 'Cancelled', 7, true, '2025-09-18 11:43:59.552632 +00:00', null);

insert into public.md_table_statuses (id, code, name, name_en, created_at, updated_at)
values  (1919968486671519744, 'AVAILABLE', 'ว่าง', 'Available', '2025-09-18 11:43:59.552632 +00:00', null),
        (1919968486843486208, 'RESERVED', 'ถูกจองล่วงหน้า', 'Reserved', '2025-09-18 11:43:59.552632 +00:00', null),
        (1919968486847680514, 'WAITING_PAYMENT', 'รอชำระเงิน', 'Waiting for Payment', '2025-09-18 11:43:59.552632 +00:00', null),
        (1919968486847680515, 'CLEANING', 'รอทำความสะอาด', 'Cleaning', '2025-09-18 11:43:59.552632 +00:00', null),
        (1919968486847680516, 'DISABLED', 'ปิดการใช้งานชั่วคราว', 'Disabled', '2025-09-18 11:43:59.552632 +00:00', null),
        (1919968486847680512, 'WAIT_ORDER', 'รอสั่ง', 'Waiting to Order', '2025-09-18 11:43:59.552632 +00:00', null),
        (1919968486847680513, 'WAIT_SERVE', 'รอเสิร์ฟ', 'Waiting to be Served', '2025-09-18 11:43:59.552632 +00:00', null),
        (1919968486847680517, 'FOOD_SERVED', 'อาหารครบแล้ว', 'Food Served', '2025-09-18 16:10:29.071839 +00:00', null),
        (1919968486847680518, 'CALL_WAITER', 'เรียกพนักงาน', 'Call Waiter', '2025-09-18 16:10:29.071839 +00:00', null);


insert into public.products (id, name, name_en, categories, description, price, is_available, image_url, created_at, updated_at)
values  (1921822053405560832, 'ข้าวผัด', 'Fried rice', 1921143886227443712, 'lorem ipso', 60.00, true, 'https://images.unsplash.com/photo-1603133872878-684f208fb84b?ixlib=rb-4.1.0&ixid=M3wxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8fA%3D%3D&auto=format&fit=crop&q=80&w=500', '2025-10-10 05:34:36.692235 +00:00', null),
        (1921828287366041600, 'ข้าวผัดกระเพรา', 'Phat kaphrao', 1921143886227443712, 'lorem ipso', 70.50, true, 'https://images.unsplash.com/photo-1694499792070-48e64a00cf0a?ixlib=rb-4.1.0&ixid=M3wxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8fA%3D%3D&auto=format&fit=crop&q=80&w=500', '2025-10-10 05:34:36.692235 +00:00', null),
        (1921821817723424768, 'ข้าวมันไก่', 'Chicken rice', 1921143886227443712, 'lorem ipso', 80.00, true, 'https://images.unsplash.com/photo-1749640566096-5d8098d452b4?ixlib=rb-4.1.0&ixid=M3wxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8fA%3D%3D&auto=format&fit=crop&q=80&w=500', '2025-10-10 05:34:36.692235 +00:00', null),
        (1921822608437809152, 'แป๊ปซี่', 'Pepsi', 1921144050476388352, 'lorem ipso', 30.00, true, 'https://images.unsplash.com/photo-1651000877733-fe2c0a70b3cd?ixlib=rb-4.1.0&ixid=M3wxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8fA%3D%3D&auto=format&fit=crop&q=80&w=500', '2025-10-10 05:34:36.692235 +00:00', null),
        (1921822481287483392, 'เค้กแครอท', 'Carrot cake', 1921144250070732800, 'lorem ipso', 120.00, true, 'https://images.unsplash.com/photo-1622926421334-6829deee4b4b?ixlib=rb-4.1.0&ixid=M3wxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8fA%3D%3D&auto=format&fit=crop&q=80&w=500', '2025-10-10 05:34:36.692235 +00:00', null);

INSERT INTO public."tables" (id,table_number,status_id,seats,created_at,updated_at) VALUES
                                                                                        (1920153361642950656,5,1919968486671519744,4,NOW(),NULL),
                                                                                        (1919972141986484224,3,1919968486671519744,4,NOW(),NULL),
                                                                                        (1919996486741921792,4,1919968486671519744,5,NOW(),NULL),
                                                                                        (1919968785813475328,1,1919968486671519744,5,NOW(),NULL),
                                                                                        (1919971956241731584,2,1919968486671519744,3,NOW(),NULL);


INSERT INTO public.md_payment_methods (id,code,"name",name_en,"enable",created_at,updated_at) VALUES
                                                                                               (1923732004537372672,'CASH','เงินสด','Cash',true,NOW(),NULL),
                                                                                               (1923732004537372675,'PROMPTPAY','พร้อมเพย์','PromptPay',true,NOW(),NULL);
