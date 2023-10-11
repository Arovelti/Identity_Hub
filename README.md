# Task 
Необходимо написать веб сервис который занимается хранением профилей пользователей и их авторизацией.

Профиль имеет набор полей:
1. id (uuid, unique)
2. email
3. username (unique)
4. password
5. admin (bool)

Сервис должeн иметь набор ручек (gRPC):
Создание пользователя
Выдача списка пользователей
Выдача юзера по id
Изменение и удаление профиля

Сервис использует basic access authentication (https://en.wikipedia.org/wiki/Basic_access_authentication)

Просматривать профили могут все зарегистрированные пользователи.
Админ может создавать, изменять и удалять профили.

Для хранения данных профилей необходимо реализовать примитивную in memory базу данных, не используя сторонние решения.

Задание не сложное, мы очень хотим посмотреть на то, как ты пишешь код.

Времени 1 рабочая неделя, если получится раньше - отлично!

Задачу можно выполнить с разной степенью глубины. Если понимаешь, что не успеешь сделать качественнее, все важные проблемы и допущения нужно указать в файле README

# Remarks
- Вместо 5. admin (bool) я бы использовал enum или сущность для ролей.
- Я бы выставил на обсуждение структуру профиля. Поскольку она хранит пароли, то я бы предложил создать отдельную стуктуру для хранения хэшированных данных пароля и имени пользователя, к примеру "Credentials", которые были бы привязаны к конкретному профилю и юзеру, и соответствующая логика бы проверяла sensitive data. 
- Мне очень не нравится цикл for в реализации методов в репозитории. Здесь есть небольшой трейд-ин между тем, чтобы сохранять уникальность id и имени юзера, и тем, как мы будем работать с данным. По ходу поиска решения пришел к тому, что начинаю создать собственную im-memory базу данных, поэтому вернулся к пройтейшему варианту со всеми его погрешностями. Мне понравился подобный вариант, который уже есть на github: https://github.com/hashicorp/go-memdb 
- В UpdateProfile можно использовать более сложную логику, с доп проверкамиБ к примеру, совпадают ли запрашиваемые id запроса и профайла. Вариаций много, зависит от требований бизнеса. 
- В пакете playground собран dirty code для разных экспириментов (туда лучше не залазить)

# TODO:
- пройтись по коду и облагородить 
- улучшить работу с ошибками 
- разбить main() 
- добавить работу с envs 
- добавить ci/cd 
- добавить метрики 
- добавить клиентскую часть
- развить устройство in-memory db

# gRPC
$ go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28
$ go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2

$ go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
$ go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@latest
$ go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest

$ go get -u  google.golang.org/genproto/googleapis/api@latest
$ go get -u  google.golang.org/genproto/googleapis/rpc@latest

$ export PATH="$PATH:$(go env GOPATH)/bin"

$ git submodule add https://github.com/googleapis/googleapis
cd ~/go/pkg/mod/github.com/grpc-ecosystem/grpc-gateway

Generate:
    protoc -I . \
    --go_out ./api --go_opt paths=source_relative \
    --go-grpc_out ./api --go-grpc_opt paths=source_relative \
    proto/profile.proto

    protoc -I . --grpc-gateway_out ./api \
    --grpc-gateway_opt logtostderr=true \
    --grpc-gateway_opt paths=source_relative \
    --grpc-gateway_opt generate_unbound_methods=true \
    proto/profile.proto

