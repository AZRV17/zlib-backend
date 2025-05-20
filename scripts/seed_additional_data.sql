-- Скрипт для добавления дополнительных данных в базу данных ZLib
-- Скрипт сохраняет существующие данные и добавляет новые

-- Добавление новых жанров
INSERT INTO genres (name, description, created_at, updated_at)
VALUES 
  ('Фэнтези', 'Литературный жанр, основанный на использовании мифологических и сказочных мотивов', NOW(), NOW()),
  ('Научная фантастика', 'Жанр литературы, описывающий вымышленные технологии и научные открытия', NOW(), NOW()),
  ('Детектив', 'Литературный жанр, описывающий процесс исследования загадочного происшествия с целью выяснения его обстоятельств', NOW(), NOW()),
  ('Приключения', 'Жанр литературы, повествующий о путешествиях и открытиях', NOW(), NOW()),
  ('Биография', 'Описание жизни и деятельности конкретного человека', NOW(), NOW()),
  ('Историческая проза', 'Художественные произведения на исторические темы', NOW(), NOW()),
  ('Поэзия', 'Художественная литература, опирающаяся в основном на поэтическую речь', NOW(), NOW()),
  ('Психология', 'Книги о человеческом поведении и психических процессах', NOW(), NOW()),
  ('Философия', 'Литература о фундаментальных вопросах бытия и познания', NOW(), NOW()),
  ('Драма', 'Литературный жанр, основывающийся на драматических ситуациях', NOW(), NOW());

-- Добавление новых издательств
INSERT INTO publishers (name, created_at, updated_at)
VALUES 
  ('Эксмо', NOW(), NOW()),
  ('АСТ', NOW(), NOW()),
  ('Азбука', NOW(), NOW()),
  ('Альпина Паблишер', NOW(), NOW()),
  ('МИФ', NOW(), NOW()),
  ('Росмэн', NOW(), NOW()),
  ('Академический проект', NOW(), NOW()),
  ('Фантом Пресс', NOW(), NOW()),
  ('Corpus', NOW(), NOW()),
  ('Лениздат', NOW(), NOW());

-- Добавление новых авторов
INSERT INTO authors (name, lastname, biography, birthdate, created_at, updated_at)
VALUES 
  ('Александр', 'Пушкин', 'Великий русский поэт, драматург и прозаик', '1799-06-06', NOW(), NOW()),
  ('Фёдор', 'Тютчев', 'Русский поэт, дипломат', '1803-12-05', NOW(), NOW()),
  ('Иван', 'Тургенев', 'Русский писатель-реалист, поэт', '1818-11-09', NOW(), NOW()),
  ('Антон', 'Чехов', 'Русский писатель, прозаик, драматург', '1860-01-29', NOW(), NOW()),
  ('Михаил', 'Булгаков', 'Русский писатель, драматург, театральный режиссёр и актёр', '1891-05-15', NOW(), NOW()),
  ('Джордж', 'Оруэлл', 'Английский писатель и публицист', '1903-06-25', NOW(), NOW()),
  ('Эрнест', 'Хемингуэй', 'Американский писатель, журналист', '1899-07-21', NOW(), NOW()),
  ('Габриэль', 'Гарсиа Маркес', 'Колумбийский писатель-прозаик', '1927-03-06', NOW(), NOW()),
  ('Джоан', 'Роулинг', 'Британская писательница, сценаристка и продюсер', '1965-07-31', NOW(), NOW()),
  ('Стивен', 'Кинг', 'Американский писатель, работающий в жанрах ужасов', '1947-09-21', NOW(), NOW()),
  ('Рэй', 'Брэдбери', 'Американский писатель-фантаст', '1920-08-22', NOW(), NOW()),
  ('Айзек', 'Азимов', 'Американский писатель-фантаст', '1920-01-02', NOW(), NOW()),
  ('Агата', 'Кристи', 'Английская писательница и драматург', '1890-09-15', NOW(), NOW()),
  ('Артур', 'Конан Дойл', 'Английский писатель, автор Шерлока Холмса', '1859-05-22', NOW(), NOW()),
  ('Оскар', 'Уайльд', 'Ирландский поэт, драматург, писатель, эссеист', '1854-10-16', NOW(), NOW()),
  ('Марк', 'Твен', 'Американский писатель, журналист и общественный деятель', '1835-11-30', NOW(), NOW()),
  ('Харуки', 'Мураками', 'Японский писатель и переводчик', '1949-01-12', NOW(), NOW()),
  ('Виктор', 'Гюго', 'Французский писатель, поэт, драматург', '1802-02-26', NOW(), NOW()),
  ('Фридрих', 'Ницше', 'Немецкий философ, культурный критик, поэт', '1844-10-15', NOW(), NOW()),
  ('Альбер', 'Камю', 'Французский писатель и философ', '1913-11-07', NOW(), NOW());

-- Добавление новых книг
INSERT INTO books (title, author_id, genre_id, description, publisher_id, isbn, year_of_publication, rating, is_available, created_at, updated_at)
VALUES 
  ('Евгений Онегин', (SELECT id FROM authors WHERE name = 'Александр' AND lastname = 'Пушкин' LIMIT 1), 
   (SELECT id FROM genres WHERE name = 'Роман' LIMIT 1), 
   'Роман в стихах, повествующий о жизни молодого дворянина Евгения Онегина', 
   (SELECT id FROM publishers WHERE name = 'Лениздат' LIMIT 1), 
   9781234567801, '1833-01-01', 4.9, true, NOW(), NOW()),
   
  ('Мастер и Маргарита', (SELECT id FROM authors WHERE name = 'Михаил' AND lastname = 'Булгаков' LIMIT 1), 
   (SELECT id FROM genres WHERE name = 'Роман' LIMIT 1), 
   'Роман о дьяволе, посетившем Москву, и о мастере, написавшем роман о Понтии Пилате', 
   (SELECT id FROM publishers WHERE name = 'АСТ' LIMIT 1), 
   9781234567802, '1966-01-01', 4.8, true, NOW(), NOW()),
   
  ('1984', (SELECT id FROM authors WHERE name = 'Джордж' AND lastname = 'Оруэлл' LIMIT 1), 
   (SELECT id FROM genres WHERE name = 'Анти утопия' LIMIT 1), 
   'Роман-антиутопия о тоталитарном обществе', 
   (SELECT id FROM publishers WHERE name = 'Эксмо' LIMIT 1), 
   9781234567803, '1949-01-01', 4.7, true, NOW(), NOW()),
   
  ('Старик и море', (SELECT id FROM authors WHERE name = 'Эрнест' AND lastname = 'Хемингуэй' LIMIT 1), 
   (SELECT id FROM genres WHERE name = 'Роман' LIMIT 1), 
   'Повесть о кубинском рыбаке, вступившем в единоборство с гигантской рыбой', 
   (SELECT id FROM publishers WHERE name = 'АСТ' LIMIT 1), 
   9781234567804, '1952-01-01', 4.5, true, NOW(), NOW()),
   
  ('Сто лет одиночества', (SELECT id FROM authors WHERE name = 'Габриэль' AND lastname = 'Гарсиа Маркес' LIMIT 1), 
   (SELECT id FROM genres WHERE name = 'Роман' LIMIT 1), 
   'Роман о семье Буэндиа на протяжении нескольких поколений', 
   (SELECT id FROM publishers WHERE name = 'Эксмо' LIMIT 1), 
   9781234567805, '1967-01-01', 4.6, true, NOW(), NOW()),
   
  ('Гарри Поттер и философский камень', (SELECT id FROM authors WHERE name = 'Джоан' AND lastname = 'Роулинг' LIMIT 1), 
   (SELECT id FROM genres WHERE name = 'Фэнтези' LIMIT 1), 
   'Первая книга о приключениях юного волшебника Гарри Поттера', 
   (SELECT id FROM publishers WHERE name = 'Росмэн' LIMIT 1), 
   9781234567806, '1997-01-01', 4.7, true, NOW(), NOW()),
   
  ('Сияние', (SELECT id FROM authors WHERE name = 'Стивен' AND lastname = 'Кинг' LIMIT 1), 
   (SELECT id FROM genres WHERE name = 'Роман' LIMIT 1), 
   'Роман ужасов о писателе Джеке Торрансе, его жене и сыне, оказавшихся в отеле', 
   (SELECT id FROM publishers WHERE name = 'АСТ' LIMIT 1), 
   9781234567807, '1977-01-01', 4.5, true, NOW(), NOW()),
   
  ('451 градус по Фаренгейту', (SELECT id FROM authors WHERE name = 'Рэй' AND lastname = 'Брэдбери' LIMIT 1), 
   (SELECT id FROM genres WHERE name = 'Научная фантастика' LIMIT 1), 
   'Роман-антиутопия о обществе, где книги находятся под запретом', 
   (SELECT id FROM publishers WHERE name = 'Эксмо' LIMIT 1), 
   9781234567808, '1953-01-01', 4.6, true, NOW(), NOW()),
   
  ('Я, робот', (SELECT id FROM authors WHERE name = 'Айзек' AND lastname = 'Азимов' LIMIT 1), 
   (SELECT id FROM genres WHERE name = 'Научная фантастика' LIMIT 1), 
   'Сборник рассказов о роботах, объединённых общим сюжетом', 
   (SELECT id FROM publishers WHERE name = 'Эксмо' LIMIT 1), 
   9781234567809, '1950-01-01', 4.5, true, NOW(), NOW()),
   
  ('Убийство в Восточном экспрессе', (SELECT id FROM authors WHERE name = 'Агата' AND lastname = 'Кристи' LIMIT 1), 
   (SELECT id FROM genres WHERE name = 'Детектив' LIMIT 1), 
   'Детективный роман об убийстве в поезде', 
   (SELECT id FROM publishers WHERE name = 'Эксмо' LIMIT 1), 
   9781234567810, '1934-01-01', 4.6, true, NOW(), NOW()),
   
  ('Приключения Шерлока Холмса', (SELECT id FROM authors WHERE name = 'Артур' AND lastname = 'Конан Дойл' LIMIT 1), 
   (SELECT id FROM genres WHERE name = 'Детектив' LIMIT 1), 
   'Сборник рассказов о знаменитом сыщике', 
   (SELECT id FROM publishers WHERE name = 'АСТ' LIMIT 1), 
   9781234567811, '1892-01-01', 4.7, true, NOW(), NOW()),
   
  ('Портрет Дориана Грея', (SELECT id FROM authors WHERE name = 'Оскар' AND lastname = 'Уайльд' LIMIT 1), 
   (SELECT id FROM genres WHERE name = 'Роман' LIMIT 1), 
   'Роман о портрете, стареющем вместо своего хозяина', 
   (SELECT id FROM publishers WHERE name = 'Азбука' LIMIT 1), 
   9781234567812, '1890-01-01', 4.5, true, NOW(), NOW()),
   
  ('Приключения Тома Сойера', (SELECT id FROM authors WHERE name = 'Марк' AND lastname = 'Твен' LIMIT 1), 
   (SELECT id FROM genres WHERE name = 'Приключения' LIMIT 1), 
   'Роман о приключениях мальчика, живущего на берегу реки Миссисипи', 
   (SELECT id FROM publishers WHERE name = 'Эксмо' LIMIT 1), 
   9781234567813, '1876-01-01', 4.6, true, NOW(), NOW()),
   
  ('Норвежский лес', (SELECT id FROM authors WHERE name = 'Харуки' AND lastname = 'Мураками' LIMIT 1), 
   (SELECT id FROM genres WHERE name = 'Роман' LIMIT 1), 
   'Роман о любви и потере в Японии конца 1960-х', 
   (SELECT id FROM publishers WHERE name = 'Эксмо' LIMIT 1), 
   9781234567814, '1987-01-01', 4.5, true, NOW(), NOW()),
   
  ('Отверженные', (SELECT id FROM authors WHERE name = 'Виктор' AND lastname = 'Гюго' LIMIT 1), 
   (SELECT id FROM genres WHERE name = 'Историческая проза' LIMIT 1), 
   'Роман о бывшем заключенном Жане Вальжане', 
   (SELECT id FROM publishers WHERE name = 'Эксмо' LIMIT 1), 
   9781234567815, '1862-01-01', 4.7, true, NOW(), NOW()),
   
  ('Так говорил Заратустра', (SELECT id FROM authors WHERE name = 'Фридрих' AND lastname = 'Ницше' LIMIT 1), 
   (SELECT id FROM genres WHERE name = 'Философия' LIMIT 1), 
   'Философское произведение о сверхчеловеке', 
   (SELECT id FROM publishers WHERE name = 'Академический проект' LIMIT 1), 
   9781234567816, '1883-01-01', 4.4, true, NOW(), NOW()),
   
  ('Посторонний', (SELECT id FROM authors WHERE name = 'Альбер' AND lastname = 'Камю' LIMIT 1), 
   (SELECT id FROM genres WHERE name = 'Роман' LIMIT 1), 
   'Роман о человеке, совершившем убийство', 
   (SELECT id FROM publishers WHERE name = 'Азбука' LIMIT 1), 
   9781234567817, '1942-01-01', 4.5, true, NOW(), NOW()),
   
  ('Вишневый сад', (SELECT id FROM authors WHERE name = 'Антон' AND lastname = 'Чехов' LIMIT 1), 
   (SELECT id FROM genres WHERE name = 'Драма' LIMIT 1), 
   'Пьеса о судьбе вишневого сада, который принадлежит разорившейся помещице', 
   (SELECT id FROM publishers WHERE name = 'АСТ' LIMIT 1), 
   9781234567818, '1903-01-01', 4.6, true, NOW(), NOW()),
   
  ('Отцы и дети', (SELECT id FROM authors WHERE name = 'Иван' AND lastname = 'Тургенев' LIMIT 1), 
   (SELECT id FROM genres WHERE name = 'Роман' LIMIT 1), 
   'Роман о взаимоотношениях двух поколений', 
   (SELECT id FROM publishers WHERE name = 'АСТ' LIMIT 1), 
   9781234567819, '1862-01-01', 4.5, true, NOW(), NOW()),
   
  ('Стихотворения', (SELECT id FROM authors WHERE name = 'Фёдор' AND lastname = 'Тютчев' LIMIT 1), 
   (SELECT id FROM genres WHERE name = 'Поэзия' LIMIT 1), 
   'Сборник стихотворений русского поэта', 
   (SELECT id FROM publishers WHERE name = 'Лениздат' LIMIT 1), 
   9781234567820, '1850-01-01', 4.4, true, NOW(), NOW());

-- Добавление уникальных кодов для новых книг
INSERT INTO unique_codes (code, book_id, is_available)
SELECT 1000 + b.id, b.id, true
FROM books b
WHERE b.id NOT IN (SELECT DISTINCT book_id FROM unique_codes);

INSERT INTO unique_codes (code, book_id, is_available)
SELECT 2000 + b.id, b.id, true
FROM books b
WHERE b.id NOT IN (SELECT DISTINCT book_id FROM unique_codes WHERE code >= 2000);

-- Теперь у каждой книги есть как минимум по 2 уникальных кода
