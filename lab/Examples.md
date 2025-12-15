# Originalniye (1-4 indexes)

./bin/bc -add \
 -course 5 -group "5.507M" -name "Иванов Иван Иванович" \
 -zachetka "202434" -subject "Математика" -grade 5

./bin/bc -add \
 -course 5 -group "5.507M" -name "Петров Петр Петрович" \
 -zachetka "202435" -subject "Физика" -grade 4

./bin/bc -add \
 -course 5 -group "5.507M" -name "Сидоров Сидор Сидорович" \
 -zachetka "202436" -subject "Химия" -grade 2

./bin/bc -add \
 -course 5 -group "5.507M" -name "Васильев Василий Васильевич" \
 -zachetka "202437" -subject "Биология" -grade 3

# SoftFork (5+)

# Dolshe mainiyatsa po vremeni i oni 000

./bin/bc -add \
 -course 5 -group "5.507M" -name "Александров Александр Александрович" \
 -zachetka "202438" -subject "История" -grade 2 \

./bin/bc -add \
 -course 5 -group "5.507M" -name "Александров Александр Александрович" \
 -zachetka "202438" -subject "История" -grade 2 \
 -teacher "Смирнова О.П."

./bin/bc -add \
 -course 5 -group "5.507M" -name "Михайлов Михаил Михайлович" \
 -zachetka "202439" -subject "География" -grade 4 \
 -teacher "Иванов П.С."

# MOZHNO PROVERIT VALIDNOST (eto uzhe 6 block)

./bin/bc -add \
 -course 5 -group "5.507M" -name "Дмитриев Дмитрий Дмитриевич" \
 -zachetka "202440" -subject "Информатика" -grade 5 \
 -teacher "Козлов А.В."

./bin/bc -add \
 -course 5 -group "5.507M" -name "Николаев Николай Николаевич" \
 -zachetka "202441" -subject "Английский" -grade 4 \
 -teacher "Петрова Е.И."

./bin/bc -add \
 -course 5 -group "5.507M" -name "Андреев Андрей Андреевич" \
 -zachetka "202442" -subject "Литература" -grade 5 \
 -teacher "Сидорова М.А."

# HARD FORK (CSV)

./bin/bc -csv -add \
 -course 2 -group "5.507M" -name "Акакиев Акакий Акакиевич" \
 -zachetka "44444" -subject "Биология" -grade 2 \
 -teacher "Колымагин A.Д."

./bin/bc -add \
 -course 4 -group "5.507M" -name "Дружок Денис Денисович" \
 -zachetka "44444" -subject "История" -grade 4 \
 -teacher "Будейко В.З."

./bin/bc -csv -add \
 -course 5 -group "5.507M" -name "Сергеев Сергей Сергеевич" \
 -zachetka "202443" -subject "Физкультура" -grade 5 \
 -teacher "Волков К.Л."
