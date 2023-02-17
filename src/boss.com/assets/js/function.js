//检测是否已登录状态
var admintoken = localStorage.getItem('admintoken');
var admincode = localStorage.getItem('admincode');
//alert(admincode+'=>'+admintoken)
var pagename = window.location.pathname;

if (pagename == "/page/login.html") {
    if (admintoken != null && admincode != null) {
        $.post(server_url + 'checkadminlogin', {admincode: admincode, admintoken: admintoken}, function (res) {
            if (res.code === 0) {
                location.replace('../../../');
                /*layer.msg(res.msg, {icon: 1, time: 500}, function () {
                    location.replace('../../../');
                });*/
            } else {
                localStorage.removeItem("admincode")
                localStorage.removeItem("admintoken")
                localStorage.removeItem("admin_id")
                localStorage.removeItem("account")
                //layer.msg(res.msg, {icon: 2, anim: 6});
            }
        }, 'json');
    }
} else {
    if (admintoken != null && admincode != null) {
        $.post(server_url + 'checkadminlogin', {admincode: admincode, admintoken: admintoken}, function (res) {
            if (res.code === 0) {
                console.log("目前登录状态");
            } else {
                localStorage.removeItem("admincode")
                localStorage.removeItem("admintoken")
                localStorage.removeItem("admin_id")
                localStorage.removeItem("account")
                top.location.replace('/page/login.html');
                /*layer.msg(res.msg, {icon: 2, anim: 6, time: 500}, function () {
                    top.location.replace('/page/login.html');
                });*/
                //layer.msg(res.msg, {icon: 2, anim: 6});
            }
        }, 'json');
    } else {
        top.location.replace('/page/login.html');
    }
}

//获取url参数
function getUrlParams(key) {
    var url = window.location.search.substr(1);
    if (url == '') {
        return false;
    }
    var paramsArr = url.split('&');
    for (var i = 0; i < paramsArr.length; i++) {
        var combina = paramsArr[i].split("=");
        if (combina[0] == key) {
            return combina[1];
        }
    }
    return false;
}

//获取url的分页码参数
function getUrlPage(key) {
    var url = window.location.hash.substr(2);
    if (url == '') {
        return 1;
    }
    var paramsArr = url.split('#!');
    for (var i = 0; i < paramsArr.length; i++) {
        var combina = paramsArr[i].split("=");
        if (combina[0] == key) {
            return combina[1];
        }
    }
    return false;
}

//图片转换base64
function getImgBase64(imgUrl) {
    return new Promise((resolve) => {
        window.URL = window.URL || window.webkitURL;
        let xhr = new XMLHttpRequest();
        xhr.open("get", imgUrl, true);
        // 至关重要
        xhr.responseType = "blob";
        xhr.onload = function () {
            if (this.status === 200) {
                //得到一个blob对象
                let blob = this.response;
                //console.log("blob", blob)
                // 至关重要
                let oFileReader = new FileReader();
                oFileReader.onloadend = function (e) {
                    let base64 = e.target.result;
                    resolve(base64);
                    //console.log("方式一》》》》》》》》》", base64);
                };
                oFileReader.readAsDataURL(blob);
                //====为了在页面显示图片，可以删除====
                /*let img = document.createElement("img");
                img.onload = function (e) {
                    window.URL.revokeObjectURL(img.src); // 清除释放
                };
                img.src = window.URL.createObjectURL(blob)
                document.getElementById("container1").appendChild(img);*/
                //====为了在页面显示图片，可以删除====
            }
        }
        xhr.send();
    });
}




