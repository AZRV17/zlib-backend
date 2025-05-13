--
-- PostgreSQL database dump
--

-- Dumped from database version 14.17 (Ubuntu 14.17-0ubuntu0.22.04.1)
-- Dumped by pg_dump version 14.17 (Ubuntu 14.17-0ubuntu0.22.04.1)

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

--
-- Name: role; Type: TYPE; Schema: public; Owner: postgres
--

CREATE TYPE public.role AS ENUM (
    'user',
    'admin',
    'librarian'
);


ALTER TYPE public.role OWNER TO postgres;

--
-- Name: sign_in(text, text); Type: FUNCTION; Schema: public; Owner: postgres
--

CREATE FUNCTION public.sign_in(login text, pass text) RETURNS TABLE(id integer, role text)
    LANGUAGE plpgsql
    AS $$
    BEGIN
        RETURN QUERY
        SELECT id, role
        FROM users
        WHERE login = login AND password = crypt(pass, password);

        IF NOT FOUND THEN
            RAISE EXCEPTION 'Invalid login or password';
        END IF;
    END;
    $$;


ALTER FUNCTION public.sign_in(login text, pass text) OWNER TO postgres;

SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: audiobook_files; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.audiobook_files (
    id bigint NOT NULL,
    book_id bigint,
    file_path text,
    chapter_title text,
    "order" bigint
);


ALTER TABLE public.audiobook_files OWNER TO postgres;

--
-- Name: audiobook_files_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.audiobook_files_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.audiobook_files_id_seq OWNER TO postgres;

--
-- Name: audiobook_files_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.audiobook_files_id_seq OWNED BY public.audiobook_files.id;


--
-- Name: authors; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.authors (
    id bigint NOT NULL,
    name text,
    lastname text,
    biography text,
    birthdate timestamp with time zone,
    created_at timestamp with time zone,
    updated_at timestamp with time zone
);


ALTER TABLE public.authors OWNER TO postgres;

--
-- Name: authors_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.authors_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.authors_id_seq OWNER TO postgres;

--
-- Name: authors_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.authors_id_seq OWNED BY public.authors.id;


--
-- Name: books; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.books (
    id bigint NOT NULL,
    title text,
    author_id bigint,
    genre_id bigint,
    description text,
    publisher_id bigint,
    isbn bigint,
    year_of_publication timestamp with time zone,
    picture text,
    rating numeric,
    is_available boolean DEFAULT true,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    epub_file text
);


ALTER TABLE public.books OWNER TO postgres;

--
-- Name: books_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.books_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.books_id_seq OWNER TO postgres;

--
-- Name: books_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.books_id_seq OWNED BY public.books.id;


--
-- Name: chats; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.chats (
    id bigint NOT NULL,
    user_id bigint,
    librarian_id bigint,
    status text,
    title text,
    last_activity timestamp with time zone,
    created_at timestamp with time zone
);


ALTER TABLE public.chats OWNER TO postgres;

--
-- Name: chats_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.chats_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.chats_id_seq OWNER TO postgres;

--
-- Name: chats_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.chats_id_seq OWNED BY public.chats.id;


--
-- Name: favorites; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.favorites (
    id bigint NOT NULL,
    user_id bigint,
    book_id bigint,
    created_at timestamp with time zone,
    updated_at timestamp with time zone
);


ALTER TABLE public.favorites OWNER TO postgres;

--
-- Name: favorites_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.favorites_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.favorites_id_seq OWNER TO postgres;

--
-- Name: favorites_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.favorites_id_seq OWNED BY public.favorites.id;


--
-- Name: genres; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.genres (
    id bigint NOT NULL,
    name text,
    description text,
    created_at timestamp with time zone,
    updated_at timestamp with time zone
);


ALTER TABLE public.genres OWNER TO postgres;

--
-- Name: genres_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.genres_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.genres_id_seq OWNER TO postgres;

--
-- Name: genres_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.genres_id_seq OWNED BY public.genres.id;


--
-- Name: logs; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.logs (
    id bigint NOT NULL,
    user_id bigint,
    action text,
    date timestamp with time zone,
    details text
);


ALTER TABLE public.logs OWNER TO postgres;

--
-- Name: logs_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.logs_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.logs_id_seq OWNER TO postgres;

--
-- Name: logs_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.logs_id_seq OWNED BY public.logs.id;


--
-- Name: messages; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.messages (
    id bigint NOT NULL,
    chat_id bigint,
    sender_id bigint,
    sender_role text,
    sender_name text,
    content text,
    read_at timestamp with time zone,
    created_at timestamp with time zone
);


ALTER TABLE public.messages OWNER TO postgres;

--
-- Name: messages_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.messages_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.messages_id_seq OWNER TO postgres;

--
-- Name: messages_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.messages_id_seq OWNED BY public.messages.id;


--
-- Name: notifications; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.notifications (
    id bigint NOT NULL,
    user_id bigint,
    message text,
    is_read boolean DEFAULT false,
    created_at timestamp with time zone,
    updated_at timestamp with time zone
);


ALTER TABLE public.notifications OWNER TO postgres;

--
-- Name: notifications_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.notifications_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.notifications_id_seq OWNER TO postgres;

--
-- Name: notifications_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.notifications_id_seq OWNED BY public.notifications.id;


--
-- Name: publishers; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.publishers (
    id bigint NOT NULL,
    name text,
    created_at timestamp with time zone,
    updated_at timestamp with time zone
);


ALTER TABLE public.publishers OWNER TO postgres;

--
-- Name: publishers_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.publishers_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.publishers_id_seq OWNER TO postgres;

--
-- Name: publishers_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.publishers_id_seq OWNED BY public.publishers.id;


--
-- Name: reservations; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.reservations (
    id bigint NOT NULL,
    user_id bigint,
    book_id bigint,
    date_of_issue timestamp with time zone,
    date_of_return timestamp with time zone,
    status text,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    unique_code_id bigint
);


ALTER TABLE public.reservations OWNER TO postgres;

--
-- Name: reservations_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.reservations_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.reservations_id_seq OWNER TO postgres;

--
-- Name: reservations_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.reservations_id_seq OWNED BY public.reservations.id;


--
-- Name: reviews; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.reviews (
    id bigint NOT NULL,
    user_id bigint,
    book_id bigint,
    rating numeric,
    message text,
    created_at timestamp with time zone,
    updated_at timestamp with time zone
);


ALTER TABLE public.reviews OWNER TO postgres;

--
-- Name: reviews_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.reviews_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.reviews_id_seq OWNER TO postgres;

--
-- Name: reviews_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.reviews_id_seq OWNED BY public.reviews.id;


--
-- Name: unique_codes; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.unique_codes (
    id bigint NOT NULL,
    code bigint,
    book_id bigint,
    is_available boolean
);


ALTER TABLE public.unique_codes OWNER TO postgres;

--
-- Name: unique_codes_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.unique_codes_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.unique_codes_id_seq OWNER TO postgres;

--
-- Name: unique_codes_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.unique_codes_id_seq OWNED BY public.unique_codes.id;


--
-- Name: users; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.users (
    id bigint NOT NULL,
    login text,
    full_name text,
    password text,
    role public.role DEFAULT 'user'::public.role,
    email text,
    phone_number text,
    passport_number bigint,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    reset_password_token character varying(100),
    reset_token_expiry timestamp with time zone
);


ALTER TABLE public.users OWNER TO postgres;

--
-- Name: users_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.users_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.users_id_seq OWNER TO postgres;

--
-- Name: users_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.users_id_seq OWNED BY public.users.id;


--
-- Name: audiobook_files id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.audiobook_files ALTER COLUMN id SET DEFAULT nextval('public.audiobook_files_id_seq'::regclass);


--
-- Name: authors id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.authors ALTER COLUMN id SET DEFAULT nextval('public.authors_id_seq'::regclass);


--
-- Name: books id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.books ALTER COLUMN id SET DEFAULT nextval('public.books_id_seq'::regclass);


--
-- Name: chats id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.chats ALTER COLUMN id SET DEFAULT nextval('public.chats_id_seq'::regclass);


--
-- Name: favorites id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.favorites ALTER COLUMN id SET DEFAULT nextval('public.favorites_id_seq'::regclass);


--
-- Name: genres id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.genres ALTER COLUMN id SET DEFAULT nextval('public.genres_id_seq'::regclass);


--
-- Name: logs id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.logs ALTER COLUMN id SET DEFAULT nextval('public.logs_id_seq'::regclass);


--
-- Name: messages id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.messages ALTER COLUMN id SET DEFAULT nextval('public.messages_id_seq'::regclass);


--
-- Name: notifications id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.notifications ALTER COLUMN id SET DEFAULT nextval('public.notifications_id_seq'::regclass);


--
-- Name: publishers id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.publishers ALTER COLUMN id SET DEFAULT nextval('public.publishers_id_seq'::regclass);


--
-- Name: reservations id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.reservations ALTER COLUMN id SET DEFAULT nextval('public.reservations_id_seq'::regclass);


--
-- Name: reviews id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.reviews ALTER COLUMN id SET DEFAULT nextval('public.reviews_id_seq'::regclass);


--
-- Name: unique_codes id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.unique_codes ALTER COLUMN id SET DEFAULT nextval('public.unique_codes_id_seq'::regclass);


--
-- Name: users id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.users ALTER COLUMN id SET DEFAULT nextval('public.users_id_seq'::regclass);


--
-- Data for Name: audiobook_files; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.audiobook_files (id, book_id, file_path, chapter_title, "order") FROM stdin;
1	11	uploads/audio/f79f6cc9-28f4-440e-a85d-437ae386869d.mp3	Первая глава	1
2	11	uploads/audio/072b50df-1e7f-4c46-a42a-e8d4f84463aa.mp3	Вторая глава	2
3	20	uploads/audio/a6689b3c-ec1b-4e12-ad84-19144dca4d41.mp3	Первая	1
\.


--
-- Data for Name: authors; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.authors (id, name, lastname, biography, birthdate, created_at, updated_at) FROM stdin;
2	Федор	Достаевский	Тоже драмматург	2024-11-22 03:00:00+03	0001-01-01 02:30:17+02:30:17	2024-11-03 15:25:56.705955+03
2289	Федор12	Достаевский12	Тоже драмматург	2024-11-22 03:00:00+03	2025-05-07 19:01:19.89153+03	2025-05-07 19:01:19.89153+03
3	Фердинанд	Селин	Француз	1999-10-10 04:00:00+04	2024-11-25 11:25:37.353251+03	2024-11-25 11:25:37.353251+03
1	Лев	Толстой	Драмматург, великий	2024-10-11 03:00:00+03	0001-01-01 02:30:17+02:30:17	2025-04-06 20:52:21.701615+03
\.


--
-- Data for Name: books; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.books (id, title, author_id, genre_id, description, publisher_id, isbn, year_of_publication, picture, rating, is_available, created_at, updated_at, epub_file) FROM stdin;
3	Преступление и наказание	2	1	Про убийство бабки	1	24114	2024-10-22 03:00:00+03		0	t	2024-10-22 20:29:37.478+03	2024-11-25 11:35:23.717802+03	\N
1	Война и мир	1	1	Про 1812	1	12312	2024-10-22 20:27:51.823+03	https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcSzS2ywJwuOso4hQdSM-mHnIGWeNpnYOvriQA&s	4.5	t	2024-10-22 20:28:32.727+03	2024-11-25 11:36:07.869008+03	\N
2	Бесы	2	2	Про бесовщину	1	12424	2024-10-22 03:00:00+03		0	t	2024-10-22 20:29:05.228+03	2024-11-03 17:19:39.629851+03	\N
8	Путешествие на край ночи	3	1	Книга про первую мировую	1	124251	2000-11-11 03:00:00+03	uploads/books/e40b63cb-b32e-46c8-8436-3da5717f4e85.jpg	0	t	2024-11-25 11:37:06.240883+03	2024-11-25 11:37:11.25027+03	\N
20	Новая1	3	1	new	1	1242125	2000-11-11 03:00:00+03	http://localhost:8080/uploads/books/c7c49dcb-2250-482f-aed8-b0454ecb719f.jpg	4	t	2025-04-14 09:36:06.879422+03	2025-05-06 14:22:26.103383+03	uploads/books/c7c49dcb-2250-482f-aed8-b0454ecb719f.epub
11	new	2	1	dasdsa	1	1241242	2021-02-02 03:00:00+03	http://localhost:8080/uploads/books/953f68ff-9e3a-45ea-99b8-4d29e997e9f6.jpg	3	t	2025-01-28 12:53:08.433643+03	2025-05-07 18:39:46.995167+03	uploads/books/ddab2652-4dc4-42da-af5b-68f23304036b.epub
\.


--
-- Data for Name: chats; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.chats (id, user_id, librarian_id, status, title, last_activity, created_at) FROM stdin;
3	1	5	closed	Обращение в поддержку	2025-03-14 20:29:30.399973+03	2025-03-14 12:11:33.557924+03
4	1	5	closed	Новое обращение	2025-03-14 20:44:04.918513+03	2025-03-14 20:42:43.183415+03
5	1	5	closed	новое	2025-03-21 12:33:47.656636+03	2025-03-21 12:33:26.065593+03
7	14	5	closed	новое	2025-04-01 12:45:55.696146+03	2025-04-01 12:45:27.757119+03
8	1	5	closed	Новое обращение об ошибке	2025-04-06 20:53:51.093755+03	2025-04-06 20:51:49.027694+03
9	16	5	closed	Проблема	2025-04-14 09:37:37.915538+03	2025-04-14 09:31:32.787005+03
6	14	5	closed	Проблема	2025-04-14 09:37:39.709378+03	2025-03-31 20:29:06.440142+03
10	1	5	closed	новое обращение	2025-04-29 13:25:39.855723+03	2025-04-29 13:24:11.302923+03
\.


--
-- Data for Name: favorites; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.favorites (id, user_id, book_id, created_at, updated_at) FROM stdin;
26	1	1	2025-03-14 20:57:08.934758+03	2025-03-14 20:57:08.934758+03
27	1	8	2025-03-14 20:57:34.897455+03	2025-03-14 20:57:34.897455+03
29	15	11	2025-04-14 09:25:55.363986+03	2025-04-14 09:25:55.363986+03
30	16	11	2025-04-14 09:30:36.449679+03	2025-04-14 09:30:36.449679+03
31	17	11	2025-04-14 09:43:52.964264+03	2025-04-14 09:43:52.964264+03
33	1	11	2025-04-29 13:22:58.609558+03	2025-04-29 13:22:58.609558+03
\.


--
-- Data for Name: genres; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.genres (id, name, description, created_at, updated_at) FROM stdin;
1	Роман	Про всякие любовные штучки	2024-10-22 20:26:55.973+03	2024-10-22 20:26:56.226+03
2	Анти утопия	Про некрасивые миры	0001-01-01 02:30:17+02:30:17	2024-11-03 17:38:17.982919+03
\.


--
-- Data for Name: logs; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.logs (id, user_id, action, date, details) FROM stdin;
1	4	Изменение роли пользователя	2024-11-07 18:02:59.253725+03	Изменение роли пользователя
2	1	Добавление книги в избранное	2024-11-08 12:09:00.780813+03	Добавление книги в избранное
3	1	Удаление из избранного	2024-11-08 12:09:02.937071+03	Удаление из избранного
4	4	Изменение роли пользователя	2024-11-08 12:10:33.356029+03	Изменение роли пользователя
5	4	Создание бэкапа	2024-11-08 12:10:42.423444+03	Создание бэкапа
6	4	Создание бэкапа	2024-11-11 22:30:53.572082+03	Создание бэкапа
7	4	Восстановление бд	2024-11-11 22:32:33.497669+03	Восстановление бд
8	5	Экспорт авторов в CSV	2024-11-11 22:45:18.93043+03	Экспорт авторов в CSV
9	5	Создание книги	2024-11-12 12:27:15.348589+03	Создание книги
10	5	Удаление книги	2024-11-13 11:11:49.068987+03	Удаление книги
11	5	Создание книги	2024-11-13 11:13:15.564241+03	Создание книги
12	4	Создание бэкапа	2024-11-15 12:18:00.173604+03	Создание бэкапа
13	4	Восстановление бд	2024-11-15 12:18:07.235874+03	Восстановление бд
14	5	Экспорт авторов в CSV	2024-11-15 12:18:19.366481+03	Экспорт авторов в CSV
15	5	Удаление книги	2024-11-19 23:08:02.953402+03	Удаление книги
16	1	Добавление книги в избранное	2024-11-19 23:15:48.335561+03	Добавление книги в избранное
17	5	Экспорт авторов в CSV	2024-11-19 23:16:31.094741+03	Экспорт авторов в CSV
18	4	Создание бэкапа	2024-11-19 23:18:25.508844+03	Создание бэкапа
19	4	Создание бэкапа	2024-11-19 23:18:31.935267+03	Создание бэкапа
20	1	Удаление из избранного	2024-11-22 12:24:19.676247+03	Удаление из избранного
24	5	Изменение статуса бронирования	2024-11-25 11:03:36.935535+03	Изменение статуса бронирования
25	5	Изменение статуса бронирования	2024-11-25 11:03:39.722592+03	Изменение статуса бронирования
26	4	Изменение роли пользователя	2024-11-25 11:03:50.847784+03	Изменение роли пользователя
27	4	Изменение роли пользователя	2024-11-25 11:03:53.112931+03	Изменение роли пользователя
35	14	Добавление книги в избранное	2024-11-25 11:24:24.425459+03	Добавление книги в избранное
36	14	Бронирование книги	2024-11-25 11:24:25.426119+03	Бронирование книги
37	14	Добавление отзыва	2024-11-25 11:24:33.077915+03	Добавление отзыва
38	14	Удаление из избранного	2024-11-25 11:24:38.879028+03	Удаление из избранного
39	5	Экспорт авторов в CSV	2024-11-25 11:25:05.700654+03	Экспорт авторов в CSV
40	5	Обновление автора	2024-11-25 11:25:14.407975+03	Обновление автора
41	5	Создание автора	2024-11-25 11:25:37.359794+03	Создание автора
42	5	Создание книги	2024-11-25 11:27:25.973456+03	Создание книги
43	5	Удаление книги	2024-11-25 11:29:12.408982+03	Удаление книги
44	5	Создание книги	2024-11-25 11:30:51.037738+03	Создание книги
45	5	Изменение книги	2024-11-25 11:33:13.741666+03	Изменение книги
46	5	Изменение книги	2024-11-25 11:33:33.426339+03	Изменение книги
47	5	Изменение книги	2024-11-25 11:35:23.720171+03	Изменение книги
48	5	Изменение книги	2024-11-25 11:36:07.871067+03	Изменение книги
49	5	Удаление книги	2024-11-25 11:36:29.816353+03	Удаление книги
50	5	Создание книги	2024-11-25 11:37:06.246563+03	Создание книги
51	5	Изменение книги	2024-11-25 11:37:11.253537+03	Изменение книги
52	5	Создание жанра	2024-11-25 11:37:32.682242+03	Создание жанра
53	5	Удаление жанра	2024-11-25 11:37:34.648595+03	Удаление жанра
54	5	Создание издателя	2024-11-25 11:37:42.123747+03	Создание издателя
55	5	Изменение статуса бронирования	2024-11-25 11:37:53.661677+03	Изменение статуса бронирования
56	4	Изменение роли пользователя	2024-11-25 11:38:12.645099+03	Изменение роли пользователя
57	4	Создание бэкапа	2024-11-25 11:38:28.025279+03	Создание бэкапа
58	4	Восстановление бд	2024-11-25 11:38:34.50952+03	Восстановление бд
59	5	Создание книги	2025-01-28 12:42:07.014433+03	Создание книги
60	5	Удаление книги	2025-01-28 12:44:33.661283+03	Удаление книги
61	5	Создание книги	2025-01-28 12:49:25.603058+03	Создание книги
62	5	Удаление книги	2025-01-28 12:50:19.84355+03	Удаление книги
63	5	Создание книги	2025-01-28 12:53:08.438363+03	Создание книги
64	1	Добавление книги в избранное	2025-01-28 13:11:28.82083+03	Добавление книги в избранное
65	1	Удаление из избранного	2025-01-28 13:11:44.665114+03	Удаление из избранного
66	1	Добавление книги в избранное	2025-02-04 12:53:01.094226+03	Добавление книги в избранное
67	5	Создание издателя	2025-03-02 13:39:28.236756+03	Создание издателя
68	5	Удаление издателя	2025-03-02 13:39:33.556012+03	Удаление издателя
69	1	Добавление книги в избранное	2025-03-14 20:57:08.939862+03	Добавление книги в избранное
70	1	Добавление книги в избранное	2025-03-14 20:57:34.898883+03	Добавление книги в избранное
71	4	Изменение роли пользователя	2025-03-31 10:56:20.222944+03	Изменение роли пользователя
72	14	Бронирование книги	2025-03-31 10:59:14.684655+03	Бронирование книги
73	5	Изменение статуса бронирования	2025-03-31 11:09:09.216486+03	Изменение статуса бронирования
74	5	Изменение статуса бронирования	2025-03-31 11:14:56.730302+03	Изменение статуса бронирования
75	5	Изменение статуса бронирования	2025-03-31 11:15:08.093482+03	Изменение статуса бронирования
76	5	Изменение книги	2025-03-31 12:09:27.323899+03	Изменение книги
77	5	Изменение книги	2025-03-31 12:12:56.305261+03	Изменение книги
78	1	Удаление из избранного	2025-04-06 20:50:58.845817+03	Удаление из избранного
79	1	Добавление книги в избранное	2025-04-06 20:51:00.255405+03	Добавление книги в избранное
80	5	Обновление автора	2025-04-06 20:52:21.706724+03	Обновление автора
81	5	Создание автора	2025-04-06 20:52:38.240714+03	Создание автора
82	5	Удаление автора	2025-04-06 20:52:40.640372+03	Удаление автора
83	5	Экспорт авторов в CSV	2025-04-06 20:52:43.830555+03	Экспорт авторов в CSV
84	5	Создание жанра	2025-04-06 20:52:53.192422+03	Создание жанра
85	5	Удаление жанра	2025-04-06 20:52:55.225029+03	Удаление жанра
86	5	Создание издателя	2025-04-06 20:53:02.909708+03	Создание издателя
87	5	Удаление издателя	2025-04-06 20:53:05.330151+03	Удаление издателя
88	5	Изменение уникального кода	2025-04-06 20:53:21.260277+03	Изменение уникального кода
89	5	Изменение уникального кода	2025-04-06 20:53:24.202517+03	Изменение уникального кода
90	5	Изменение статуса бронирования	2025-04-06 20:53:31.524463+03	Изменение статуса бронирования
91	5	Изменение статуса бронирования	2025-04-06 20:53:34.44689+03	Изменение статуса бронирования
92	5	Изменение статуса бронирования	2025-04-06 20:53:36.937756+03	Изменение статуса бронирования
93	4	Изменение роли пользователя	2025-04-06 20:54:14.06507+03	Изменение роли пользователя
94	4	Изменение роли пользователя	2025-04-06 20:54:17.484886+03	Изменение роли пользователя
95	4	Изменение роли пользователя	2025-04-06 20:54:19.571365+03	Изменение роли пользователя
96	4	Создание бэкапа	2025-04-06 20:54:25.056088+03	Создание бэкапа
97	5	Изменение книги	2025-04-14 09:13:16.66873+03	Изменение книги
98	5	Изменение книги	2025-04-14 09:13:24.597361+03	Изменение книги
99	5	Создание уникального кода	2025-04-14 09:13:38.062506+03	Создание уникального кода
100	15	Добавление книги в избранное	2025-04-14 09:25:55.365183+03	Добавление книги в избранное
101	15	Бронирование книги	2025-04-14 09:25:58.459366+03	Бронирование книги
102	5	Изменение статуса бронирования	2025-04-14 09:29:45.828924+03	Изменение статуса бронирования
103	16	Добавление книги в избранное	2025-04-14 09:30:36.45163+03	Добавление книги в избранное
104	16	Бронирование книги	2025-04-14 09:30:37.934531+03	Бронирование книги
105	5	Создание автора	2025-04-14 09:31:56.121574+03	Создание автора
106	5	Обновление автора	2025-04-14 09:32:00.236668+03	Обновление автора
107	5	Удаление автора	2025-04-14 09:32:01.298621+03	Удаление автора
108	5	Экспорт авторов в CSV	2025-04-14 09:32:02.324269+03	Экспорт авторов в CSV
109	5	Создание жанра	2025-04-14 09:32:11.346072+03	Создание жанра
110	5	Удаление жанра	2025-04-14 09:32:13.095009+03	Удаление жанра
111	5	Создание автора	2025-04-14 09:35:33.124373+03	Создание автора
112	5	Обновление автора	2025-04-14 09:35:35.935039+03	Обновление автора
113	5	Удаление автора	2025-04-14 09:35:37.278276+03	Удаление автора
114	5	Экспорт авторов в CSV	2025-04-14 09:35:38.024561+03	Экспорт авторов в CSV
115	5	Создание книги	2025-04-14 09:36:06.882841+03	Создание книги
116	5	Создание жанра	2025-04-14 09:36:46.569367+03	Создание жанра
117	5	Удаление жанра	2025-04-14 09:36:48.0658+03	Удаление жанра
118	5	Создание издателя	2025-04-14 09:36:51.387607+03	Создание издателя
119	5	Изменение издателя	2025-04-14 09:36:53.968987+03	Изменение издателя
120	5	Удаление издателя	2025-04-14 09:36:54.782627+03	Удаление издателя
121	5	Создание уникального кода	2025-04-14 09:37:05.997079+03	Создание уникального кода
122	5	Изменение статуса бронирования	2025-04-14 09:37:15.368565+03	Изменение статуса бронирования
123	5	Изменение статуса бронирования	2025-04-14 09:37:17.819886+03	Изменение статуса бронирования
124	4	Изменение роли пользователя	2025-04-14 09:37:55.911387+03	Изменение роли пользователя
125	4	Изменение роли пользователя	2025-04-14 09:37:58.031095+03	Изменение роли пользователя
126	4	Изменение роли пользователя	2025-04-14 09:38:03.120122+03	Изменение роли пользователя
127	4	Создание бэкапа	2025-04-14 09:38:07.42663+03	Создание бэкапа
128	17	Добавление книги в избранное	2025-04-14 09:43:52.966307+03	Добавление книги в избранное
129	17	Бронирование книги	2025-04-14 09:43:54.077576+03	Бронирование книги
2439	1	Добавление отзыва	2025-05-06 14:22:16.662589+03	Добавление отзыва
2440	1	Добавление отзыва	2025-05-06 14:22:26.105562+03	Добавление отзыва
2441	5	Импорт авторов из CSV	2025-05-07 18:32:27.467849+03	Импорт авторов из CSV
2442	5	Удаление автора	2025-05-07 18:32:49.653812+03	Удаление автора
2443	5	Изменение книги	2025-05-07 18:34:02.098853+03	Изменение книги
2444	5	Изменение книги	2025-05-07 18:35:03.990627+03	Изменение книги
2445	1	Добавление отзыва	2025-05-07 18:39:46.998754+03	Добавление отзыва
2446	1	Добавление книги в избранное	2025-05-07 18:56:16.902885+03	Добавление книги в избранное
2447	1	Удаление из избранного	2025-05-07 18:56:17.811004+03	Удаление из избранного
2448	5	Импорт авторов из CSV	2025-05-07 19:01:19.892973+03	Импорт авторов из CSV
2409	1	Удаление из избранного	2025-04-29 13:08:14.945438+03	Удаление из избранного
2410	1	Добавление книги в избранное	2025-04-29 13:08:15.886219+03	Добавление книги в избранное
2411	1	Удаление из избранного	2025-04-29 13:22:57.917848+03	Удаление из избранного
2412	1	Добавление книги в избранное	2025-04-29 13:22:58.611159+03	Добавление книги в избранное
2413	5	Экспорт авторов в CSV	2025-04-29 13:24:35.24778+03	Экспорт авторов в CSV
2414	5	Изменение статуса бронирования	2025-04-29 13:25:23.499812+03	Изменение статуса бронирования
2415	4	Создание бэкапа	2025-04-29 13:26:08.02635+03	Создание бэкапа
2416	5	Экспорт авторов в CSV	2025-05-03 14:22:47.009841+03	Экспорт авторов в CSV
2417	5	Изменение книги	2025-05-03 14:27:22.836519+03	Изменение книги
2418	5	Экспорт книг в CSV	2025-05-03 14:29:29.623891+03	Экспорт книг в CSV
2419	5	Экспорт жанров в CSV	2025-05-03 14:33:13.634474+03	Экспорт жанров в CSV
2420	5	Экспорт жанров в CSV	2025-05-03 14:34:00.210391+03	Экспорт жанров в CSV
2421	5	Экспорт жанров в CSV	2025-05-03 14:34:29.16138+03	Экспорт жанров в CSV
2422	5	Экспорт книг в CSV	2025-05-03 14:44:06.674768+03	Экспорт книг в CSV
2423	5	Экспорт книг в CSV	2025-05-03 14:44:36.184523+03	Экспорт книг в CSV
2424	5	Экспорт бронирований в CSV	2025-05-03 14:56:43.535903+03	Экспорт бронирований в CSV
2425	5	Экспорт книг в CSV	2025-05-03 15:41:01.97081+03	Экспорт книг в CSV
2426	5	Импорт книг из CSV	2025-05-03 15:41:44.450539+03	Импорт книг из CSV
2427	5	Импорт книг из CSV	2025-05-03 15:43:26.527769+03	Импорт книг из CSV
2428	5	Импорт книг из CSV	2025-05-03 15:45:02.265772+03	Импорт книг из CSV
2429	5	Удаление книги	2025-05-03 15:45:39.494629+03	Удаление книги
2430	5	Импорт книг из CSV	2025-05-03 15:47:24.122488+03	Импорт книг из CSV
2431	5	Импорт книг из CSV	2025-05-03 15:47:46.15618+03	Импорт книг из CSV
2432	5	Удаление книги	2025-05-03 15:47:51.173998+03	Удаление книги
2433	5	Экспорт авторов в CSV	2025-05-03 16:09:54.988515+03	Экспорт авторов в CSV
2434	5	Импорт авторов из CSV	2025-05-03 16:11:41.404134+03	Импорт авторов из CSV
2435	5	Удаление автора	2025-05-03 16:11:52.124938+03	Удаление автора
2436	5	Экспорт жанров в CSV	2025-05-03 16:13:28.200861+03	Экспорт жанров в CSV
2437	5	Импорт жанров из CSV	2025-05-03 16:13:50.209796+03	Импорт жанров из CSV
2438	5	Удаление жанра	2025-05-03 16:13:53.131647+03	Удаление жанра
\.


--
-- Data for Name: messages; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.messages (id, chat_id, sender_id, sender_role, sender_name, content, read_at, created_at) FROM stdin;
3	3	1	user	user	1	2025-03-14 12:23:12.263882+03	2025-03-14 12:11:45.63046+03
4	3	5	librarian	libr	2	2025-03-14 12:23:28.0982+03	2025-03-14 12:23:20.390909+03
5	3	1	user	user	3	\N	2025-03-14 20:29:30.399973+03
6	4	5	librarian	libr	Добрый день! У меня проблема	2025-03-14 20:44:44.788117+03	2025-03-14 20:43:38.715833+03
7	4	5	librarian	libr	Какая?	2025-03-14 20:44:44.788117+03	2025-03-14 20:43:53.798131+03
8	5	1	user	user	12312	2025-03-21 12:33:34.18165+03	2025-03-21 12:33:28.246193+03
9	5	5	librarian	libr	выфаыфа	2025-03-21 12:33:43.665236+03	2025-03-21 12:33:37.9589+03
10	6	14	user	newUser	У меня возникла проблема	2025-04-01 12:45:48.36771+03	2025-03-31 20:29:20.782902+03
11	7	5	librarian	libr	sdaf	\N	2025-04-01 12:45:53.507964+03
12	8	1	user	user	У меня проблема	2025-04-06 20:53:43.588617+03	2025-04-06 20:51:55.516723+03
13	8	5	librarian	libr	Какая?	2025-04-14 09:12:05.109295+03	2025-04-06 20:53:48.245997+03
14	9	16	user	newUser3	У меня проблема	2025-04-14 09:37:26.91222+03	2025-04-14 09:31:36.630128+03
15	9	5	librarian	libr	Опишите проблему	\N	2025-04-14 09:37:35.698997+03
16	10	1	user	user	ыалвфодаф	2025-04-29 13:25:35.01582+03	2025-04-29 13:24:16.424567+03
17	10	5	librarian	libr	dkhsakdjhsakj	2025-05-12 18:03:31.501811+03	2025-04-29 13:25:38.434699+03
\.


--
-- Data for Name: notifications; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.notifications (id, user_id, message, is_read, created_at, updated_at) FROM stdin;
\.


--
-- Data for Name: publishers; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.publishers (id, name, created_at, updated_at) FROM stdin;
1	Питер	2024-10-22 20:26:37.677+03	2024-11-03 17:48:09.058064+03
2	Новое издание	2024-11-25 11:37:42.120214+03	2024-11-25 11:37:42.120214+03
\.


--
-- Data for Name: reservations; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.reservations (id, user_id, book_id, date_of_issue, date_of_return, status, created_at, updated_at, unique_code_id) FROM stdin;
4	1	2	2024-11-02 14:10:34.540667+03	2024-11-09 14:10:34.540668+03	returned	2024-11-02 14:10:34.540704+03	2024-11-03 23:45:51.830777+03	2
8	14	3	2024-11-25 11:24:25.42319+03	2024-12-02 11:24:25.423191+03	returned	2024-11-25 11:24:25.423212+03	2024-11-25 11:37:53.658698+03	3
9	14	3	2025-03-31 10:59:14.682192+03	2025-04-07 10:59:14.682192+03	returned	2025-03-31 10:59:14.682206+03	2025-03-31 11:15:08.090966+03	3
3	1	1	2024-10-26 20:13:48.288232+03	2024-11-02 20:13:48.288232+03	returned	2024-10-26 20:13:48.288261+03	2025-04-06 20:53:36.93212+03	1
10	15	11	2025-04-14 09:25:58.45427+03	2025-04-21 09:25:58.45427+03	returned	2025-04-14 09:25:58.454383+03	2025-04-14 09:29:45.827164+03	6
11	16	11	2025-04-14 09:30:37.931824+03	2025-04-21 09:30:37.931824+03	returned	2025-04-14 09:30:37.931849+03	2025-04-14 09:37:17.815555+03	6
12	17	11	2025-04-14 09:43:54.074714+03	2025-04-21 09:43:54.074714+03	returned	2025-04-14 09:43:54.074741+03	2025-04-29 13:25:23.497535+03	6
\.


--
-- Data for Name: reviews; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.reviews (id, user_id, book_id, rating, message, created_at, updated_at) FROM stdin;
1	1	1	4	wow	2024-10-29 23:31:09.697841+03	2024-10-29 23:31:09.697841+03
5	14	3	5	Хорошая книга	2024-11-25 11:24:33.07548+03	2024-11-25 11:24:33.07548+03
6	1	20	4	Хорошая книга	2025-05-06 14:22:16.661147+03	2025-05-06 14:22:16.661147+03
8	1	11	3	Плохая книга	2025-05-07 18:39:46.99658+03	2025-05-07 18:39:46.99658+03
\.


--
-- Data for Name: unique_codes; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.unique_codes (id, code, book_id, is_available) FROM stdin;
4	4	1	t
2	2	2	t
5	5	2	t
3	3	3	t
1	1	1	t
7	1214	1	t
6	2131	11	t
\.


--
-- Data for Name: users; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.users (id, login, full_name, password, role, email, phone_number, passport_number, created_at, updated_at, reset_password_token, reset_token_expiry) FROM stdin;
4	admin	admin admin amdin	$2a$10$wOr4qEaNwr3WFl1cbsgDR.RsY/HTFaPZWfEXuzy1.bD9dz8otoxk2	admin	admin@gmail.com	+7 (123) 123-12-41	1241241412	2024-10-26 00:09:39.371052+03	2024-10-26 00:09:39.371052+03	\N	\N
5	libr	libr	$2a$10$v5R5P/wwNIZTAgV9hADn1.hJcZGeqHbYIdwWBF5U6LCOTlCv49q/6	librarian	libr@gmail.com	+7 (112) 312-41-24	1241241221	2024-10-30 18:40:55.884689+03	2024-10-30 18:40:55.884689+03	\N	\N
16	newUser3	newUser3 newUser3 newUser4	$2a$10$uUamuSnBuDbWEAYvVCBbMO/zLArABKqLNgUtCszq7dzM4eR2zPBIy	user	newUser3@gmail.com	+7 (925) 777-44-77	1111234212	2025-04-14 09:30:19.847048+03	2025-04-14 09:38:03.118191+03	\N	\N
6	user2	user2 user user	$2a$10$rUEbSPoQBvPCVLupZDOqm.eDBdTb3c9uKbG2ZJr7vPG3AhPXUjc2m	librarian	user@mail.com	+7 (123) 123-13-31	1231241242	2024-11-07 14:17:36.709479+03	2024-11-08 12:10:33.353623+03	\N	\N
7	user3	user user user	$2a$10$IAyAwede1a.LbgSqYORQUeDKREVTmIcsOJlUO3Yl03KAiRuFqNyFS	user	user3@gmail.com	+7 (123) 123-12-31	1242144555	2024-11-22 12:27:12.171752+03	2024-11-22 12:27:12.171752+03	\N	\N
17	newUser4	newUser4 newUser4 newUser4	$2a$10$grOogFyeTuwWvdoqTG4mT.KcGIkN/8mg8MDyU/3TP6sneYwaEASOe	user	newUser4@gmail.com	+7 (925) 777-45-95	1234125192	2025-04-14 09:43:36.00722+03	2025-04-14 09:43:36.00722+03	\N	\N
14	newUser	new new new	$2a$10$QwioXJYThIY751.s9nHEjOcEw79biducxgZG32lbqbtTXzTcvtOMO	user	newUser@gmail.com	+7 (925) 777-44-95	1231241214	2024-11-25 11:24:07.552186+03	2025-04-06 20:54:19.568534+03	\N	\N
15	newUser2	new2 new2 new2	$2a$10$vwu4UYZ0OU8AUSsGIlJDt.G2nrxfRhVlV5OVu2nD1f/g02d2nq.52	user	newUser2@gmail.com	+7 (925) 717-44-90	1111123412	2025-04-14 09:25:33.86955+03	2025-04-14 09:25:33.86955+03	\N	\N
1	user	Александр Азаров	$2a$10$1hDJ9UB.qBSa3m2u3mnqbek.U9mxiUjxHl758Z740L2oJWbvEtNti	user	user@gmail.com	+79851594884	12321230	2024-10-22 17:35:51.22645+03	2024-10-22 17:35:51.22645+03	\N	\N
968	skylang	skylang skylang skylang	$2a$10$d8Z/wwBx6oquS8oQKeHYOuoLKrYMxbT9Wpd4bf0FDkEBNVd/wNQnG	user	skylang@inbox.ru	+7 (123) 125-12-21	1242232241	2025-05-02 22:19:38.58253+03	2025-05-12 18:06:30.528952+03	N___PdvJ0ZZb-FxQp5p8sjhBEJC3zIeQd9zY0n0qR3k=	2025-05-13 18:06:30+03
\.


--
-- Name: audiobook_files_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.audiobook_files_id_seq', 3, true);


--
-- Name: authors_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.authors_id_seq', 2289, true);


--
-- Name: books_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.books_id_seq', 22, true);


--
-- Name: chats_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.chats_id_seq', 10, true);


--
-- Name: favorites_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.favorites_id_seq', 34, true);


--
-- Name: genres_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.genres_id_seq', 7, true);


--
-- Name: logs_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.logs_id_seq', 2448, true);


--
-- Name: messages_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.messages_id_seq', 17, true);


--
-- Name: notifications_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.notifications_id_seq', 1, false);


--
-- Name: publishers_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.publishers_id_seq', 5, true);


--
-- Name: reservations_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.reservations_id_seq', 12, true);


--
-- Name: reviews_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.reviews_id_seq', 8, true);


--
-- Name: unique_codes_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.unique_codes_id_seq', 7, true);


--
-- Name: users_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.users_id_seq', 968, true);


--
-- Name: audiobook_files audiobook_files_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.audiobook_files
    ADD CONSTRAINT audiobook_files_pkey PRIMARY KEY (id);


--
-- Name: authors authors_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.authors
    ADD CONSTRAINT authors_pkey PRIMARY KEY (id);


--
-- Name: books books_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.books
    ADD CONSTRAINT books_pkey PRIMARY KEY (id);


--
-- Name: chats chats_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.chats
    ADD CONSTRAINT chats_pkey PRIMARY KEY (id);


--
-- Name: favorites favorites_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.favorites
    ADD CONSTRAINT favorites_pkey PRIMARY KEY (id);


--
-- Name: genres genres_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.genres
    ADD CONSTRAINT genres_pkey PRIMARY KEY (id);


--
-- Name: logs logs_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.logs
    ADD CONSTRAINT logs_pkey PRIMARY KEY (id);


--
-- Name: messages messages_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.messages
    ADD CONSTRAINT messages_pkey PRIMARY KEY (id);


--
-- Name: notifications notifications_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.notifications
    ADD CONSTRAINT notifications_pkey PRIMARY KEY (id);


--
-- Name: publishers publishers_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.publishers
    ADD CONSTRAINT publishers_pkey PRIMARY KEY (id);


--
-- Name: reservations reservations_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.reservations
    ADD CONSTRAINT reservations_pkey PRIMARY KEY (id);


--
-- Name: reviews reviews_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.reviews
    ADD CONSTRAINT reviews_pkey PRIMARY KEY (id);


--
-- Name: users uni_users_email; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT uni_users_email UNIQUE (email);


--
-- Name: users uni_users_passport_number; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT uni_users_passport_number UNIQUE (passport_number);


--
-- Name: users uni_users_phone_number; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT uni_users_phone_number UNIQUE (phone_number);


--
-- Name: unique_codes unique_codes_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.unique_codes
    ADD CONSTRAINT unique_codes_pkey PRIMARY KEY (id);


--
-- Name: users users_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_pkey PRIMARY KEY (id);


--
-- Name: idx_book_title; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_book_title ON public.books USING gin (to_tsvector('simple'::regconfig, title));


--
-- Name: idx_books_title; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_books_title ON public.books USING gin (to_tsvector('russian'::regconfig, title));


--
-- Name: books fk_authors_books; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.books
    ADD CONSTRAINT fk_authors_books FOREIGN KEY (author_id) REFERENCES public.authors(id) ON DELETE CASCADE;


--
-- Name: audiobook_files fk_books_audiobook_files; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.audiobook_files
    ADD CONSTRAINT fk_books_audiobook_files FOREIGN KEY (book_id) REFERENCES public.books(id);


--
-- Name: favorites fk_books_favorites; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.favorites
    ADD CONSTRAINT fk_books_favorites FOREIGN KEY (book_id) REFERENCES public.books(id);


--
-- Name: reservations fk_books_reservations; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.reservations
    ADD CONSTRAINT fk_books_reservations FOREIGN KEY (book_id) REFERENCES public.books(id);


--
-- Name: reviews fk_books_reviews; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.reviews
    ADD CONSTRAINT fk_books_reviews FOREIGN KEY (book_id) REFERENCES public.books(id);


--
-- Name: unique_codes fk_books_unique_codes; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.unique_codes
    ADD CONSTRAINT fk_books_unique_codes FOREIGN KEY (book_id) REFERENCES public.books(id);


--
-- Name: books fk_genres_books; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.books
    ADD CONSTRAINT fk_genres_books FOREIGN KEY (genre_id) REFERENCES public.genres(id) ON DELETE CASCADE;


--
-- Name: books fk_publishers_books; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.books
    ADD CONSTRAINT fk_publishers_books FOREIGN KEY (publisher_id) REFERENCES public.publishers(id) ON DELETE CASCADE;


--
-- Name: reservations fk_reservations_unique_code; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.reservations
    ADD CONSTRAINT fk_reservations_unique_code FOREIGN KEY (unique_code_id) REFERENCES public.unique_codes(id);


--
-- Name: unique_codes fk_unique_codes_book; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.unique_codes
    ADD CONSTRAINT fk_unique_codes_book FOREIGN KEY (book_id) REFERENCES public.books(id);


--
-- Name: favorites fk_users_favorites; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.favorites
    ADD CONSTRAINT fk_users_favorites FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE;


--
-- Name: logs fk_users_logs; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.logs
    ADD CONSTRAINT fk_users_logs FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE;


--
-- Name: notifications fk_users_notifications; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.notifications
    ADD CONSTRAINT fk_users_notifications FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE;


--
-- Name: reservations fk_users_reservations; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.reservations
    ADD CONSTRAINT fk_users_reservations FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE;


--
-- Name: reviews fk_users_reviews; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.reviews
    ADD CONSTRAINT fk_users_reviews FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE;


--
-- PostgreSQL database dump complete
--

