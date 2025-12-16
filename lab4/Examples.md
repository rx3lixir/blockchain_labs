# 3. Добавить первый блок с транзакциями

./bin/bc -add \
 -names "Иванов И.И.,Петров П.П.,Сидоров С.С." \
 -grades "5,4,3" \
 -courses "5,5,5" \
 -groups "5.507M,5.507M,5.507M" \
 -zachetkas "202434,202435,202436" \
 -subjects "Математика,Физика,Химия"

# 4. Добавить второй блок

./bin/bc -add \
 -names "Васильев В.В.,Николаев Н.Н." \
 -grades "5,4" \
 -courses "5,5" \
 -groups "5.507M,5.507M" \
 -zachetkas "202437,202438" \
 -subjects "Биология,История"

# 5. Проверить блокчейн

./bin/bc -validate

# 6. Посмотреть все блоки

./bin/bc -list

---

ДЕРЕВО логика

# 1. Визуализировать дерево

./bin/bc -merkle-build 1

# 2. Получить proof для каждой транзакции

./bin/bc -merkle-proof 1,0 # первая транзакция
./bin/bc -merkle-proof 1,1 # вторая транзакция
./bin/bc -merkle-proof 1,2 # третья транзакция

# 3. Проверить транзакции

./bin/bc -merkle-verify 1,0
./bin/bc -merkle-verify 1,1
./bin/bc -merkle-verify 1,2
