#!/bin/bash
#
# test script of cjoin0
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
name=cjoin1 # test対象コマンドの名前
testpath=$(dirname $0) # 本スクリプト実行コマンドの先頭部($0)から本スクリプトのディレトリー名をとりだす
cd $testpath # 本スクリプトのあるディレクトリーへ移動
if test "$1" = ""; # <command-path>($1)がなければ
	then commandpath="." # test対象コマンドは現ディレクトリーにある
	else commandpath="$1" # <command-path>($1)があればtest対象コマンドは指定のディレクトリーにある
fi
com="go run ${commandpath}/${name}.go"
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

cat << FIN > $tmp-tran
0000000 浜地______ 50 F 91 59 20 76 54
0000004 白土______ 40 M 58 71 20 10 6
0000003 杉山______ 26 F 30 50 71 36 30
0000001 鈴田______ 50 F 46 39 8  5  21
0000005 崎村______ 50 F 82 79 16 21 80
FIN

cat << FIN > $tmp-master
0000001 B
0000004 A
FIN

cat << FIN > $tmp-out
0000004 A 白土______ 40 M 58 71 20 10 6
0000001 B 鈴田______ 50 F 46 39 8 5 21
FIN

${com} key=1 $tmp-master $tmp-tran 	|
sed 's/  */ /g'				> $tmp-ans
diff $tmp-ans $tmp-out
[ $? -eq 0 ] ; ERROR_CHECK "TEST1 error"

###########################################
#TEST2

cat << FIN > $tmp-tran
0000000 浜地______ 50 F 91 59 20 76 54
0000004 白土______ 40 M 58 71 20 10 6
0000005 崎村______ 50 F 82 79 16 21 80
0000003 杉山______ 26 F 30 50 71 36 30
0000001 鈴田______ 50 F 46 39 8  5  21
FIN

cat << FIN > $tmp-master
0000001 A
0000004 B
FIN

cat << FIN > $tmp-out
0000004 B 白土______ 40 M 58 71 20 10 6
0000001 A 鈴田______ 50 F 46 39 8 5 21
FIN

cat << FIN > $tmp-ng
0000000 浜地______ 50 F 91 59 20 76 54
0000005 崎村______ 50 F 82 79 16 21 80
0000003 杉山______ 26 F 30 50 71 36 30
FIN

${com} +ng key=1 $tmp-master $tmp-tran 2> $tmp-ans2	|
sed 's/  */ /g'	> $tmp-ans
diff $tmp-ans $tmp-out
[ $? -eq 0 ] ; ERROR_CHECK "TEST2-1 error"
diff $tmp-ans2 $tmp-ng
[ $? -eq 0 ] ; ERROR_CHECK "TEST2-2 error"

###########################################
#TEST3

cat << FIN > $tmp-tran
CCC 003 太田
AAA 001 山田
BBB 002 上田
DDD 004 堅田
FIN

cat << FIN > $tmp-master
002 上田 富山
003 太田 石川
FIN

cat << FIN > $tmp-out
CCC 003 太田 石川
BBB 002 上田 富山
FIN

${com} key=2/3 $tmp-master $tmp-tran > $tmp-ans
diff $tmp-ans $tmp-out
[ $? -eq 0 ] ; ERROR_CHECK "TEST3 error"

###########################################
#TEST4

cat << FIN > $tmp-tran
DDD 004 堅田 へへへ
AAA 001 山田 あはは
CCC 003 太田 ふふふ
BBB 002 上田 おほほ
FIN

cat << FIN > $tmp-master
002 上田 富山
003 太田 石川
FIN

cat << FIN > $tmp-out
CCC 003 太田 石川 ふふふ
BBB 002 上田 富山 おほほ
FIN

${com} key=2/3 $tmp-master $tmp-tran > $tmp-ans
diff $tmp-ans $tmp-out
[ $? -eq 0 ] ; ERROR_CHECK "TEST4 error"

rm -f $tmp-*
echo "${name}" OK
exit 0