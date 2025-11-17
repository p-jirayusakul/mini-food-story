CREATE TABLE public.md_table_statuses (
                                          id BIGINT NOT NULL PRIMARY KEY
    ,code VARCHAR(15) NOT NULL UNIQUE
    ,name VARCHAR(100) NOT NULL UNIQUE
    ,name_en VARCHAR(100) NOT NULL UNIQUE
    ,created_at TIMESTAMP WITH TIME zone DEFAULT now() NOT NULL
    ,updated_at TIMESTAMP WITH TIME zone
);

ALTER TABLE public.md_table_statuses OWNER TO postgres;

CREATE INDEX md_table_statuses_id_idx ON public.md_table_statuses (id);

CREATE INDEX md_table_statuses_code_idx ON public.md_table_statuses (code);

CREATE TABLE public.md_categories (
                                      id BIGINT NOT NULL PRIMARY KEY
    ,name VARCHAR(100) NOT NULL UNIQUE
    ,name_en VARCHAR(100) NOT NULL UNIQUE
    ,created_at TIMESTAMP WITH TIME zone DEFAULT now() NOT NULL
    ,updated_at TIMESTAMP WITH TIME zone
    ,icon_name VARCHAR(50)
    ,sort_order INTEGER DEFAULT 1 NOT NULL UNIQUE
    ,is_visible boolean DEFAULT true NOT NULL
    ,code VARCHAR(50) UNIQUE
);

ALTER TABLE public.md_categories OWNER TO postgres;

CREATE INDEX md_categories_id_idx ON public.md_categories (id);

CREATE TABLE public.md_order_statuses (
                                          id BIGINT NOT NULL PRIMARY KEY
    ,code VARCHAR(15) NOT NULL UNIQUE
    ,name VARCHAR(100) NOT NULL UNIQUE
    ,name_en VARCHAR(100) NOT NULL UNIQUE
    ,sort_order INTEGER NOT NULL UNIQUE
    ,is_final boolean DEFAULT false NOT NULL
    ,created_at TIMESTAMP WITH TIME zone DEFAULT now() NOT NULL
    ,updated_at TIMESTAMP WITH TIME zone
);

ALTER TABLE public.md_order_statuses OWNER TO postgres;

CREATE INDEX md_order_statuses_id_idx ON public.md_order_statuses (id);

CREATE INDEX md_order_statuses_code_idx ON public.md_order_statuses (code);

CREATE TABLE public.md_payment_methods (
                                           id BIGINT NOT NULL PRIMARY KEY
    ,code VARCHAR(15) NOT NULL UNIQUE
    ,name VARCHAR(100) NOT NULL UNIQUE
    ,name_en VARCHAR(100) NOT NULL UNIQUE
    ,enable boolean DEFAULT false NOT NULL
    ,created_at TIMESTAMP WITH TIME zone DEFAULT now() NOT NULL
    ,updated_at TIMESTAMP WITH TIME zone
);

ALTER TABLE public.md_payment_methods OWNER TO postgres;

CREATE INDEX md_payment_methods_id_idx ON public.md_payment_methods (id);

CREATE INDEX md_payment_methods_code_idx ON public.md_payment_methods (code);

CREATE TABLE public.md_payment_statuses (
                                            id BIGINT NOT NULL PRIMARY KEY
    ,code VARCHAR(15) NOT NULL UNIQUE
    ,name VARCHAR(100) NOT NULL UNIQUE
    ,name_en VARCHAR(100) NOT NULL UNIQUE
    ,is_final boolean DEFAULT false NOT NULL
    ,created_at TIMESTAMP WITH TIME zone DEFAULT now() NOT NULL
    ,updated_at TIMESTAMP WITH TIME zone
);

ALTER TABLE public.md_payment_statuses OWNER TO postgres;

CREATE INDEX md_payment_statuses_id_idx ON public.md_payment_statuses (id);

CREATE INDEX md_payment_statuses_code_idx ON public.md_payment_statuses (code);

CREATE TABLE public.tables (
                               id BIGINT NOT NULL PRIMARY KEY
    ,table_number INTEGER NOT NULL UNIQUE
    ,status_id BIGINT NOT NULL REFERENCES public.md_table_statuses
    ,seats INTEGER DEFAULT 0 NOT NULL
    ,created_at TIMESTAMP WITH TIME zone DEFAULT now() NOT NULL
    ,updated_at TIMESTAMP WITH TIME zone
);

ALTER TABLE public.tables OWNER TO postgres;

CREATE INDEX tables_id_idx ON public.tables (id);

CREATE INDEX tables_table_number_idx ON public.tables (table_number);

CREATE INDEX tables_status_id_idx ON public.tables (status_id);

create type public.table_session_status as enum ('active', 'closed', 'expired');

alter type public.table_session_status owner to postgres;

CREATE TABLE public.table_session (
                                      id BIGINT NOT NULL PRIMARY KEY
    ,table_id BIGINT NOT NULL REFERENCES public.tables
    ,session_id uuid DEFAULT gen_random_uuid() NOT NULL UNIQUE
    ,number_of_people INTEGER DEFAULT 1 NOT NULL
    ,STATUS table_session_status
    ,started_at TIMESTAMP WITH TIME zone DEFAULT now() NOT NULL
    ,expires_at TIMESTAMP WITH TIME zone NOT NULL
    ,max_extend_minutes INTEGER DEFAULT 120 NOT NULL
    ,extend_count INTEGER DEFAULT 0 NOT NULL
    ,extend_total_minutes INTEGER DEFAULT 0 NOT NULL
    ,last_reason_code TEXT
    ,lock_version INTEGER DEFAULT 1 NOT NULL
    ,ended_at TIMESTAMP WITH TIME zone
);

ALTER TABLE public.table_session OWNER TO postgres;

CREATE INDEX table_session_id_idx ON public.table_session (id);

CREATE INDEX table_session_table_id_idx ON public.table_session (table_id);

CREATE INDEX table_session_session_id_idx ON public.table_session (session_id);

CREATE INDEX idx_session_status_expires_at ON public.table_session (
                                                                    STATUS
    ,expires_at
    );

CREATE TABLE public.products (
                                 id BIGINT NOT NULL PRIMARY KEY
    ,name VARCHAR(255) NOT NULL UNIQUE
    ,name_en VARCHAR(255) NOT NULL UNIQUE
    ,categories BIGINT NOT NULL REFERENCES public.md_categories
    ,description TEXT
    ,price NUMERIC(10, 2) DEFAULT 0 NOT NULL
    ,is_available boolean DEFAULT false NOT NULL
    ,image_url TEXT
    ,created_at TIMESTAMP WITH TIME zone DEFAULT now() NOT NULL
    ,updated_at TIMESTAMP WITH TIME zone
    ,is_visible boolean DEFAULT true NOT NULL
);

ALTER TABLE public.products OWNER TO postgres;

CREATE INDEX products_id_idx ON public.products (id);

CREATE INDEX products_name_idx ON public.products (name);

CREATE INDEX products_name_en_idx ON public.products (name_en);

CREATE TABLE public.orders (
                               id BIGINT NOT NULL PRIMARY KEY
    ,order_number VARCHAR(50) NOT NULL UNIQUE
    ,session_id uuid
    ,table_id BIGINT NOT NULL REFERENCES public.tables
    ,status_id BIGINT NOT NULL REFERENCES public.md_order_statuses
    ,total_amount NUMERIC(10, 2) DEFAULT 0 NOT NULL
    ,created_at TIMESTAMP WITH TIME zone DEFAULT now() NOT NULL
    ,updated_at TIMESTAMP WITH TIME zone
);

ALTER TABLE public.orders OWNER TO postgres;

CREATE INDEX orders_id_idx ON public.orders (id);

CREATE INDEX orders_order_number_idx ON public.orders (order_number);

CREATE INDEX orders_status_id_idx ON public.orders (status_id);

CREATE INDEX orders_created_at_idx ON public.orders (created_at);

CREATE TABLE public.order_items (
                                    id BIGINT NOT NULL PRIMARY KEY
    ,order_id BIGINT NOT NULL REFERENCES public.orders
    ,product_id BIGINT NOT NULL REFERENCES public.products
    ,status_id BIGINT NOT NULL REFERENCES public.md_order_statuses
    ,product_name VARCHAR(255) NOT NULL
    ,product_name_en VARCHAR(255) NOT NULL
    ,price NUMERIC(10, 2) DEFAULT 0 NOT NULL
    ,quantity INTEGER DEFAULT 1 NOT NULL
    ,note TEXT
    ,prepared_at TIMESTAMP WITH TIME zone
    ,created_at TIMESTAMP WITH TIME zone DEFAULT now() NOT NULL
    ,updated_at TIMESTAMP WITH TIME zone
    ,product_image_url TEXT
    ,is_visible boolean DEFAULT true NOT NULL
);

comment ON COLUMN public.order_items.prepared_at IS 'เวลาที่ทำอาหารเสร็จ';

ALTER TABLE public.order_items OWNER TO postgres;

CREATE INDEX order_items_id_idx ON public.order_items (id);

CREATE INDEX order_items_order_id ON public.order_items (
                                                         order_id
    ,id
    );

CREATE INDEX order_items_order_id_idx ON public.order_items (order_id);

CREATE TABLE public.order_sequences (
                                        order_date DATE NOT NULL PRIMARY KEY
    ,current_number INTEGER NOT NULL
);

ALTER TABLE public.order_sequences OWNER TO postgres;

CREATE TABLE public.payments (
                                 id BIGINT NOT NULL PRIMARY KEY
    ,order_id BIGINT NOT NULL REFERENCES public.orders
    ,amount NUMERIC(10, 2) DEFAULT 0 NOT NULL
    ,method BIGINT NOT NULL REFERENCES public.md_payment_methods
    ,STATUS BIGINT NOT NULL REFERENCES public.md_payment_statuses
    ,paid_at TIMESTAMP WITH TIME zone
    ,transaction_id TEXT NOT NULL UNIQUE
    ,ref_code VARCHAR(150) NOT NULL UNIQUE
    ,note TEXT
    ,created_at TIMESTAMP WITH TIME zone DEFAULT now() NOT NULL
    ,updated_at TIMESTAMP WITH TIME zone
);

ALTER TABLE public.payments OWNER TO postgres;

CREATE INDEX payments_id_idx ON public.payments (id);

CREATE INDEX payments_order_id_idx ON public.payments (order_id);

CREATE INDEX payments_transaction_id_idx ON public.payments (transaction_id);

CREATE INDEX transaction_id_status ON public.payments (
                                                       transaction_id
    ,STATUS
    );

CREATE TABLE public.md_session_extension_mode (
                                                  id BIGINT NOT NULL PRIMARY KEY
    ,code VARCHAR(50) NOT NULL UNIQUE
    ,name VARCHAR(100) NOT NULL UNIQUE
    ,name_en VARCHAR(100) NOT NULL UNIQUE
    ,sort_order INTEGER DEFAULT 1 NOT NULL UNIQUE
    ,created_at TIMESTAMP WITH TIME zone DEFAULT now() NOT NULL
    ,updated_at TIMESTAMP WITH TIME zone
);

ALTER TABLE public.md_session_extension_mode OWNER TO postgres;

CREATE INDEX md_session_extension_mode_id_idx ON public.md_session_extension_mode (id);

CREATE TABLE public.md_session_extension_reason (
                                                    id BIGINT NOT NULL PRIMARY KEY
    ,code VARCHAR(50) NOT NULL UNIQUE
    ,name VARCHAR(100) NOT NULL UNIQUE
    ,name_en VARCHAR(100) NOT NULL UNIQUE
    ,category VARCHAR(50)
    ,mode_code VARCHAR(50)
    ,is_active boolean DEFAULT true NOT NULL
    ,sort_order INTEGER DEFAULT 1 NOT NULL
    ,created_at TIMESTAMP WITH TIME zone DEFAULT now() NOT NULL
    ,updated_at TIMESTAMP WITH TIME zone
);

ALTER TABLE public.md_session_extension_reason OWNER TO postgres;

CREATE TABLE public.session_extension (
                                          id BIGINT NOT NULL PRIMARY KEY
    ,session_id uuid NOT NULL REFERENCES public.table_session(session_id)
    ,requested_minutes INTEGER NOT NULL
    ,created_at TIMESTAMP WITH TIME zone DEFAULT now() NOT NULL
    ,mode_id BIGINT CONSTRAINT session_extension_md_session_extension_mode_id_fk REFERENCES public.md_session_extension_mode
    ,reason_id BIGINT CONSTRAINT session_extension_md_session_extension_reason_id_fk REFERENCES public.md_session_extension_reason
);

ALTER TABLE public.session_extension OWNER TO postgres;

CREATE INDEX session_extension_session_id_idx ON public.session_extension (session_id);

CREATE INDEX session_extension_created_at_idx ON public.session_extension (created_at);

CREATE INDEX md_session_extension_reason_id_idx ON public.md_session_extension_reason (id);

CREATE INDEX md_session_extension_reason_code_idx ON public.md_session_extension_reason (code);

CREATE TABLE public.product_time_extension (
                                               id BIGINT NOT NULL CONSTRAINT product_time_extension_pk PRIMARY KEY
    ,duration_minutes INTEGER DEFAULT 0 NOT NULL
    ,products_id BIGINT CONSTRAINT product_time_extension_products_id_fk REFERENCES public.products
    ,created_at TIMESTAMP WITH TIME zone DEFAULT now() NOT NULL
    ,updated_at TIMESTAMP WITH TIME zone
);

ALTER TABLE public.product_time_extension OWNER TO postgres;

insert into public.md_categories (id, name, name_en, created_at, updated_at, icon_name, sort_order, is_visible, code)
values  (1921143886227443712, 'อาหาร', 'Food', '2025-10-10 05:34:36.643495 +00:00', null, 'soup', 1, true, 'FOOD'),
        (1921144050476388352, 'เครื่องดื่ม', 'Drink', '2025-10-10 05:34:36.643495 +00:00', null, 'glass-water', 2, true, 'DRINK'),
        (1921144250070732800, 'ขนม', 'Dessert', '2025-10-10 05:34:36.643495 +00:00', null, 'cake-slice', 3, true, 'DESSERT'),
        (1921144250070732801, 'ต่อเวลา', 'Time extension', '2025-10-10 05:34:36.643495 +00:00', null, 'timer', 4, false, 'TIME_EXTENSION');

insert into public.md_order_statuses (id, code, name, name_en, sort_order, is_final, created_at, updated_at)
values  (1921868485739155456, 'PENDING', 'รอยืนยันออเดอร์', 'Pending', 1, false, '2025-09-18 11:43:59.552632 +00:00', null),
        (1921868485739155457, 'CONFIRMED', 'ยืนยันออเดอร์', 'Confirmed', 2, false, '2025-09-18 11:43:59.552632 +00:00', null),
        (1921868485739155458, 'PREPARING', 'กำลังเตรียมอาหาร', 'Preparing', 3, false, '2025-09-18 11:43:59.552632 +00:00', null),
        (1921868485739155459, 'SERVED', 'เสิร์ฟอาหารแล้ว', 'Served', 4, false, '2025-09-18 11:43:59.552632 +00:00', null),
        (1921868485739155460, 'WAITING_PAYMENT', 'รอชำระเงิน', 'Waiting for Payment', 5, false, '2025-09-18 11:43:59.552632 +00:00', null),
        (1921868485739155461, 'COMPLETED', 'เสร็จสิ้น', 'Completed', 6, true, '2025-09-18 11:43:59.552632 +00:00', null),
        (1921868485739155462, 'CANCELLED', 'ยกเลิก', 'Cancelled', 7, true, '2025-09-18 11:43:59.552632 +00:00', null);

insert into public.md_payment_methods (id, code, name, name_en, enable, created_at, updated_at)
values  (1923732004537372672, 'CASH', 'เงินสด', 'Cash', true, '2025-11-14 08:41:17.586336 +00:00', null),
        (1923732004537372675, 'PROMPTPAY', 'พร้อมเพย์', 'PromptPay', true, '2025-11-14 08:41:17.586336 +00:00', null);

insert into public.md_session_extension_mode (id, code, name, name_en, sort_order, created_at, updated_at)
values  (1989547119583523322, 'COMP_FREE', 'ขยายเวลา (ฟรี)', 'Complimentary Free Extension', 1, '2025-11-14 08:59:08.870989 +00:00', null),
        (1989547119580791291, 'PAID', 'ขยายเวลา (เสียเงิน)', 'Paid Extension', 2, '2025-11-14 08:59:08.870989 +00:00', null);

insert into public.md_session_extension_reason (id, code, name, name_en, category, mode_code, is_active, sort_order, created_at, updated_at)
values  (1989546983509878172, 'DELAYED_SERVING', 'เสิร์ฟอาหารล่าช้า', 'Food Serving Delay', 'SERVICE', 'COMP_FREE', true, 2, '2025-11-14 09:01:39.538663 +00:00', null),
        (1989546983510861211, 'CUSTOMER_REQUEST', 'ลูกค้าขอขยายเวลาเอง', 'Customer Requested Extension', 'CUSTOMER', 'PAID', true, 1, '2025-11-14 09:01:39.538663 +00:00', null),
        (1989546983509939613, 'SYSTEM_ISSUE', 'ระบบมีปัญหา', 'System Issue', 'SYSTEM', 'COMP_FREE', true, 3, '2025-11-14 09:01:39.538663 +00:00', null);

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

insert into public.products (id, name, name_en, categories, description, price, is_available, image_url, created_at, updated_at, is_visible)
values  (1921822053405560832, 'ข้าวผัด', 'Fried rice', 1921143886227443712, 'lorem ipso', 60.00, true, 'https://images.unsplash.com/photo-1603133872878-684f208fb84b?ixlib=rb-4.1.0&ixid=M3wxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8fA%3D%3D&auto=format&fit=crop&q=80&w=500', '2025-10-10 05:34:36.692235 +00:00', null, true),
        (1921828287366041600, 'ข้าวผัดกระเพรา', 'Phat kaphrao', 1921143886227443712, 'lorem ipso', 70.50, true, 'https://images.unsplash.com/photo-1694499792070-48e64a00cf0a?ixlib=rb-4.1.0&ixid=M3wxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8fA%3D%3D&auto=format&fit=crop&q=80&w=500', '2025-10-10 05:34:36.692235 +00:00', null, true),
        (1921821817723424768, 'ข้าวมันไก่', 'Chicken rice', 1921143886227443712, 'lorem ipso', 80.00, true, 'https://images.unsplash.com/photo-1749640566096-5d8098d452b4?ixlib=rb-4.1.0&ixid=M3wxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8fA%3D%3D&auto=format&fit=crop&q=80&w=500', '2025-10-10 05:34:36.692235 +00:00', null, true),
        (1921822608437809152, 'แป๊ปซี่', 'Pepsi', 1921144050476388352, 'lorem ipso', 30.00, true, 'https://images.unsplash.com/photo-1651000877733-fe2c0a70b3cd?ixlib=rb-4.1.0&ixid=M3wxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8fA%3D%3D&auto=format&fit=crop&q=80&w=500', '2025-10-10 05:34:36.692235 +00:00', null, true),
        (1921822481287483392, 'เค้กแครอท', 'Carrot cake', 1921144250070732800, 'lorem ipso', 120.00, true, 'https://images.unsplash.com/photo-1622926421334-6829deee4b4b?ixlib=rb-4.1.0&ixid=M3wxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8fA%3D%3D&auto=format&fit=crop&q=80&w=500', '2025-10-10 05:34:36.692235 +00:00', null, true),
        (1921822481287483393, 'ต่อเวลา 15 นาที', 'Extension 15min', 1921144250070732801, null, 15.00, true, null, '2025-11-14 09:05:04.123826 +00:00', null, false),
        (1921822481287483394, 'ต่อเวลา 30 นาที', 'Extension 30min', 1921144250070732801, null, 30.00, true, null, '2025-11-14 09:05:04.123826 +00:00', null, false),
        (1921822481287483395, 'ต่อเวลา 60 นาที', 'Extension 60min', 1921144250070732801, null, 60.00, true, null, '2025-11-14 09:05:04.123826 +00:00', null, false),
        (1921829000000000104, 'ข้าวหน้าไก่เทอริยากิ', 'Teriyaki Chicken Rice', 1921143886227443712, 'Rice with teriyaki chicken', 90.00, true, 'https://images.unsplash.com/photo-1646809156467-6e825869b29f?q=80&w=2670&auto=format&fit=crop&ixlib=rb-4.1.0&ixid=M3wxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8fA%3D%3D&auto=format&fit=crop&q=80&w=500', '2025-11-14 09:05:04.123826 +00:00', null, true),
        (1921829000000000101, 'ผัดไทกุ้ง', 'Pad Thai Shrimp', 1921143886227443712, 'Signature Thai stir-fried noodles', 85.00, true, 'https://images.unsplash.com/photo-1729708475167-71a6eb3cd741?q=80&w=1470&auto=format&fit=crop&ixlib=rb-4.1.0&ixid=M3wxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8fA%3D%3D&auto=format&fit=crop&q=80&w=500', '2025-11-14 09:05:04.123826 +00:00', null, true),
        (1921829000000000108, 'บราวนี่', 'Brownie', 1921144250070732800, 'Chocolate brownie', 55.00, true, 'https://images.unsplash.com/photo-1560788843-8928d12b1600?q=80&w=1706&auto=format&fit=crop&ixlib=rb-4.1.0&ixid=M3wxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8fA%3D%3D&auto=format&fit=crop&q=80&w=500', '2025-11-14 09:05:04.123826 +00:00', null, true),
        (1921829000000000105, 'แกงเผ็ดเป็ดย่าง', 'Red Curry Roasted Duck', 1921143886227443712, 'Thai red curry duck', 140.00, true, 'https://images.unsplash.com/photo-1708782344490-9026aaa5eec7?q=80&w=1760&auto=format&fit=crop&ixlib=rb-4.1.0&ixid=M3wxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8fA%3D%3D&auto=format&fit=crop&q=80&w=500', '2025-11-14 09:05:04.123826 +00:00', null, true),
        (1921829000000000114, 'ชาเขียวมะลิ', 'Jasmine Green Tea', 1921144050476388352, 'Fragrant jasmine tea', 35.00, true, 'https://images.unsplash.com/photo-1630209712184-00101e51f374?q=80&w=1760&auto=format&fit=crop&ixlib=rb-4.1.0&ixid=M3wxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8fA%3D%3D&auto=format&fit=crop&q=80&w=500', '2025-11-14 09:05:04.123826 +00:00', null, true),
        (1921829000000000111, 'อเมริกาโน่เย็น', 'Iced Americano', 1921144050476388352, 'Strong black coffee', 45.00, true, 'https://images.unsplash.com/photo-1632789395770-20e6f63be806?q=80&w=1678&auto=format&fit=crop&ixlib=rb-4.1.0&ixid=M3wxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8fA%3D%3D&auto=format&fit=crop&q=80&w=500', '2025-11-14 09:05:04.123826 +00:00', null, true),
        (1921829000000000107, 'ฮันนี่โทสต์', 'Honey Toast', 1921144250070732800, 'Crispy toast with ice cream', 150.00, true, 'https://images.unsplash.com/photo-1737053566442-81cfe61f3357?q=80&w=1760&auto=format&fit=crop&ixlib=rb-4.1.0&ixid=M3wxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8fA%3D%3D&auto=format&fit=crop&q=80&w=500', '2025-11-14 09:05:04.123826 +00:00', null, true),
        (1921829000000000113, 'มัทฉะลาเต้', 'Matcha Latte', 1921144050476388352, 'Green tea latte', 65.00, true, 'https://images.unsplash.com/photo-1686794154608-e45c831ef567?q=80&w=1760&auto=format&fit=crop&ixlib=rb-4.1.0&ixid=M3wxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8fA%3D%3D&auto=format&fit=crop&q=80&w=500', '2025-11-14 09:05:04.123826 +00:00', null, true),
        (1921829000000000112, 'คาปูชิโน่ร้อน', 'Hot Cappuccino', 1921144050476388352, 'Milk coffee with foam', 50.00, true, 'https://images.unsplash.com/photo-1751976803077-2da275a7ad2c?q=80&w=1760&auto=format&fit=crop&ixlib=rb-4.1.0&ixid=M3wxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8fA%3D%3D&auto=format&fit=crop&q=80&w=500', '2025-11-14 09:05:04.123826 +00:00', null, true),
        (1921829000000000110, 'ชีสพายสตรอว์เบอร์รี', 'Strawberry Cheese Pie', 1921144250070732800, 'Fresh strawberry pie', 75.00, true, 'https://images.unsplash.com/photo-1622621746668-59fb299bc4d7?q=80&w=1833&auto=format&fit=crop&ixlib=rb-4.1.0&ixid=M3wxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8fA%3D%3D&auto=format&fit=crop&q=80&w=500', '2025-11-14 09:05:04.123826 +00:00', null, true),
        (1921829000000000102, 'ต้มยำกุ้ง', 'Tom Yum Kung', 1921143886227443712, 'Spicy Thai soup with shrimp', 120.00, true, 'https://images.unsplash.com/photo-1628430043175-0e8820df47c3?q=80&w=2671&auto=format&fit=crop&ixlib=rb-4.1.0&ixid=M3wxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8fA%3D%3D&auto=format&fit=crop&q=80&w=500', '2025-11-14 09:05:04.123826 +00:00', null, true),
        (1921829000000000103, 'ลาบหมู', 'Larb Moo', 1921143886227443712, 'Spicy minced pork salad', 70.00, true, 'https://images.unsplash.com/photo-1673238111115-18d3da6ec22b?q=80&w=1760&auto=format&fit=crop&ixlib=rb-4.1.0&ixid=M3wxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8fA%3D%3D&auto=format&fit=crop&q=80&w=500', '2025-11-14 09:05:04.123826 +00:00', null, true),
        (1921829000000000106, 'เครปเค้กสายรุ้ง', 'Rainbow Crepe Cake', 1921144250070732800, 'Colorful crepe cake with sauce', 95.00, true, 'https://images.unsplash.com/photo-1587131791788-ec0cb04e0a52?q=80&w=1760&auto=format&fit=crop&ixlib=rb-4.1.0&ixid=M3wxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8fA%3D%3D&auto=format&fit=crop&q=80&w=500', '2025-11-14 09:05:04.123826 +00:00', null, true),
        (1921829000000000109, 'ไอศกรีมชาไทย', 'Thai Tea Ice Cream', 1921144250070732800, 'Thai tea flavored ice cream', 45.00, true, 'https://images.unsplash.com/photo-1561230101-2c841778f9ae?q=80&w=1760&auto=format&fit=crop&ixlib=rb-4.1.0&ixid=M3wxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8fA%3D%3D&auto=format&fit=crop&q=80&w=500', '2025-11-14 09:05:04.123826 +00:00', null, true),
        (1921829000000000115, 'น้ำแดงโซดา', 'Red Soda', 1921144050476388352, 'Refreshing soda drink', 25.00, true, 'https://images.unsplash.com/photo-1631503190221-0f6a15367926?q=80&w=1760&auto=format&fit=crop&ixlib=rb-4.1.0&ixid=M3wxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8fA%3D%3D&auto=format&fit=crop&q=80&w=500', '2025-11-14 09:05:04.123826 +00:00', null, true),
        (1921829000000000201, 'ผัดกะเพราไก่ไข่ดาว', 'Basil Chicken with Fried Egg', 1921143886227443712, 'Spicy stir-fried chicken with basil and fried egg', 75.00, true, 'https://images.unsplash.com/photo-1707897634981-39bcfe435268?q=80&w=1796&auto=format&fit=crop&ixlib=rb-4.1.0&ixid=M3wxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8fA%3D%3D&auto=format&fit=crop&q=80&w=500', '2025-11-14 09:05:04.123826 +00:00', null, true),
        (1921829000000000203, 'เกี๊ยวน้ำหมูสับ', 'Pork Dumpling Soup', 1921143886227443712, 'Dumplings with minced pork in clear soup', 65.00, true, 'https://images.unsplash.com/photo-1571809864118-d0a73b090d6e?q=80&w=1760&auto=format&fit=crop&ixlib=rb-4.1.0&ixid=M3wxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8fA%3D%3D&auto=format&fit=crop&q=80&w=500', '2025-11-14 09:05:04.123826 +00:00', null, true),
        (1921829000000000202, 'หมูทอดกระเทียม', 'Garlic Fried Pork', 1921143886227443712, 'Deep-fried pork with garlic', 80.00, true, 'https://images.unsplash.com/photo-1709392975965-00889c6aa545?q=80&w=1724&auto=format&fit=crop&ixlib=rb-4.1.0&ixid=M3wxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8fA%3D%3D&auto=format&fit=crop&q=80&w=500', '2025-11-14 09:05:04.123826 +00:00', null, true),
        (1921829000000000204, 'ปลาทอดสมุนไพร', 'Fried Fish with Herbs', 1921143886227443712, 'Crispy fried fish with Thai herbs', 150.00, true, 'https://images.unsplash.com/photo-1700760933848-194ad3408fc4?q=80&w=1760&auto=format&fit=crop&ixlib=rb-4.1.0&ixid=M3wxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8fA%3D%3D&auto=format&fit=crop&q=80&w=500', '2025-11-14 09:05:04.123826 +00:00', null, true),
        (1921829000000000205, 'คอหมูย่าง', 'Grilled Pork Neck', 1921143886227443712, 'Thai-style grilled pork neck with dipping sauce', 95.00, true, 'https://images.unsplash.com/photo-1708615017161-2eff302d0389?q=80&w=1760&auto=format&fit=crop&ixlib=rb-4.1.0&ixid=M3wxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8fA%3D%3D&auto=format&fit=crop&q=80&w=500', '2025-11-14 09:05:04.123826 +00:00', null, true);

insert into public.product_time_extension (id, duration_minutes, products_id, created_at, updated_at)
values  (1921822481287483396, 15, 1921822481287483393, '2025-11-14 09:05:56.706476 +00:00', null),
        (1921822481287483397, 30, 1921822481287483394, '2025-11-14 09:05:56.706476 +00:00', null),
        (1921822481287483398, 60, 1921822481287483395, '2025-11-14 09:05:56.706476 +00:00', null);

insert into public.tables (id, table_number, status_id, seats, created_at, updated_at)
values  (1920153361642950656, 5, 1919968486671519744, 4, '2025-11-14 09:06:41.530788 +00:00', null),
        (1919972141986484224, 3, 1919968486671519744, 4, '2025-11-14 09:06:41.530788 +00:00', null),
        (1919996486741921792, 4, 1919968486671519744, 5, '2025-11-14 09:06:41.530788 +00:00', null),
        (1919971956241731584, 2, 1919968486671519744, 3, '2025-11-14 09:06:41.530788 +00:00', null),
        (1919968785813475328, 1, 1919968486847680512, 5, '2025-11-14 09:06:41.530788 +00:00', '2025-11-15 04:34:22.802407 +00:00'),
        (1989992002053308903, 15, 1919968486671519744, 8, '2025-11-14 09:06:41.530788 +00:00', null),
        (1989992002056393188, 12, 1919968486671519744, 2, '2025-11-14 09:06:41.530788 +00:00', null),
        (1989992002055795167, 7, 1919968486671519744, 4, '2025-11-14 09:06:41.530788 +00:00', null),
        (1989992002055377385, 17, 1919968486671519744, 4, '2025-11-14 09:06:41.530788 +00:00', null),
        (1989992002056450528, 8, 1919968486671519744, 4, '2025-11-14 09:06:41.530788 +00:00', null),
        (1989992002053304803, 11, 1919968486671519744, 5, '2025-11-14 09:06:41.530788 +00:00', null),
        (1989992002055102949, 13, 1919968486671519744, 4, '2025-11-14 09:06:41.530788 +00:00', null),
        (1989992002057384417, 9, 1919968486671519744, 6, '2025-11-14 09:06:41.530788 +00:00', null),
        (1989992002054472158, 6, 1919968486671519744, 2, '2025-11-14 09:06:41.530788 +00:00', null),
        (1989992002054017514, 18, 1919968486671519744, 5, '2025-11-14 09:06:41.530788 +00:00', null),
        (1989992002055758310, 14, 1919968486671519744, 6, '2025-11-14 09:06:41.530788 +00:00', null),
        (1989992002055954920, 16, 1919968486671519744, 2, '2025-11-14 09:06:41.530788 +00:00', null),
        (1989992002055152098, 10, 1919968486671519744, 3, '2025-11-14 09:06:41.530788 +00:00', null),
        (1989992002056503787, 19, 1919968486671519744, 3, '2025-11-14 09:06:41.530788 +00:00', null),
        (1989992002055467500, 20, 1919968486671519744, 4, '2025-11-14 09:06:41.530788 +00:00', null);