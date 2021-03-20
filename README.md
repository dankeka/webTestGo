# Runok

Просто учусь писать сайты на Go
Вот задумка моего сайта:
Сайт, где будут публиковать предложения о продажи, который будет включать в себя:
1. регестрацияю/авторизацию
2. публикацию/просмотор/удалине и тп товара
3. чат для продавца/покупателя  
и т. п.

Надеюсь у меня получится это сделать)

## Структура папок на сайте
### cmd
Папка, в которой хранится главный main.go 

### pkg
Папка, в которой хранится основная логика сайта.
Включает в себя папки:

  #### 1. handler
  Папка, в которой хранить файл handler.go, отвечающий за обработку роутеров(адресов сайта)

  #### 2. views
  Папка, в которой хранится логика частей сайта(главная страница, регестрация, авторизация и т. п.)

### types
Папка, в которой хранится большенство структур
dbtypes.go - структуры для таблий db,
types.go - структуры для [views](#2.-views)

### web
Папка, в которой хранятся [pugs](#1.-pugs), [static](#2.-static), [templates](#3.-templates)

  #### 1. pugs
  Как понятно по названию, папка, в которой хранятся файлы pug для [templates](#3.-templates)

  #### 2. static
  Папка, в которой хранятся статичные файлы(css, js, изображения)

  #### 3. templates
  Папка, в которой хранятся HTML шаблоны для сайта