В качестве СУБД используется SQLite. Файл с БД называется tracker.db. В БД всего одна таблица parcel со следующими колонками:
number — номер посылки, целое число, автоинкрементное поле.
client — идентификатор клиента, целое число.
status — статус посылки, строка.
address — адрес посылки, строка.
created_at — дата и время создания посылки, строка.
