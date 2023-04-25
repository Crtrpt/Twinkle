export default class api {
    static api = import.meta.env.VITE_API_BASE_URL;

    static headers = {
        "Content-Type": "application/json",
        Token: "",
    };

    // GET 请求
    static get(
        path: string,
        query: Object
    ): Promise<any> {
        // query=JSON.parse(JSON.stringify(query));
        var queryStr = Object.keys(query)
            .map((key) => `${key}=${query[key]}`)
            .join("&");
        var header = { ...api.headers };

        return fetch(api.api + path + "?" + queryStr, {
            method: "GET",
            headers: header,
        }).then((res) => {
            return res.json();
        });
    }

     // POST 请求
    static post(
        path: string,
        param: any
    ): Promise<any> {
        var header = { ...api.headers };
        return fetch(api.api + path, {
            method: "POST",
            headers: header,
            body: JSON.stringify(param),
        }).then((res) => {
            return res.json();
        });
    }

     // 上传 请求
    static upload(
        path: string,
        param: any,
    ): Promise<any> {
        var headers = JSON.parse(JSON.stringify(api.headers));

        // headers["Content-Type"] = "multipart/form-data; boundary=something";
        delete headers["Content-Type"];
        var form = new FormData();

        for (const name in param) {
            form.append(name, param[name]);
        }
        return fetch(api.api + path, {
            method: "POST",
            headers: headers,
            body: form,
        }).then((res) => {
            return res.json();
        });
    }

    // 下载
    static download(
        path: string,
        query: any,
        option: {} = { service: "api" }
    ): Promise<any> {
        var headers = JSON.parse(JSON.stringify(api.headers));
        delete headers["Content-Type"];
        var queryStr = Object.keys(query)
            .map((key) => `${key}=${query[key]}`)
            .join("&");
        return fetch(api.api + path + "?" + queryStr, {
            method: "GET",
            headers: headers,
        }).then((res) => {
            return res.text();
        });
    }
}