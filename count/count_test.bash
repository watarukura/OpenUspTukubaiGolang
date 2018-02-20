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
name=count # test対象コマンドの名前
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

cat << FIN > $tmp-in
001 1 942
001 1.3 421
002 -123.0 111
002 123.0 11
FIN

cat << FIN > $tmp-out
001 2
002 2
FIN

${com} 1 1 $tmp-in > $tmp-ans
diff $tmp-ans $tmp-out
[ $? -eq 0 ] ; ERROR_CHECK "TEST1.1 error"

cat $tmp-in | ${com} 1 1 > $tmp-ans
diff $tmp-ans $tmp-out
[ $? -eq 0 ] ; ERROR_CHECK "TEST1.2 error"

###########################################
#TEST2

cat << FIN > $tmp-in
001 江頭 1 942
001 江頭 1.3 421
002 上山田 -123.0 111
002 上田 123.0 11
FIN

cat << FIN > $tmp-out
001 江頭 2
002 上山田 1
002 上田 1
FIN

${com} 1 2 $tmp-in > $tmp-ans
diff $tmp-ans $tmp-out
[ $? -eq 0 ] ; ERROR_CHECK "TEST2.1 error"

rm -Rf $tmp-*
echo "${name}" OK
exit 0