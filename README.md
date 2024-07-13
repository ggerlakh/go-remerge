# go-remerge
![graph demo](https://github.com/ggerlakh/go-remerge/blob/master/images/first_neo4j_graph.png)
В данном репозиторий находится исходный код консольной утилиты `go-remerge`, которое было написано на языке [Golang](https://go.dev/), для создания графов зависимостей исходного кода с поддержкой выгрузки созданных графов в нереляционные графовые СУБД [Neo4j](https://neo4j.com/) и [ArangoDB](https://arangodb.com/).<br>
Утилита в ходе ананлиза исходного кода создает 5 графов:
* `filesystem` - граф иерархии файловой системы анализируемого проекта.  
  Вершины данного графа состоят и из всех файлов и директорий, которые есть в проекте.
  Ребра данного графа формируются на основе вложенности файлов и директорий в другие директории.
* `file_dependencу` - граф зависимостей между файлами и модулями с исходным кодом анализируемого проекта.  
  Вершины данного графа состоят из файлов с исходным кодом анализируемого проекта и внешних зависимостей, которые также могут импортироваться в файлах с исходным кодом.  
  Ребра данного графа формируются на основе импортов между файлами и модулями.
* `entity_dependency` - граф зависимостей между всеми классами (структурами) для заданного проекта.  
  Вершины данного графа состоят из всех классы или структур, которые извлекаются при помощи парсера (своего для соответвующего ЯП) из файлов с исходным кодом.
  Ребра данного графа формируются на основе импортов между всеми классами и структурами из анализируемого исходного кода. 
* `entity_inheritance` - граф зависимостей отражающий иерархию наследования между всеми классам (структурами) анализируемого проекта.  
  Вершины данного графа состоят из всех классы или структур, которые извлекаются при помощи парсера (своего для соответвующего ЯП) из файлов с исходным кодом.
  Ребра данного графа формируются на основе отношения наследования между всеми классами и структурами из анализируемого исходного кода. 
* `entity_complete` - граф консолидирующий всю информацию о зависимостях между классами (структурами), который представляет собой объединение графов `entity_dependency` и `entity_inheritance`.  
  Вершины данного графа состоят из объединения вершин графов `entity_dependency` и `entity_inheritance`.
  Ребра данного графа состоят из объединения ребер графов `entity_dependency` и `entity_inheritance`.

В данный момент подерживается анализ исходного кода который написан на языках:
* [Python](https://www.python.org/)
* [Golang](https://go.dev/)

Для настройики параметров анализа исходного кода и построения графов зависимостей, при запуске утилиты указывается конфигурационный файл в формате [YAML](https://yaml.org/).
<details>
<summary>Описание формата конфигурационного файла</summary>

* `project_name:` - обязательное поле в значении которого указывается название сканируемого проекта
* `analysis_name:` - обязательное поле в значении которого указывается название запускаемой задачи на анализ исходного кода
* `source_directory:` - обязательное поле в значении которого указывается путь до директории с анализируемым исходным кодом
* `languages:` - обязательное поле в значении которого указывается язык программирования на котором написан анализируемый исходный код (на данный момент допустимы только значения "golang" или "python")
* `extensions:` - обязательное поле, представляющее собой массив строк, в значении которого указываются допустимые расширения анализируемых файлов с исходным кодом для соответствующего ЯП 
* `ignore_directories:` - необязательное поле, представляющее собой массив строк, в значении которого указываются названия директорий, которые нужно пропустить в ходе анализа исходного кода
* `ignore_files:` - необязательное поле, представляющее собой массив строк, в значении которого указываются названия файлов, которые нужно пропустить в ходе анализа исходного кода
* `graphs:` - обязательное поле, представляющее собой список из объектов с полями `graph` и `direction`, в значении которого указываются типы создаваемых графов в результате анализа исходного кода.  
  - `graph:` - обязательное поле в котором указывается тип создаваемого графа (возможные значения: `filesystem`, `file_dependency`, `entity_dependency`, `entity_inheritance`, `entity_complete`)  
    `direction:` - обязательное поле в котором указывается направленность графа (возможные значения: `directed` (по умолчанию), `undirected`)
* `export:` - обязательное поле в значении которого указывается об экспорте о создаваемых результирующих графов. Поле имеет три необязательных вложенных поля (но обязательно указать хотя бы одно из них)
  * `as_file:` - необязательное поле, которое нужно для указания выгрузки созданных графов в виде JSON файлов.  
    Имеет обязательные вложенные поля `output_dir` и `formats` и нужно для выгрузки графов в виде JSON файлов
    * `output_dir:` - обязательное поле, в котором указывается путь до директории куда будут экспортироваться созданные графы в формате JSON
    * `formats:` - обязательное поле, представляющее собой список (возможные значения: `json`, `arango_format`, обязательно указать хотя бы одно из них), в котором указывается формат представления графов в виде JSON файлов  
      - `json` - стандартный формат представления графа который предусмотрен утилитой в виде JSON файла ([пример](https://github.com/ggerlakh/go-remerge/blob/master/configs/go-remerge_file_dependency.json))
      - `arango_format` -  представление графа в формате базы данных [ArangoDB](https://arangodb.com/) в виде JSON файла ([пример](https://github.com/ggerlakh/go-remerge/blob/master/configs/go-remerge_file_dependencyArangoFormat.json))
  * `arango:` - необязательное поле, которое нужно для указания выгрузки созданных графов в базу данных [ArangoDB](https://arangodb.com/)  
    * `username:` - обязательное поле, в котором указывается пользователь под которым нужно подключиться к базе данных
    * `password:` - обязательное поле, в котором указывается пароль для подключения под пользователем `username`
    * `endpoints:` - обязательное поле, представляющее собой список адресов URL для подключения к базе данных [ArangoDB](https://arangodb.com/) (нужно указать хотя бы одно значение в списке)
    * `database:` - обязательное поле, нзвание базы данных в ArangoDB](https://arangodb.com/) в которую нужно выгрузить созданные графы
  * `neo4j:` - необязательное поле, которое нужно для указания выгрузки созданных графов в базу данных [Neo4j](https://neo4j.com/)     
    * `username:` - обязательное поле, в котором указывается пользователь под которым нужно подключиться к базе данных
    * `password:` - обязательное поле, в котором указывается пароль для подключения под пользователем `username`  
    * `uri:` - обязательное поле, в котором укаазывется URI для подключения к базе данных [Neo4j](https://neo4j.com/)
</details>

Примеры конфигурационных файлов в данном формате представлены вот [здесь](https://github.com/ggerlakh/go-remerge/tree/master/configs).
## Структура репозитория
* `cmd` - основная директория содержащая точку входа для приложения `cmd/app/app.go`
* `internal` - директория содержащая исходных код всей внутренней логики приложения
* `tools` - вспомогательные инструсенты которые используются в исходном коде приложения
* `test` - директория с тестами приложения
* `scripts` - директория со вспомогательными скриптами
* `configs` - директория, содержащая примеры конфигурационных файлов
* `images` - директория содержазая изображения, которые используются в README.md
## Инструкция по запуску
Для запуска утилиты нужно выполнить следующие шаги:
1. Установить [Golang](https://go.dev/) не ниже версии 1.19, ([ссылка](https://go.dev/doc/install) на инструкцию по установке)  
   Проверить наличие Golang нужно версии на локальной машине можно с помощью данной команды
   ```bash
   go version
   ```
2. Склонировать репозиторий с исходным кодом и собрать утилиту:  
   ```bash
   git clone https://github.com/ggerlakh/go-remerge.git && cd go-remerge && go build -o go-remerge-app cmd/app/app.go && chmod +x go-remerge-app
   ```
3. Добавить путь до склонированного репозитория в переменную `PATH`  
   ```bash
   export PATH="${PATH}:<path_to_cloned_repo>"
   ```
4. Создать конфигурационный файл в заданном формате ([пример](https://github.com/ggerlakh/go-remerge/tree/master/configs))
5. Если требуется выгрузка графов в графовые базы данных [Neo4j](https://neo4j.com/) и [ArangoDB](https://arangodb.com/), их можно развернуть локально через docker-compose ([ссылка](https://docs.docker.com/compose/install/) на инструкции по установке), пример конфигурационного файла docker-composе есть в [репозитории](https://github.com/ggerlakh/go-remerge/blob/master/docker-compose.yml)  
   Поднять контейнеры с базами данных можно через следующую команду:
   ```bash
   docker-compose up -d
   ```
6. Запусить собранную утилиту для анализа исходного кода  
   Для вывода короткой справки о запуске утилиты можно выполнить следущую команду  
   ```bash
   go-remerge % go-remerge-app -h
   Usage: go-remerge-app -c <path> [-h] [-v] [--async]:
   -h, --help      print help information
   --async         asynchronous task execution
   -c              path to yaml config
   -v              produce verbose output
   ```
   * Через опцию `-c` указывается путь до конфигурационного файла или до директории с несколькими файлами, если нужно запустить сразу несколько задач  
   * Если запускается сразу несколько задач, их выполнение можно ускорить путем их асинхронного выполнение, задав флаг `--async`  
   * Для более подробного вывода логов о выполнении задач на анализ кода утилитой нужно задать опцию `-v`
## Пример запуска утилиты
Ниже представлен пример запуска утилиты для создания графов зависимостей
```bash
go-remerge % go-remerge-app -c configs/go-template.yml -v
Starting task from config configs/go-template.yml...
skipping a dir without errors: .git 
skipping a file without errors: .gitignore 
skipping a dir without errors: export 
skipping a dir without errors: scripts 
filesystem graph exported as JSON file in /Users/ggerlakh/go-remerge/export/go-remerge_filesystem.json
filesystem graph exported as ArangoDB formatted JSON file in /Users/ggerlakh/go-remerge/export/go-remerge_filesystemArangoFormat.json
starting export filesystem graph in ArangoDB...
graph filesystem exported in ArangoDB
starting export filesystem graph in Neo4j...
filesystem graph exported in Neo4j
file_dependency graph exported as JSON file in /Users/ggerlakh/go-remerge/export/go-remerge_file_dependency.json
file_dependency graph exported as ArangoDB formatted JSON file in /Users/ggerlakh/go-remerge/export/go-remerge_file_dependencyArangoFormat.json
starting export file_dependency graph in ArangoDB...
graph file_dependency exported in ArangoDB
starting export file_dependency graph in Neo4j...
file_dependency graph exported in Neo4j
entity_dependency graph exported as JSON file in /Users/ggerlakh/go-remerge/export/go-remerge_entity_dependency.json
entity_dependency graph exported as ArangoDB formatted JSON file in /Users/ggerlakh/go-remerge/export/go-remerge_entity_dependencyArangoFormat.json
starting export entity_dependency graph in ArangoDB...
graph entity_dependency exported in ArangoDB
starting export entity_dependency graph in Neo4j...
entity_dependency graph exported in Neo4j
entity_inheritance graph exported as JSON file in /Users/ggerlakh/go-remerge/export/go-remerge_entity_inheritance.json
entity_inheritance graph exported as ArangoDB formatted JSON file in /Users/ggerlakh/go-remerge/export/go-remerge_entity_inheritanceArangoFormat.json
starting export entity_inheritance graph in ArangoDB...
graph entity_inheritance exported in ArangoDB
starting export entity_inheritance graph in Neo4j...
entity_inheritance graph exported in Neo4j
entity_complete graph exported as JSON file in /Users/ggerlakh/go-remerge/export/go-remerge_entity_complete.json
entity_complete graph exported as ArangoDB formatted JSON file in /Users/ggerlakh/go-remerge/export/go-remerge_entity_completeArangoFormat.json
starting export entity_complete graph in ArangoDB...
graph entity_complete exported in ArangoDB
starting export entity_complete graph in Neo4j...
entity_complete graph exported in Neo4j
Task from config configs/go-template.yml was completed successfully
```