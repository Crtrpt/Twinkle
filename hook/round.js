
//获取请求的url
console.log(GetUrl())
//获取请求的方法
console.log(GetMethod())
//获取请求体
console.log(GetRequestBody())
abc = 2 + 2;
//设置请求参数
SetRequestHeader("AA","bb")
//获取请求参数
console.log(GetRequestHeader("Aa"))
//设置返回参数
SetResponseHeader("Test","111")
//获取返回参数
console.log(GetResponseHeader("Test"))

Run()
//设置返回码
//SetCode(404)
//SetBody("1111")