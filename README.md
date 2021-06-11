# data exporter

data exporter - командна програма написана на go, що дозволяє витягти дані з сайту. Наприклад свої оцінки з "кинопоиск" у файл csv.


## Базове використання
1. Для базового використання профіль в "кинопоиск" потрібно щоб був публічний на час роботи програми (10-20 секунд).
2. У файлі `config.json` змінити значення `baseUrl`, замість `XXX` вписати ваш "кинопоиск" ідентифікатор

```
"baseUrl": "https://www.kinopoisk.ru/user/XXX/votes/list/ord/date/page/%d/"
```
3. Запустити `data-exporter.exe` (для windows) чи `data-exporter` (для linux)

Перегляди будуть записані з оцінкою 0.

## Просунуте використання

#### Запуск програми з перепризначенням файлу конфігурації
```
data-exporter.exe -config="config.json"
```

#### Запуск програми з перепризначенням baseUrl
```
data-exporter.exe -baseUrl="https://www.kinopoisk.ru/user/XXX/votes/list/ord/date/page/%d/"
```

#### Запуск програми з перепризначенням fileName
```
data-exporter.exe -fileName="mykp_votes.csv"
```

Можна комбінувати

## Додаткові можливості
У файлі `config.json` є багато значень:
- `regexpPattern` - паттерн регулярного виразу, який збирає id, nameEn, nameRu, date, vote зі сторінки кп
- `regexpPatternPages` - паттерн регулярного виразу, який знаходить сторінки. Це для визначення кількості сторінок.
- `csvHeaders` - це те, що буде в експортованому файлі в першому рядку.
- `regexpIndexes` - індекси результатів regexpPattern, які записувати і в якому порядку.
- `regexpIndexVote` - індекс результату regexpPattern, де є ваша оцінка (потрібно тільки якщо regexpConvertVoteToNumber true)
- `regexpConvertVoteToNumber` - чи конвертувати оцінку в число
- `fileName` - назва файлу з вашими оцінками, що програма створить
- `httpHeaders` - якщо не знаєте, що це, то вам це не потрібно :)

