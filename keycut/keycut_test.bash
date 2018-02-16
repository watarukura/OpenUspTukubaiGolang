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
name=keycut # test対象コマンドの名前
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

cat << FIN > $tmp-input
0000000 浜地______ 50 F 91 59 20 76 54
0000001 鈴田______ 50 F 46 39 8  5  21
0000003 杉山______ 26 F 30 50 71 36 30
0000004 白土______ 40 M 58 71 20 10 6
0000005 崎村______ 50 F 82 79 16 21 80
FIN

cat << FIN > $tmp-out
$tmp-0000000
$tmp-0000001
$tmp-0000003
$tmp-0000004
$tmp-0000005
FIN

${com} "$tmp-%1" $tmp-input
echo $tmp-000000? | tr ' ' '\n' > $tmp-ans
diff $tmp-ans $tmp-out
[ $? -eq 0 ] ; ERROR_CHECK "TEST1.1 error"

cat $tmp-out		|
xargs cat		|
diff $tmp-input -
[ $? -eq 0 ] ; ERROR_CHECK "TEST1.2 error"

cat $tmp-out | xargs rm
###########################################
#TEST2
cat << FIN > $tmp-out
$tmp-0000000/0000000
$tmp-0000001/0000001
$tmp-0000003/0000003
$tmp-0000004/0000004
$tmp-0000005/0000005
FIN

${com} "$tmp-%1/%1" $tmp-input
echo $tmp-000000?/* | tr ' ' '\n' > $tmp-ans
diff $tmp-ans $tmp-out
[ $? -eq 0 ] ; ERROR_CHECK "TEST2 error"

rm -Rf $tmp-*
###########################################
#TEST3
cat << FIN > $tmp-input
0000000 浜地______ 50 F 91 59 20 76 54
0000001 鈴田______ 50 F 46 39 8  5  21
0000003 杉山______ 26 F 30 50 71 36 30
0000004 白土______ 40 M 58 71 20 10 6
0000005 崎村______ 50 F 82 79 16 21 80
FIN

cat << FIN > $tmp-out
$tmp-0000000
$tmp-0000001
$tmp-0000003
$tmp-0000004
$tmp-0000005
FIN

${com} -d "$tmp-%1" $tmp-input
echo $tmp-000000? | tr ' ' '\n' > $tmp-ans
diff $tmp-ans $tmp-out
[ $? -eq 0 ] ; ERROR_CHECK "TEST3.1 error"

cat $tmp-out		|
xargs cat		|
diff <(sed 's/^[^ ]* //' $tmp-input | sed 's/  */ /g') -
[ $? -eq 0 ] ; ERROR_CHECK "TEST3.2 error"

rm -Rf $tmp-*
echo "${name}" OK
exit 0