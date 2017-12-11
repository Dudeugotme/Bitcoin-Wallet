# Bitcoin-Wallet

Модифицированный SPV Wallet для Bitcoin, использует наработки от Lighting Network, OpenBazaar  и btcd проектов.

1. Установить версию go. https://golang.org/doc/install
2. Настроить GOPATH и для workspace проекта. Переходим в папку <workspace>/src
3. git clone https://github.com/OpenDomido/Bitcoin-Wallet.git или go get github.com/OpenDomido/Bitcoin-Wallet
4. В папке workspace (там где GoPath) или в папке github переходим cd Bitcoin-Wallet/cmd/bitcoin
5. go get установит все зависимости
6. go build -  соберет проект
7. go install - установить проект  
9. Команда bitcoin запустить проект и grpc сервис.

В папке API лежит описание сервисов grpc, файл API.proto

# API
  1. rpc Stop (Empty) returns (Empty) {} - остановить сервис
  2. rpc CurrentAddress (KeySelection) returns (Address) {} - получить текущий адрес. Адрес бывает двух видов внешний для получения денег и внутренний адрес. Для получения средств нужно указывать внешний - external в качестве параметра или 0
  enum KeyPurpose {
    INTERNAL = 0;
    EXTERNAL = 1;
  }
  
  3. Получить новый адрес. Параметры аналогично External и Internal. rpc NewAddress (KeySelection) returns (Address) {} - основная функция для получения адреса
  
  4. Получить высоту цепочку (как правило используется для получения количества блоков в цепочке по отношению к текущему адреса (факт подтверждения)
  rpc ChainTip (Empty) returns (Height) {}
  
  5. Вернуть баланс кошелька
  rpc Balance (Empty) returns (Balances) {}
  6. Получить приватный master-ключ кошелька - требуется для бекапа
  rpc MasterPrivateKey (Empty) returns (Key) {}
  
  7. Получить публичный мастер-ключ кошелька - SID. 
  rpc MasterPublicKey (Empty) returns (Key) {}
  8. Проверка наличия у адреса приватного адреса
  rpc HasKey (Address) returns (BoolResponse) {}
  
  9. Получить параметры конфигуарции
  rpc Params (Empty) returns (NetParams) {}
  10. Получить список транзакций (входящие и исходящие из кошелька)
  rpc Transactions (Empty) returns (TransactionList) {}
  11. Получить конкретную транзакцию
  rpc GetTransaction (Txid) returns (Tx) {}
  12. Получить размер комиссии оптимальный в данный момент в зависимости от размера транзакции в байтах
  rpc GetFeePerByte (FeeLevelSelection) returns (FeePerByte) {}
  13. Информация относительно сдачи в транзакции
  rpc Spend (SpendInfo) returns (Txid) {}
  14. Информация о комиссии в транзакции
  rpc BumpFee (Txid) returns (Txid) {}
  15. Список пиров
  rpc Peers (Empty) returns (PeerList) {}
  16. Добавить отслеживающий скрипт на адрес
  rpc AddWatchedScript (Address) returns (Empty) {}
  17. Получить количество подтверждений к транзакции
  rpc GetConfirmations (Txid) returns (Confirmations) {}
  18. Очистить адрес
  rpc SweepAddress (SweepInfo) returns (Txid) {}
  19. Синхронизировать блокчейн
  rpc ReSyncBlockchain (Height) returns (Empty) {}
  20. Создать мультиподпись
  rpc CreateMultisigSignature (CreateMultisigInfo) returns (SignatureList) {}
  21. Добавить мультиподпись
  rpc Multisign (MultisignInfo) returns (RawTx) {}
  22. Примерная стоимость комисии
  rpc EstimateFee (EstimateFeeData) returns (Fee) {}
  23. Получить приватный ключ
  rpc GetKey (Address) returns (Key) {}
  24. Список ключей
  rpc ListKeys (Empty) returns (Keys) {}
  25. Список адресов
  rpc ListAddresses (Empty) returns (Addresses) {}

