#!/bin/bash
#
# test script of keycut
#
# usage: [<test-path>/]cjoin0_test.bash [<command-path>]
#
#            <test-path>は
#                    「現ディレクトリーからみた」本スクリプトの相対パス
#                    または本スクリプトの完全パス
#                    省略時は現ディレクトリーを仮定する
#            <command-path>は
#                    「本スクリプトのディレクトリーからみた」test対象コマンドの相対パス
#                    またはtest対象コマンドの完全パス
#                    省略時は本スクリプトと同じディレクトリーを仮定する
#                    値があるときまたは空値（""）で省略を示したときはあとにつづく<python-version>を指定できる
name=sm2 # test対象コマンドの名前
testpath=$(dirname $0) # 本スクリプト実行コマンドの先頭部($0)から本スクリプトのディレトリー名をとりだす
cd $testpath # 本スクリプトのあるディレクトリーへ移動
if test "$1" = ""; # <command-path>($1)がなければ
	then commandpath="." # test対象コマンドは現ディレクトリーにある
	else commandpath="$1" # <command-path>($1)があればtest対象コマンドは指定のディレクトリーにある
fi
com="go run ${commandpath}/${name}.go"
marumecom=marume # marumeコマンドを使用するため
tmp=/tmp/$$

ERROR_CHECK(){
	[ "$(echo ${PIPESTATUS[@]} | tr -d ' 0')" = "" ] && return
	echo $1
	echo "${name}" NG
	rm -f $tmp-*
	exit 1
}

###########################################
#TEST1

cat << FIN > $tmp-in
001 1
001 1.11
001 -2.1
002 0.0
002 1.101
FIN

cat << FIN > $tmp-out
001 0.010
002 1.101
FIN

${com} 1 1 2 2 $tmp-in 		|
${marumecom} 2.3			> $tmp-ans
diff $tmp-ans $tmp-out
[ $? -eq 0 ] ; ERROR_CHECK "TEST1 error"

###########################################
#TEST2

cat << FIN > $tmp-in
001 1
001 1.11
001 -2.1
002 0.0
002 1.101
FIN

cat << FIN > $tmp-out
1.111
FIN

${com} 0 0 2 2 $tmp-in		|
${marumecom} 1.3			> $tmp-ans
diff $tmp-ans $tmp-out
[ $? -eq 0 ] ; ERROR_CHECK "TEST2 error"

###########################################
#TEST3

# cat << FIN > $tmp-in
# 001 1
# 001 1.11
# 001 -2.1
# 002 0.0
# 002 1.101
# FIN

# cat << FIN > $tmp-out
# 001 3 0.010
# 002 2 1.101
# FIN

# ${com} +count 1 1 2 2 $tmp-in 	|
# ${marumecom} 3.3			> $tmp-ans
# diff $tmp-ans $tmp-out
# [ $? -eq 0 ] ; ERROR_CHECK "TEST3.1 error"

# cat $tmp-in			|
# ${com} +count 1 1 2 2 		|
# ${marumecom} 3.3			> $tmp-ans
# diff $tmp-ans $tmp-out
# [ $? -eq 0 ] ; ERROR_CHECK "TEST3.2 error"

###########################################
#TEST4 Support of Scientific Representation

cat << FIN > $tmp-in
-1.0e+1
-1.0e+0
FIN

cat << FIN > $tmp-out
-11
FIN

cat $tmp-in			|
${com} 0 0 1 1			|
${marumecom} 1.0 > $tmp-ans
diff $tmp-ans $tmp-out
[ $? -eq 0 ] ; ERROR_CHECK "TEST4 error"

###########################################
#TEST5 Invalid word

cat << FIN > $tmp-in
あ
FIN

cat $tmp-in	|
${com} 0 0 1 1	> $tmp-ans	2> /dev/null
#if exit status is zero, it's an error
[ $? -eq 0 ] && false || true
ERROR_CHECK "TEST5 error"

###########################################
#TEST6 a bugfix (The value disappears when the value is zero. )

cat << FIN > $tmp-in
a 1
a 3
a 2
b 0
b 0
c 5
c 6
d 0
d 2
FIN

cat << FIN > $tmp-out
a 6
b 0
c 11
d 2
FIN

cat $tmp-in	|
${com} 1 1 2 2	|
diff - $tmp-out
[ $? -eq 0 ] ; ERROR_CHECK "TEST6 error"

rm -Rf $tmp-*
echo "${name}" OK
exit 0