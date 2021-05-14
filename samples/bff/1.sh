errorCount=`grep 'error:' SwiftCheck.log | wc -l`
warningCount=`grep 'warning:' SwiftCheck.log | wc -l`
contentStr=""
if (($errorCount > 0)); then
  contentStr="代码检查完毕:\n"$errorCount"个swift错误\n"$warningCount"个swift警告,代码检查未通过"
  echo $contentStr
else
  contentStr="代码检查完毕:没有错误,代码检查通过"
  echo $contentStr
fi

json="{\"msgtype\":\"text\",\"text\": {\"content\":\"$contentStr\"}}"

  echo $json
curl 'https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=2b3125e6-1a98-4b45-87a1-1c2627fb87b2' \
-H 'Content-Type: application/json' \
-d "$json"