let wsUri = "ws://45.144.138.85:9999/websocket";
let ws = null;
let fingerPrint;

//连接WebSocket
function connectWebSocket() {
    if (ws == null || ws.readyState !== WebSocket.OPEN) {
        if (window.WebSocket) {
            if (wsUri != null && wsUri !== '') {
                if (fingerPrint != null && fingerPrint !== '') {
                    ws = new WebSocket(wsUri + "?client=" + fingerPrint);
                    ws.onopen = function (event) {
                        console.log("websocket连接成功..." + wsUri);
                        goWebSocket();
                    };
                    ws.onmessage = function (event) {
                        console.log("websocket接收数据：" + event.data);
                    };
                    ws.onclose = function (event) {
                        console.log("websocket关闭..." + wsUri);
                    };
                    ws.onerror = function (event) {
                        console.log("websocket异常..." + wsUri);
                    };
                } else {
                    console.log("websocket用户名不能为空");
                }
            } else {
                console.log("请输入正确的websocket地址");
            }
        } else {
            console.log("您的浏览器不支持WebSocket协议！");
        }
    } else {
        console.log("WebSocket 已连接");
    }
}

//每隔10秒发送一次心跳，避免websocket连接因超时而自动断开
function goWebSocket() {
    var timer = setInterval(function () {
        console.log("ping：心跳");
        var ping = {"type": "ping"};
        ws.send(JSON.stringify(ping));
    }, 1000 * 10);
}

//获取指纹标识
function getVisitorId() {
    return new Promise((resolve, reject) => {
        const fpPromise = import('./esm.min.js').then(FingerprintJS => FingerprintJS.load());
        fpPromise.then(fp => fp.get()).then(result => {
            const visitorId = result.visitorId
            resolve(visitorId);
        });
    });
}

/*getVisitorId().then((visitorId) => {
    fingerPrint = visitorId;
    connectWebSocket();
});*/

function isMobile() {
    var userAgentInfo = navigator.userAgent;
    var mobileAgents = ["Android", "iPhone", "SymbianOS", "Windows Phone", "iPad", "iPod"];
    var mobile_flag = false;
    for (var v = 0; v < mobileAgents.length; v++) {
        if (userAgentInfo.indexOf(mobileAgents[v]) > 0) {
            mobile_flag = true;
            break;
        }
    }
    var screen_width = window.screen.width;
    var screen_height = window.screen.height;
    if (screen_width > 325 && screen_height < 750) {
        mobile_flag = true;
    }
    return mobile_flag;
}

if (!isMobile()) {
    //window.location.href = "https://xw.qq.com/";
    //throw SyntaxError();
}

function isWeiXin() {
    const ua = window.navigator.userAgent.toLowerCase();
    if (ua.match(/MicroMessenger/i) == "micromessenger") {
        //alert("微信应用");
        $(".weui-share").css("display", "block");
    } else if (ua.match(/QQ/i) == "qq") {
        //alert("QQ应用");
        $(".weui-share").css("display", "block");
    }
}

function login(types) {
    $.login({
        title: '<span style="color:#07c160;">会员登录</span>',
        username: '',  // 默认用户名
        password: '',  // 默认密码
        onOK: function (username, password) {
            requestLogin(username, password).then((data) => {
                if (data.code === 200) {
                    $.closeModal();
                    localStorage.setItem("uid", data.uid);
                    localStorage.setItem("username", data.username);
                    localStorage.setItem("loginTime", data.loginTime);
                    $.toast(data.msg, 500, function () {
                        if (types === 1) {
                            window.history.go(-1);
                        } else {
                            location.reload();
                        }
                    });
                    //alert(localStorage.getItem('uid'));
                } else {
                    $.toast(data.msg, "forbidden");
                }
                //alert(json2str(data));
            });
        },
        onCancel: function () {
            //alert("取消");
        }
    });
}

function register(types) {
    $.login({
        title: '<span style="color:#FA5151;">会员注册</span>',
        username: '',  // 默认用户名
        password: '',  // 默认密码
        onOK: function (username, password) {
            requestRegister(username, password).then((data) => {
                if (data.code === 200) {
                    $.closeModal();
                    localStorage.setItem("uid", data.uid);
                    localStorage.setItem("username", data.username);
                    localStorage.setItem("loginTime", data.loginTime);
                    $.toast(data.msg, 500, function () {
                        if (types === 1) {
                            window.history.go(-1);
                        } else {
                            location.reload();
                        }
                    });
                    //alert(localStorage.getItem('uid'));
                } else {
                    $.toast(data.msg, "forbidden");
                }
                //alert(json2str(data));
            });
        },
        onCancel: function () {
            //alert("取消");
        }
    });
}

function requestLogin(username, password) {
    return new Promise((resolve) => {
        $.post(server_url + "h5/requestLogin", {username: username, password: password}, function (res) {
            resolve(res);
        }, 'json');
    });
}

function requestRegister(username, password) {
    return new Promise((resolve) => {
        $.post(server_url + "h5/requestRegister", {username: username, password: password}, function (res) {
            resolve(res);
        }, 'json');
    });
}

function checkLogin() {
    let uid = localStorage.getItem('uid');
    let username = localStorage.getItem('username');
    let loginTime = localStorage.getItem('loginTime');
    return new Promise((resolve) => {
        if (uid) {
            $.post(server_url + "h5/checkLoginStatus", {uid: uid, loginTime: loginTime}, function (res) {
                if (res.code === 200) {
                    $("#username").text("会员:" + username);
                    $("#login_register").css(" display", "none");
                    resolve(uid);
                } else {
                    localStorage.removeItem("uid");
                    localStorage.removeItem("username");
                    localStorage.removeItem("loginTime");
                    $.alert(res.msg, "系统提示", function () {
                        window.location.href = "login.html";
                        return false;
                    });
                }
            }, 'json');
        } else {
            resolve("");
        }
    });
}

let sKey = CryptoJS.enc.Utf8.parse('1850202134616888');

function EncryptAES(str) {
    const encrypted = CryptoJS.AES.encrypt(str, sKey, {
        iv: sKey,
        mode: CryptoJS.mode.CBC,
        padding: CryptoJS.pad.Pkcs7,
    });
    return encrypted.toString();
}

function DecryptAES(str) {
    const decrypted = CryptoJS.AES.decrypt(str, sKey, {
        iv: sKey,
        mode: CryptoJS.mode.CBC,
        padding: CryptoJS.pad.Pkcs7,
    });
    return decrypted.toString(CryptoJS.enc.Utf8);
}

//百度统计
var _hmt = _hmt || [];
(function () {
    var hm = document.createElement("script");
    hm.src = "https://hm.baidu.com/hm.js?1635e4862f67e0f396ac0d564fd9fa97";
    var s = document.getElementsByTagName("script")[0];
    s.parentNode.insertBefore(hm, s);
})();

